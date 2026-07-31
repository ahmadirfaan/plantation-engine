package usecase

import (
	"strings"
	"testing"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func buildTreeRequest(x int, y int, height int) generated.CreateTreeRequest {
	return generated.CreateTreeRequest{
		X:      x,
		Y:      y,
		Height: height,
	}
}

func assertErrorMessage(t *testing.T, err error, treeId *string, expectedErrorMessage string) {
	assert.NotNil(t, err)
	assert.Nil(t, treeId)
	assert.Equal(t, expectedErrorMessage, err.Error())
}

func TestAddTreeQueryEstateErrorOrNotExist(t *testing.T) {

	//error when query from DB
	newTreeUseCase := NewTreeUseCase(&mockEstateRepositoryError{}, &mockTreeRepositoryError{}, &mockBlockRepositoryError{}, &mockStatsEstateUseCase{})
	id := uuid.New().String()
	treeId, err := newTreeUseCase.AddTreeToEstate(id, buildTreeRequest(5, 1, 25))
	assertErrorMessage(t, err, treeId, "error Connection Timeout")

	//error when estate not exist
	newTreeUseCase = NewTreeUseCase(&mockEstateRepositoryQueryIdNotExist{}, &mockTreeRepositoryError{}, &mockBlockRepositoryError{}, &mockStatsEstateUseCase{})
	treeId, err = newTreeUseCase.AddTreeToEstate(id, buildTreeRequest(5, 1, 25))
	assertErrorMessage(t, err, treeId, "404|estate does not exist")

	//error when add tree is above the width estate
	newTreeUseCase = NewTreeUseCase(&mockEstateRepositoryQueryIdEstateExist{}, &mockTreeRepositoryError{}, &mockBlockRepositoryError{}, &mockStatsEstateUseCase{})
	treeId, err = newTreeUseCase.AddTreeToEstate(id, buildTreeRequest(5, 1, 25))
	assertErrorMessage(t, err, treeId, "400|X request above from data")

	//error when add tree is above the length estate
	newTreeUseCase = NewTreeUseCase(&mockEstateRepositoryQueryIdEstateExist{}, &mockTreeRepositoryError{}, &mockBlockRepositoryError{}, &mockStatsEstateUseCase{})
	treeId, err = newTreeUseCase.AddTreeToEstate(id, buildTreeRequest(3, 4, 25))
	assertErrorMessage(t, err, treeId, "400|Y request above from data")

}

func TestAddTreeQueryBlockErrorOrAlreadyExist(t *testing.T) {

	//connection timeout when query by id block
	newTreeUseCase := NewTreeUseCase(&mockEstateRepositoryQueryIdEstateExist{}, &mockTreeRepositoryError{}, &mockBlockRepositoryError{}, &mockStatsEstateUseCase{})
	estateId := uuid.New().String()

	treeId, err := newTreeUseCase.AddTreeToEstate(estateId, buildTreeRequest(3, 3, 25))
	assertErrorMessage(t, err, treeId, "connection timeout")

	//block already has tree
	newTreeUseCase = NewTreeUseCase(&mockEstateRepositoryQueryIdEstateExist{}, &mockTreeRepositoryError{}, &mockBlockRepositoryBlockExist{}, &mockStatsEstateUseCase{})
	treeId, err = newTreeUseCase.AddTreeToEstate(estateId, buildTreeRequest(3, 3, 25))
	assertErrorMessage(t, err, treeId, "409|block already has a tree")
}

func TestAddTreeSavingBlockError(t *testing.T) {
	//connection timeout when saving block
	newTreeUseCase := NewTreeUseCase(&mockEstateRepositoryQueryIdEstateExist{}, &mockTreeRepositoryError{}, &mockBlockRepositorySavingBlockError{}, &mockStatsEstateUseCase{})
	estateId := uuid.New().String()

	treeId, err := newTreeUseCase.AddTreeToEstate(estateId, buildTreeRequest(3, 3, 25))
	assertErrorMessage(t, err, treeId, "connection timeout")
}

func TestAddTreeSavingTreeError(t *testing.T) {
	//connection timeout when saving tree
	newTreeUseCase := NewTreeUseCase(&mockEstateRepositoryQueryIdEstateExist{}, &mockTreeRepositoryError{}, &mockBlockRepositorySuccess{}, &mockStatsEstateUseCase{})
	estateId := uuid.New().String()

	treeId, err := newTreeUseCase.AddTreeToEstate(estateId, buildTreeRequest(3, 3, 25))
	assertErrorMessage(t, err, treeId, "connection timeout")
}

func TestAddTreeSavingTreeSuccess(t *testing.T) {
	//connection timeout when query by id block
	newTreeUseCase := NewTreeUseCase(&mockEstateRepositoryQueryIdEstateExist{}, &mockTreeRepositorySaveTreeSuccess{}, &mockBlockRepositorySuccess{}, &mockStatsEstateUseCase{})
	estateId := uuid.New().String()

	treeId, err := newTreeUseCase.AddTreeToEstate(estateId, buildTreeRequest(3, 3, 25))
	assert.Nil(t, err)
	assert.NotNil(t, treeId)
	assert.True(t, strings.TrimSpace(*treeId) != "", "tree id should not be empty")
}
