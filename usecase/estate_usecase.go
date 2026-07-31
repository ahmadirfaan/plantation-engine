package usecase

import (
	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/repository"
)

type EstateUseCase interface {
	CreateEstate(request generated.CreateEstateRequest) (*string, error)
}

func NewEstateUseCase(estateRepo repository.EstateRepository) EstateUseCase {
	return &estateUseCase{
		estateRepository: estateRepo,
	}
}

type estateUseCase struct {
	estateRepository repository.EstateRepository
}

func (e estateUseCase) CreateEstate(request generated.CreateEstateRequest) (*string, error) {

	id, err := e.estateRepository.SaveEstate(request.Width, request.Length)
	if err != nil {
		return nil, err
	}

	return id, nil
}
