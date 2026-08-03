package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/google/uuid"
)

// mockEstateRepositoryError returns errors on every estate operation.
type mockEstateRepositoryError struct{}

func (m *mockEstateRepositoryError) QueryByEstateId(ctx context.Context, id string) (*model.Estate, error) {
	return nil, errors.New("error Connection Timeout")
}

func (m *mockEstateRepositoryError) SaveEstate(ctx context.Context, width int, length int) (*string, error) {
	return nil, errors.New("mock create estate error")
}

// mockEstateRepositoryQueryIdNotExist returns no estate and no error.
type mockEstateRepositoryQueryIdNotExist struct{}

func (m *mockEstateRepositoryQueryIdNotExist) QueryByEstateId(ctx context.Context, id string) (*model.Estate, error) {
	return nil, nil
}

func (m *mockEstateRepositoryQueryIdNotExist) SaveEstate(ctx context.Context, width int, length int) (*string, error) {
	return nil, nil
}

// mockEstateRepositoryQueryIdEstateExist returns an existing 3x3 estate.
type mockEstateRepositoryQueryIdEstateExist struct{}

func (m *mockEstateRepositoryQueryIdEstateExist) QueryByEstateId(ctx context.Context, id string) (*model.Estate, error) {
	timeNow := time.Now()
	return &model.Estate{
		Id:        id,
		Width:     3,
		Length:    3,
		CreatedAt: &timeNow,
		UpdatedAt: &timeNow,
	}, nil
}

func (m *mockEstateRepositoryQueryIdEstateExist) SaveEstate(ctx context.Context, width int, length int) (*string, error) {
	return nil, nil
}

// mockEstateRepositorySuccess returns a fresh id on save, no estate on query.
type mockEstateRepositorySuccess struct{}

func (m *mockEstateRepositorySuccess) QueryByEstateId(ctx context.Context, id string) (*model.Estate, error) {
	return nil, nil
}

func (m *mockEstateRepositorySuccess) SaveEstate(ctx context.Context, width int, length int) (*string, error) {
	uuidString := uuid.New().String()
	return &uuidString, nil
}

// mockTreeRepositoryError returns an error on SaveBlockAndTree.
type mockTreeRepositoryError struct{}

func (m *mockTreeRepositoryError) SaveBlockAndTree(ctx context.Context, estateId string, x int, y int, height int) (*string, error) {
	return nil, errors.New("connection timeout")
}

// mockTreeRepositorySaveTreeSuccess returns a fresh id on SaveBlockAndTree.
type mockTreeRepositorySaveTreeSuccess struct{}

func (m *mockTreeRepositorySaveTreeSuccess) SaveBlockAndTree(ctx context.Context, estateId string, x int, y int, height int) (*string, error) {
	uuidString := uuid.New().String()
	return &uuidString, nil
}

// mockBlockRepositoryError returns an error when querying a block.
type mockBlockRepositoryError struct{}

func (m *mockBlockRepositoryError) QueryByEstateIdAndBlockCoordinate(ctx context.Context, estateId string, x int, y int) (*model.Block, error) {
	return nil, errors.New("connection timeout")
}

// mockBlockRepositoryBlockExist returns an existing block.
type mockBlockRepositoryBlockExist struct{}

func (m *mockBlockRepositoryBlockExist) QueryByEstateIdAndBlockCoordinate(ctx context.Context, estateId string, x int, y int) (*model.Block, error) {
	timeNow := time.Now()
	id := uuid.New().String()
	return &model.Block{
		Id:        id,
		EstateId:  estateId,
		BlockX:    x,
		BlockY:    y,
		CreatedAt: &timeNow,
		UpdatedAt: &timeNow,
	}, nil
}

// mockBlockRepositorySuccess returns no existing block.
type mockBlockRepositorySuccess struct{}

func (m *mockBlockRepositorySuccess) QueryByEstateIdAndBlockCoordinate(ctx context.Context, estateId string, x int, y int) (*model.Block, error) {
	return nil, nil
}

// mockStatsEstateUseCase is a no-op stats use case.
type mockStatsEstateUseCase struct{}

func (m *mockStatsEstateUseCase) GetDronePlanDistance(ctx context.Context, estateId string, params generated.GetDistanceForDronePlanParams) (generated.GetDronePlanDistance, error) {
	panic("not expected to be called")
}

func (m *mockStatsEstateUseCase) GetStatsEstate(ctx context.Context, estateId string) (generated.GetStatsEstateResponse, error) {
	panic("not expected to be called")
}

func (m *mockStatsEstateUseCase) PublishCalculationStatsEstate(estateId string) {
	// no-op
}

// mockStatsEstateRepositoryError returns an error on QueryById.
type mockStatsEstateRepositoryError struct{}

func (m mockStatsEstateRepositoryError) QueryById(ctx context.Context, estateId string) (*model.EstateStats, error) {
	return nil, errors.New("connection timeout")
}

func (m mockStatsEstateRepositoryError) QueryAllTree(ctx context.Context, estateId string) ([]model.Tree, error) {
	panic("not expected to be called")
}

func (m mockStatsEstateRepositoryError) SaveStatsEstate(ctx context.Context, estateStats model.EstateStats) error {
	panic("not expected to be called")
}

// mockStatsEstateQueryByIdNil returns no stats and no error.
type mockStatsEstateQueryByIdNil struct{}

func (m mockStatsEstateQueryByIdNil) QueryById(ctx context.Context, estateId string) (*model.EstateStats, error) {
	return nil, nil
}

func (m mockStatsEstateQueryByIdNil) QueryAllTree(ctx context.Context, estateId string) ([]model.Tree, error) {
	panic("not expected to be called")
}

func (m mockStatsEstateQueryByIdNil) SaveStatsEstate(ctx context.Context, estateStats model.EstateStats) error {
	panic("not expected to be called")
}
