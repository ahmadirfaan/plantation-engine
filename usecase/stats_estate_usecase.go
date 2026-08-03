package usecase

import (
	"context"
	"errors"
	"log/slog"
	"sort"
	"sync"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/helper"
	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/ahmadirfaan/plantation-engine/repository"
)

type StatsEstateUseCase interface {
	PublishCalculationStatsEstate(estateId string)
	GetStatsEstate(ctx context.Context, estateId string) (generated.GetStatsEstateResponse, error)
	GetDronePlanDistance(ctx context.Context, estateId string, params generated.GetDistanceForDronePlanParams) (generated.GetDronePlanDistance, error)
}

func NewStatsEstateUseCase(statsEstateRepository repository.StatsEstateRepository, estateRepository repository.EstateRepository) StatsEstateUseCase {
	s := &statsEstateUseCase{
		statsEstateRepository: statsEstateRepository,
		estateRepository:      estateRepository,
		dronePlanCache:        repository.NewDronePlanCache(),
		dirty:                 make(map[string]struct{}),
		notifyCh:              make(chan struct{}, 1),
	}
	go s.startWorker()
	return s
}

type statsEstateUseCase struct {
	statsEstateRepository repository.StatsEstateRepository
	estateRepository      repository.EstateRepository
	dronePlanCache        *repository.DronePlanCache

	mu       sync.Mutex
	dirty    map[string]struct{}
	notifyCh chan struct{}
}

// PublishCalculationStatsEstate marks the estate as dirty and signals the
// worker. Multiple calls for the same estate are coalesced into a single
// recalculation per drain cycle, and the send never blocks the caller.
func (s *statsEstateUseCase) PublishCalculationStatsEstate(estateId string) {
	s.mu.Lock()
	s.dirty[estateId] = struct{}{}
	s.mu.Unlock()

	select {
	case s.notifyCh <- struct{}{}:
	default:
	}
}

func (s *statsEstateUseCase) startWorker() {
	for range s.notifyCh {
		s.mu.Lock()
		pending := make([]string, 0, len(s.dirty))
		for id := range s.dirty {
			pending = append(pending, id)
		}
		s.dirty = make(map[string]struct{})
		s.mu.Unlock()

		for _, estateId := range pending {
			s.calculateStatsForEstate(context.Background(), estateId)
		}
	}
}

func (s *statsEstateUseCase) GetDronePlanDistance(ctx context.Context, estateId string, params generated.GetDistanceForDronePlanParams) (generated.GetDronePlanDistance, error) {
	response := generated.GetDronePlanDistance{}

	statsEstate, err := s.queryModelStatsEstate(ctx, estateId)
	if err != nil {
		return generated.GetDronePlanDistance{}, err
	}
	if statsEstate == nil {
		return response, nil
	}

	if params.MaxDistance != nil && *params.MaxDistance > 0 {
		maxDistance := *params.MaxDistance
		restX, restY := s.findRestCoordinate(ctx, estateId, maxDistance, statsEstate)
		if restX > 0 && restY > 0 {
			response.Rest = &struct {
				X int `json:"x"`
				Y int `json:"y"`
			}{
				X: restX,
				Y: restY,
			}
		}
	}
	response.Distance = statsEstate.TotalDistanceDrone

	return response, nil
}

func (s *statsEstateUseCase) findRestCoordinate(ctx context.Context, estateId string, maxDistance int, statsEstate *model.EstateStats) (int, int) {
	if maxDistance > statsEstate.TotalDistanceDrone {
		x, y := helper.LastBlockFinish(statsEstate.Length, statsEstate.Width)
		return x, y
	}

	if coordinate, exists := s.dronePlanCache.Get(estateId, maxDistance); exists {
		return coordinate.X, coordinate.Y
	}

	trees, err := s.statsEstateRepository.QueryAllTree(ctx, estateId)
	if err != nil || len(trees) == 0 {
		slog.Error("error when query all trees", "estateId", estateId, "error", err)
		return 0, 0
	}

	restCoordinate := helper.CheckRestCoordinate(maxDistance, trees, statsEstate.Length, statsEstate.Width)
	s.dronePlanCache.Set(estateId, maxDistance, restCoordinate)

	return restCoordinate.X, restCoordinate.Y
}

func (s *statsEstateUseCase) GetStatsEstate(ctx context.Context, estateId string) (generated.GetStatsEstateResponse, error) {
	statsEstate, err := s.queryModelStatsEstate(ctx, estateId)
	if err != nil {
		return generated.GetStatsEstateResponse{}, err
	}
	if statsEstate == nil {
		return generated.GetStatsEstateResponse{}, nil
	}

	response := generated.GetStatsEstateResponse{
		Median: float32(statsEstate.MedianHeightTree),
		Count:  statsEstate.SumTree,
		Max:    statsEstate.MaxHeightTree,
		Min:    statsEstate.MinHeightTree,
	}

	return response, nil
}

func (s *statsEstateUseCase) queryModelStatsEstate(ctx context.Context, estateId string) (*model.EstateStats, error) {
	estate, err := s.estateRepository.QueryByEstateId(ctx, estateId)
	if err != nil {
		return nil, err
	}
	if estate == nil {
		return nil, errors.New("404|estate not found")
	}

	statsEstate, err := s.statsEstateRepository.QueryById(ctx, estateId)
	if err != nil {
		return nil, err
	}
	if statsEstate == nil {
		return nil, nil
	}
	statsEstate.Width = estate.Width
	statsEstate.Length = estate.Length

	return statsEstate, nil
}

func (s *statsEstateUseCase) calculateStatsForEstate(ctx context.Context, estateId string) {
	estate, err := s.estateRepository.QueryByEstateId(ctx, estateId)
	if err != nil || estate == nil {
		if err != nil {
			slog.Error("failed to query estate for stats", "estateId", estateId, "error", err)
		}
		return
	}

	trees, err := s.statsEstateRepository.QueryAllTree(ctx, estateId)
	if err != nil || len(trees) == 0 {
		if err != nil {
			slog.Error("failed to query all trees for stats", "estateId", estateId, "error", err)
		}
		return
	}

	statsEstate := s.calculateSummaryStatsEstate(trees)
	statsEstate.EstateId = estateId

	helper.CalculateTotalDistanceDrone(trees, *estate, &statsEstate)

	if err := s.statsEstateRepository.SaveStatsEstate(ctx, statsEstate); err != nil {
		slog.Error("failed to save stats estate", "estateId", estateId, "error", err)
		return
	}

	s.dronePlanCache.ClearByEstateId(estateId)
}

func (s *statsEstateUseCase) calculateSummaryStatsEstate(trees []model.Tree) model.EstateStats {
	return calculateSummaryStatsEstate(trees)
}

func calculateSummaryStatsEstate(trees []model.Tree) model.EstateStats {
	return model.EstateStats{
		SumTree:          calculateHeight(trees),
		MinHeightTree:    calculateMinHeightTree(trees),
		MaxHeightTree:    calculateMaxHeightTree(trees),
		MedianHeightTree: calculateMedianHeightTree(trees),
	}
}

func calculateHeight(trees []model.Tree) int {
	return len(trees)
}

func calculateMinHeightTree(trees []model.Tree) int {
	if len(trees) == 0 {
		return 0
	}
	min := trees[0].Height
	for _, tree := range trees {
		if tree.Height < min {
			min = tree.Height
		}
	}
	return min
}

func calculateMaxHeightTree(trees []model.Tree) int {
	if len(trees) == 0 {
		return 0
	}
	max := trees[0].Height
	for _, tree := range trees {
		if tree.Height > max {
			max = tree.Height
		}
	}
	return max
}

func calculateMedianHeightTree(trees []model.Tree) float64 {
	if len(trees) == 0 {
		return 0
	}

	heights := make([]float64, len(trees))
	for i, tree := range trees {
		heights[i] = float64(tree.Height)
	}

	sort.Float64s(heights)

	n := len(heights)
	if n%2 == 1 {
		return heights[n/2]
	}

	middle := n / 2
	return (heights[middle-1] + heights[middle]) / 2
}
