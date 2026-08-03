package usecase

import (
	"context"
	"errors"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/repository"
)

type TreeUseCase interface {
	AddTreeToEstate(ctx context.Context, estateId string, request generated.CreateTreeRequest) (*string, error)
}

func NewTreeUseCase(estateRepo repository.EstateRepository, treeRepo repository.TreeRepository, blockRepo repository.BlockRepository, statsEstateUseCase StatsEstateUseCase) TreeUseCase {
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

func (t treeUseCase) AddTreeToEstate(ctx context.Context, estateId string, request generated.CreateTreeRequest) (*string, error) {
	estate, err := t.estateRepository.QueryByEstateId(ctx, estateId)
	if err != nil {
		return nil, err
	}
	if estate == nil {
		return nil, errors.New("404|estate does not exist")
	}
	if estate.Width < request.Y {
		return nil, errors.New("400|Y request above from data")
	}
	if estate.Length < request.X {
		return nil, errors.New("400|X request above from data")
	}

	block, err := t.blockRepository.QueryByEstateIdAndBlockCoordinate(ctx, estateId, request.X, request.Y)
	if err != nil {
		return nil, err
	}
	if block != nil {
		return nil, errors.New("409|block already has a tree")
	}

	treeId, err := t.treeRepository.SaveBlockAndTree(ctx, estateId, request.X, request.Y, request.Height)
	if err != nil {
		if errors.Is(err, repository.ErrBlockHasTree) {
			return nil, errors.New("409|block already has a tree")
		}
		return nil, err
	}

	t.statsEstateUseCase.PublishCalculationStatsEstate(estateId)

	return treeId, nil
}
