package usecase

import (
	"errors"
	"time"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/google/uuid"
)

type mockEstateRepositoryError struct{}

func (m *mockEstateRepositoryError) QueryByEstateId(id string) (*model.Estate, error) {
	return nil, errors.New("error Connection Timeout")
}

func (m *mockEstateRepositoryError) SaveEstate(width int, length int) (*string, error) {
	return nil, errors.New("mock create estate error")
}

type mockEstateRepositoryQueryIdNotExist struct{}

func (m *mockEstateRepositoryQueryIdNotExist) QueryByEstateId(id string) (*model.Estate, error) {
	return nil, nil
}

func (m *mockEstateRepositoryQueryIdNotExist) SaveEstate(width int, length int) (*string, error) {
	return nil, nil
}

type mockEstateRepositoryQueryIdEstateExist struct{}

func (m *mockEstateRepositoryQueryIdEstateExist) QueryByEstateId(id string) (*model.Estate, error) {
	timeNow := time.Now()
	return &model.Estate{
		Id:        id,
		Width:     3,
		Length:    3,
		CreatedAt: &timeNow,
		UpdatedAt: &timeNow,
	}, nil
}

func (m *mockEstateRepositoryQueryIdEstateExist) SaveEstate(width int, length int) (*string, error) {
	return nil, nil
}

type mockEstateRepositorySuccess struct{}

func (m *mockEstateRepositorySuccess) QueryByEstateId(id string) (*model.Estate, error) {
	return nil, nil
}

func (m *mockEstateRepositorySuccess) SaveEstate(width int, length int) (*string, error) {
	id := uuid.New().String()
	return &id, nil
}

type mockTreeRepositoryError struct {
}

func (m *mockTreeRepositoryError) SaveTree(blockId string, estateId string, height int) (*string, error) {
	return nil, errors.New("connection timeout")
}

type mockTreeRepositorySaveTreeSuccess struct {
}

func (m *mockTreeRepositorySaveTreeSuccess) SaveTree(blockId string, estateId string, height int) (*string, error) {
	treeId := uuid.New().String()
	return &treeId, nil
}

type mockBlockRepositoryError struct {
}

func (m *mockBlockRepositoryError) QueryByEstateIdAndBlockCoordinate(id string, x int, y int) (*model.Block, error) {
	return nil, errors.New("connection timeout")
}

func (m *mockBlockRepositoryError) SaveBlock(estateId string, x int, y int) (*string, error) {
	return nil, errors.New("connection timeout")
}

type mockBlockRepositoryBlockExist struct {
}

func (m *mockBlockRepositoryBlockExist) QueryByEstateIdAndBlockCoordinate(id string, x int, y int) (*model.Block, error) {
	now := time.Now()
	return &model.Block{
		EstateId:  id,
		BlockX:    x,
		BlockY:    y,
		CreatedAt: &now,
		UpdatedAt: &now,
		Id:        uuid.New().String(),
	}, nil
}

func (m *mockBlockRepositoryBlockExist) SaveBlock(estateId string, x int, y int) (*string, error) {
	blockId := uuid.New().String()
	return &blockId, nil
}

type mockBlockRepositorySavingBlockError struct {
}

func (m *mockBlockRepositorySavingBlockError) QueryByEstateIdAndBlockCoordinate(id string, x int, y int) (*model.Block, error) {
	return nil, nil
}

func (m *mockBlockRepositorySavingBlockError) SaveBlock(estateId string, x int, y int) (*string, error) {
	return nil, errors.New("connection timeout")
}

type mockBlockRepositorySuccess struct {
}

func (m *mockBlockRepositorySuccess) QueryByEstateIdAndBlockCoordinate(id string, x int, y int) (*model.Block, error) {
	return nil, nil
}

func (m *mockBlockRepositorySuccess) SaveBlock(estateId string, x int, y int) (*string, error) {
	blockId := uuid.New().String()
	return &blockId, nil
}

type mockStatsEstateUseCase struct {
}

func (s *mockStatsEstateUseCase) GetDronePlanDistance(estateId string, params generated.GetDistanceForDronePlanParams) (generated.GetDronePlanDistance, error) {
	//TODO implement me
	panic("implement me")
}

func (s *mockStatsEstateUseCase) GetStatsEstate(estateId string) (generated.GetStatsEstateResponse, error) {
	panic("implement me")
}

func (s *mockStatsEstateUseCase) PublishCalculationStatsEstate(estateId string) {

}

type mockStatsEstateRepositoryError struct {
}

func (m mockStatsEstateRepositoryError) QueryAllTree(estateId string) ([]model.Tree, error) {
	panic("implement me")
}

func (m mockStatsEstateRepositoryError) SaveStatsEstate(estateStats model.EstateStats) error {
	panic("implement me")
}

func (m mockStatsEstateRepositoryError) QueryById(estateId string) (*model.EstateStats, error) {
	return nil, errors.New("connection timeout")
}

type mockStatsEstateQueryByIdNil struct {
}

func (m mockStatsEstateQueryByIdNil) QueryAllTree(estateId string) ([]model.Tree, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockStatsEstateQueryByIdNil) SaveStatsEstate(estateStats model.EstateStats) error {
	//TODO implement me
	panic("implement me")
}

func (m mockStatsEstateQueryByIdNil) QueryById(estateId string) (*model.EstateStats, error) {
	return nil, nil
}
