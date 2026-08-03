package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockEstateCalculation struct{}

func (m mockEstateCalculation) QueryById(ctx context.Context, estateId string) (*model.EstateStats, error) {
	panic("implement me")
}

func (m mockEstateCalculation) QueryAllTree(ctx context.Context, estateId string) ([]model.Tree, error) {
	panic("implement me")
}

func (m mockEstateCalculation) SaveStatsEstate(ctx context.Context, estateStats model.EstateStats) error {
	panic("implement me")
}

type mockEstateCalculationQueryAllTreeError struct{}

func (m mockEstateCalculationQueryAllTreeError) QueryById(ctx context.Context, estateId string) (*model.EstateStats, error) {
	panic("implement me")
}

func (m mockEstateCalculationQueryAllTreeError) QueryAllTree(ctx context.Context, estateId string) ([]model.Tree, error) {
	return nil, errors.New("mock query error")
}

func (m mockEstateCalculationQueryAllTreeError) SaveStatsEstate(ctx context.Context, estateStats model.EstateStats) error {
	panic("implement me")
}

type mockEstateCalculationWhereSavingSummaryError struct{}

func (m mockEstateCalculationWhereSavingSummaryError) QueryById(ctx context.Context, estateId string) (*model.EstateStats, error) {
	panic("implement me")
}

func (m mockEstateCalculationWhereSavingSummaryError) QueryAllTree(ctx context.Context, estateId string) ([]model.Tree, error) {
	trees := make([]model.Tree, 0)
	trees = append(trees, model.Tree{
		XAxis:  2,
		YAxis:  1,
		Height: 5,
	})
	trees = append(trees, model.Tree{
		XAxis:  3,
		YAxis:  1,
		Height: 3,
	})
	trees = append(trees, model.Tree{
		XAxis:  4,
		YAxis:  1,
		Height: 4,
	})
	return trees, nil
}

func (m mockEstateCalculationWhereSavingSummaryError) SaveStatsEstate(ctx context.Context, estateStats model.EstateStats) error {
	return errors.New("save query error")
}

type mockStatsEstateCalculationSuccess struct{}

func (m mockStatsEstateCalculationSuccess) QueryById(ctx context.Context, estateId string) (*model.EstateStats, error) {
	return &model.EstateStats{
		MinHeightTree:      3,
		MaxHeightTree:      5,
		MedianHeightTree:   4.0,
		SumTree:            3,
		TotalDistanceDrone: 54,
		EstateId:           estateId,
	}, nil
}

func (m mockStatsEstateCalculationSuccess) QueryAllTree(ctx context.Context, estateId string) ([]model.Tree, error) {
	trees := make([]model.Tree, 0)
	trees = append(trees, model.Tree{
		XAxis:  2,
		YAxis:  1,
		Height: 5,
	})
	trees = append(trees, model.Tree{
		XAxis:  3,
		YAxis:  1,
		Height: 3,
	})
	trees = append(trees, model.Tree{
		XAxis:  4,
		YAxis:  1,
		Height: 4,
	})
	trees = append(trees, model.Tree{
		XAxis:  1,
		YAxis:  1,
		Height: 3,
	})
	return trees, nil
}

func (m mockStatsEstateCalculationSuccess) SaveStatsEstate(ctx context.Context, estateStats model.EstateStats) error {
	return errors.New("save query error")
}

func TestCalculateStatsSummary_NegativeTestCase(t *testing.T) {
	ctx := context.Background()

	//stop process worker when estate not exist
	useCase := NewStatsEstateUseCase(mockEstateCalculation{}, &mockEstateRepositoryQueryIdNotExist{})
	useCase.PublishCalculationStatsEstate(uuid.New().String())

	//query all tree error timeout
	useCase = NewStatsEstateUseCase(mockEstateCalculationQueryAllTreeError{}, &mockEstateRepositoryQueryIdEstateExist{})
	useCase.PublishCalculationStatsEstate(uuid.New().String())

	//save stats estate error
	useCase = NewStatsEstateUseCase(mockEstateCalculationWhereSavingSummaryError{}, &mockEstateRepositoryQueryIdEstateExist{})
	useCase.PublishCalculationStatsEstate(uuid.New().String())

	_ = ctx
}

func TestCalculateStatsSummary_PositiveCase(t *testing.T) {
	ctx := context.Background()

	//success calculate summary
	useCase := NewStatsEstateUseCase(mockStatsEstateCalculationSuccess{}, &mockEstateRepositoryQueryIdEstateExist{})
	useCase.PublishCalculationStatsEstate(uuid.New().String())

	_ = ctx
}

func TestCalculateStatsSummary_GetSummary_Success(t *testing.T) {
	ctx := context.Background()

	//success get summary
	useCase := NewStatsEstateUseCase(mockStatsEstateCalculationSuccess{}, &mockEstateRepositoryQueryIdEstateExist{})
	estateId := uuid.New().String()
	response, err := useCase.GetStatsEstate(ctx, estateId)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 3, response.Count)
	assert.Equal(t, 4.0, float64(response.Median))

	//get summary when stats is not exist then will return 0 for all data
	useCase = NewStatsEstateUseCase(mockStatsEstateQueryByIdNil{}, &mockEstateRepositoryQueryIdEstateExist{})
	estateId = uuid.New().String()
	response, err = useCase.GetStatsEstate(ctx, estateId)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, response.Count)
	assert.Equal(t, 0.0, float64(response.Median))
	assert.Equal(t, 0, response.Max)
	assert.Equal(t, 0, response.Min)

}

func TestCalculateStatsSummary_GetSummary_Error(t *testing.T) {
	ctx := context.Background()

	//query estate by id is get error
	useCase := NewStatsEstateUseCase(mockStatsEstateCalculationSuccess{}, &mockEstateRepositoryError{})
	estateId := uuid.New().String()
	_, err := useCase.GetStatsEstate(ctx, estateId)

	assert.NotNil(t, err)

	//query estate by id is not exist
	useCase = NewStatsEstateUseCase(mockStatsEstateQueryByIdNil{}, &mockEstateRepositoryQueryIdNotExist{})
	estateId = uuid.New().String()
	_, err = useCase.GetStatsEstate(ctx, estateId)

	assert.NotNil(t, err)
	assert.Equal(t, "404|estate not found", err.Error())

	//query stats by estate get error
	useCase = NewStatsEstateUseCase(mockStatsEstateRepositoryError{}, &mockEstateRepositoryQueryIdEstateExist{})
	estateId = uuid.New().String()
	_, err = useCase.GetStatsEstate(ctx, estateId)

	assert.NotNil(t, err)
	assert.Equal(t, "connection timeout", err.Error())

}

func TestStatsEstateUseCase_GetDronePlanDistance_Error(t *testing.T) {
	ctx := context.Background()

	//query estate by id is get error
	useCase := NewStatsEstateUseCase(mockStatsEstateCalculationSuccess{}, &mockEstateRepositoryError{})
	estateId := uuid.New().String()
	_, err := useCase.GetDronePlanDistance(ctx, estateId, generated.GetDistanceForDronePlanParams{})

	assert.NotNil(t, err)

}

func TestGetDronePlanDistance_Success(t *testing.T) {
	ctx := context.Background()

	//success get summary
	useCase := NewStatsEstateUseCase(mockStatsEstateCalculationSuccess{}, &mockEstateRepositoryQueryIdEstateExist{})
	estateId := uuid.New().String()
	response, err := useCase.GetDronePlanDistance(ctx, estateId, generated.GetDistanceForDronePlanParams{})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 54, response.Distance)

	//get drone plan when stats is not exist then will return 0
	useCase = NewStatsEstateUseCase(mockStatsEstateQueryByIdNil{}, &mockEstateRepositoryQueryIdEstateExist{})
	estateId = uuid.New().String()
	response, err = useCase.GetDronePlanDistance(ctx, estateId, generated.GetDistanceForDronePlanParams{})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, response.Distance)

}

func TestGetDronePlanDistance_WithMaxDistance(t *testing.T) {
	ctx := context.Background()
	estateId := uuid.New().String()

	//max distance beyond total -> rest is last block finish, no tree query
	useCase := NewStatsEstateUseCase(mockStatsEstateCalculationSuccess{}, &mockEstateRepositoryQueryIdEstateExist{})
	maxDistance := 100
	response, err := useCase.GetDronePlanDistance(ctx, estateId, generated.GetDistanceForDronePlanParams{MaxDistance: &maxDistance})

	assert.Nil(t, err)
	assert.NotNil(t, response.Rest)
	assert.Equal(t, 54, response.Distance)

	//max distance within total -> falls back to all-trees query + cache
	useCase = NewStatsEstateUseCase(mockStatsEstateCalculationSuccess{}, &mockEstateRepositoryQueryIdEstateExist{})
	maxDistance = 27
	response, err = useCase.GetDronePlanDistance(ctx, estateId, generated.GetDistanceForDronePlanParams{MaxDistance: &maxDistance})

	assert.Nil(t, err)
	assert.NotNil(t, response.Rest)
	assert.Equal(t, 54, response.Distance)
}
