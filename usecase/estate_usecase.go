package usecase

import (
	"context"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/repository"
)

type EstateUseCase interface {
	CreateEstate(ctx context.Context, request generated.CreateEstateRequest) (*string, error)
}

func NewEstateUseCase(estateRepo repository.EstateRepository) EstateUseCase {
	return &estateUseCase{
		estateRepository: estateRepo,
	}
}

type estateUseCase struct {
	estateRepository repository.EstateRepository
}

func (e estateUseCase) CreateEstate(ctx context.Context, request generated.CreateEstateRequest) (*string, error) {
	return e.estateRepository.SaveEstate(ctx, request.Width, request.Length)
}
