package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTreeRepository_SaveTreeSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	blockId := uuid.New().String()
	estateId := uuid.New().String()
	mock.ExpectExec("INSERT INTO tree").
		WithArgs(sqlmock.AnyArg(), blockId, estateId, 20).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewTreeRepository(db)
	id, err := repo.SaveTree(blockId, estateId, 20)

	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.NotEmpty(t, *id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTreeRepository_SaveTreeFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	blockId := uuid.New().String()
	estateId := uuid.New().String()
	mock.ExpectExec("INSERT INTO tree").
		WithArgs(sqlmock.AnyArg(), blockId, estateId, 20).
		WillReturnError(errors.New("insert failed"))

	repo := NewTreeRepository(db)
	id, err := repo.SaveTree(blockId, estateId, 20)

	assert.Nil(t, id)
	assert.EqualError(t, err, "insert failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}
