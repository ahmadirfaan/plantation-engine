package usecase

import (
	"context"
	"testing"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/stretchr/testify/assert"
)

func Test_CreateEstateError_WhenInsertingToDB(t *testing.T) {
	ctx := context.Background()

	estateUseCase := NewEstateUseCase(&mockEstateRepositoryError{})
	_, err := estateUseCase.CreateEstate(ctx, generated.CreateEstateRequest{
		Width:  50,
		Length: 10,
	})
	assert.NotNil(t, err)

}

func Test_CreateEstateSuccess_WhenInsertingToDB(t *testing.T) {
	ctx := context.Background()

	estateUseCase := NewEstateUseCase(&mockEstateRepositorySuccess{})
	id, err := estateUseCase.CreateEstate(ctx, generated.CreateEstateRequest{
		Width:  50,
		Length: 10,
	})
	assert.Nil(t, err)
	assert.NotNil(t, id)
}
