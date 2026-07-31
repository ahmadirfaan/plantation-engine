package usecase

import (
	"errors"
	"strings"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/repository"
)

type TreeUseCase interface {
	AddTreeToEstate(estateId string, request generated.CreateTreeRequest) (*string, error)
}

func NewTreeUseCase(estateRepo repository.EstateRepository, treeRepo repository.TreeRepository,
	blockRepo repository.BlockRepository, statsEstateUseCase StatsEstateUseCase) TreeUseCase {
	return &treeUseCase{
		estateRepository:   estateRepo,
		treeRepository:     treeRepo,
		blockRepository:    blockRepo,
		statsEstateUseCase: statsEstateUseCase,
	}
}

type treeUseCase struct {
	estateRepository   repository.EstateRepository
	treeRepository     repository.TreeRepository
	blockRepository    repository.BlockRepository
	statsEstateUseCase StatsEstateUseCase
}

func (t treeUseCase) AddTreeToEstate(estateId string, request generated.CreateTreeRequest) (*string, error) {

	estate, err := t.estateRepository.QueryByEstateId(estateId)
	if err != nil {
		return nil, err
	}

	if estate == nil || strings.TrimSpace(estate.Id) == "" {
		return nil, errors.New("404|estate does not exist")
	}

	if estate.Width < request.Y {
		return nil, errors.New("400|Y request above from data")
	}

	if estate.Length < request.X {
		return nil, errors.New("400|X request above from data")
	}

	block, err := t.blockRepository.QueryByEstateIdAndBlockCoordinate(estateId, request.X, request.Y)
	if err != nil {
		return nil, err
	}

	if block != nil {
		return nil, errors.New("409|block already has a tree")
	}

	blockId, err := t.blockRepository.SaveBlock(estateId, request.X, request.Y)
	if err != nil {
		return nil, err
	}
	treeId, err := t.treeRepository.SaveTree(*blockId, estateId, request.Height)
	if err != nil {
		return nil, err
	}

	t.statsEstateUseCase.PublishCalculationStatsEstate(estateId)

	return treeId, nil
}
