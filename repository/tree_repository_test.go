package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTreeRepository_SaveBlockAndTreeSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	estateId := uuid.New().String()
	treeId := uuid.New().String()
	mock.ExpectQuery("WITH new_block").
		WithArgs(sqlmock.AnyArg(), estateId, 20, 50, sqlmock.AnyArg(), 20).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(treeId))

	repo := NewTreeRepository(db)
	id, err := repo.SaveBlockAndTree(context.Background(), estateId, 20, 50, 20)

	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.NotEmpty(t, *id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTreeRepository_SaveBlockAndTreeFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	estateId := uuid.New().String()
	mock.ExpectQuery("WITH new_block").
		WithArgs(sqlmock.AnyArg(), estateId, 20, 50, sqlmock.AnyArg(), 20).
		WillReturnError(errors.New("insert failed"))

	repo := NewTreeRepository(db)
	id, err := repo.SaveBlockAndTree(context.Background(), estateId, 20, 50, 20)

	assert.Nil(t, id)
	assert.EqualError(t, err, "insert failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTreeRepository_SaveBlockAndTreeUniqueViolation(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	estateId := uuid.New().String()
	pqErr := &pq.Error{Code: "23505"}
	mock.ExpectQuery("WITH new_block").
		WithArgs(sqlmock.AnyArg(), estateId, 20, 50, sqlmock.AnyArg(), 20).
		WillReturnError(pqErr)

	repo := NewTreeRepository(db)
	id, err := repo.SaveBlockAndTree(context.Background(), estateId, 20, 50, 20)

	assert.Nil(t, id)
	assert.ErrorIs(t, err, ErrBlockHasTree)
	assert.NoError(t, mock.ExpectationsWereMet())
}
