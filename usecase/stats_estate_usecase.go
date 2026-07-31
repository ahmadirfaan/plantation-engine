package usecase

import (
	"errors"
	"log"
	"math"
	"sort"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/helper"
	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/ahmadirfaan/plantation-engine/repository"
)

type StatsEstateUseCase interface {
	PublishCalculationStatsEstate(estateId string)
	GetStatsEstate(estateId string) (generated.GetStatsEstateResponse, error)
	GetDronePlanDistance(estateId string, params generated.GetDistanceForDronePlanParams) (generated.GetDronePlanDistance, error)
}

func NewStatsEstateUseCase(statsEstateRepository repository.StatsEstateRepository, estateRepo repository.EstateRepository) StatsEstateUseCase {
	statsChannel := make(chan string, 20)
	dronePlanCache := repository.NewDronePlanCache()
	statsEstateUseCase := &statsEstateUseCase{
		statsEstateRepository: statsEstateRepository,
		estateRepository:      estateRepo,
		statsEstateChannel:    statsChannel,
		dronePlanCache:        dronePlanCache,
	}
	go statsEstateUseCase.startWorker()
	return statsEstateUseCase
}

type statsEstateUseCase struct {
	statsEstateRepository repository.StatsEstateRepository
	estateRepository      repository.EstateRepository
	statsEstateChannel    chan string
	dronePlanCache        *repository.DronePlanCache
}

func (s *statsEstateUseCase) GetDronePlanDistance(estateId string, params generated.GetDistanceForDronePlanParams) (generated.GetDronePlanDistance, error) {
	response := generated.GetDronePlanDistance{}
	statsEstate, err := s.queryModelStatsEstate(estateId)
	if err != nil {
		return response, err
	}
	if statsEstate == nil {
		return generated.GetDronePlanDistance{
			Distance: 0,
		}, nil
	}
	if params.MaxDistance != nil && *params.MaxDistance > 0 {
		maxDistance := *params.MaxDistance
		restX, restY := s.findRestCoordinate(estateId, maxDistance, statsEstate)
		if restX > 0 || restY > 0 {
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

func (s *statsEstateUseCase) findRestCoordinate(estateId string, maxDistance int, statsEstate *model.EstateStats) (int, int) {
	if maxDistance > statsEstate.TotalDistanceDrone {
		return helper.LastBlockFinish(statsEstate.Length, statsEstate.Width)
	}
	restCoordinateExist, isCacheExist := s.dronePlanCache.Get(estateId, maxDistance)
	if isCacheExist {
		return restCoordinateExist.X, restCoordinateExist.Y
	}

	trees, err := s.statsEstateRepository.QueryAllTree(estateId)
	if err != nil || len(trees) == 0 {
		log.Println("error when query all trees")
		return 0, 0
	}

	restCoordinate := helper.CheckRestCoordinate(maxDistance, trees, statsEstate.Length, statsEstate.Width)
	s.dronePlanCache.Set(estateId, maxDistance, restCoordinate)
	return restCoordinate.X, restCoordinate.Y

}

func (s *statsEstateUseCase) GetStatsEstate(estateId string) (generated.GetStatsEstateResponse, error) {
	response := generated.GetStatsEstateResponse{}
	statsEstate, err := s.queryModelStatsEstate(estateId)
	if err != nil {
		return response, err
	}
	if statsEstate == nil {
		response := generated.GetStatsEstateResponse{
			Count:  0,
			Max:    0,
			Min:    0,
			Median: 0,
		}
		return response, nil
	}

	response.Median = float32(statsEstate.MedianHeightTree)
	response.Count = statsEstate.SumTree
	response.Max = statsEstate.MaxHeightTree
	response.Min = statsEstate.MinHeightTree
	return response, nil
}

func (s *statsEstateUseCase) queryModelStatsEstate(estateId string) (*model.EstateStats, error) {
	estate, err := s.estateRepository.QueryByEstateId(estateId)
	if err != nil {
		return nil, err
	}

	if estate == nil {
		return nil, errors.New("404|estate not found")
	}

	statsEstate, err := s.statsEstateRepository.QueryById(estateId)
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

func (s *statsEstateUseCase) PublishCalculationStatsEstate(estateId string) {
	s.statsEstateChannel <- estateId
}

func (s *statsEstateUseCase) startWorker() {
	for estateId := range s.statsEstateChannel {
		// Process the stats calculation when notified
		s.calculateStatsForEstate(estateId)
	}
}

func (s *statsEstateUseCase) calculateStatsForEstate(estateId string) {
	estate, _ := s.estateRepository.QueryByEstateId(estateId)
	if estate == nil {
		return
	}

	trees, err := s.statsEstateRepository.QueryAllTree(estateId)
	if err != nil || len(trees) == 0 {
		log.Println("error when query all trees")
		return
	}

	statsEstate := calculateSummaryStatsEstate(trees)
	statsEstate.EstateId = estateId

	helper.CalculateTotalDistanceDrone(trees, *estate, &statsEstate)

	err = s.statsEstateRepository.SaveStatsEstate(statsEstate)
	if err != nil {
		log.Println("error when save summary estate")
		return
	}

	s.dronePlanCache.ClearByEstateId(estateId)

}

func calculateSummaryStatsEstate(trees []model.Tree) model.EstateStats {

	heights, minHeight, maxHeight := calculateHeight(trees)
	median := calculateMedianHeightTree(heights)

	return model.EstateStats{
		SumTree:          len(trees),
		MinHeightTree:    minHeight,
		MaxHeightTree:    maxHeight,
		MedianHeightTree: median,
	}
}

func calculateMedianHeightTree(heights []float64) float64 {
	sort.Float64s(heights)
	var median float64
	n := len(heights)
	if n%2 == 1 {
		// Odd count
		median = heights[n/2]
	} else {
		// Even count
		median = (heights[n/2-1] + heights[n/2]) / 2
	}
	return median
}

func calculateHeight(trees []model.Tree) ([]float64, int, int) {
	heights := make([]float64, len(trees))
	minHeight := math.MaxInt
	maxHeight := -math.MaxInt
	for i, tree := range trees {
		h := tree.Height
		heights[i] = float64(h)
		if tree.Height < minHeight {
			minHeight = tree.Height
		}
		if tree.Height > maxHeight {
			maxHeight = tree.Height
		}
	}
	return heights, minHeight, maxHeight
}
