package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBlockRepository_QueryByPlot_QueryFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	estateId := uuid.New().String()
	mock.ExpectQuery("FROM block").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("query failed"))

	repo := NewBlockRepository(db)
	block, err := repo.QueryByEstateIdAndBlockCoordinate(context.Background(), estateId, 5, 10)

	assert.Nil(t, block)
	assert.EqualError(t, err, "query failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlockRepository_QueryByPlot_QueryNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	estateId := uuid.New().String()
	mock.ExpectQuery("FROM block").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "estate_id", "x_coordinate", "y_coordinate"}))

	repo := NewBlockRepository(db)
	block, err := repo.QueryByEstateIdAndBlockCoordinate(context.Background(), estateId, 5, 10)
	assert.Nil(t, block)
	assert.Nil(t, err)
}

func TestBlockRepository_QueryByEstateId_QueryExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	estateId := uuid.New().String()
	blockId := uuid.New().String()
	mock.ExpectQuery("FROM block").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "estate_id", "x_coordinate", "y_coordinate", "ext_info", "created_at", "updated_at"}).AddRow(blockId, estateId, 5, 10, nil, time.Now(), time.Now()),
		) // empty result s

	repo := NewBlockRepository(db)
	block, err := repo.QueryByEstateIdAndBlockCoordinate(context.Background(), estateId, 5, 10)

	assert.Nil(t, err)
	assert.NotNil(t, block)
	assert.Equal(t, blockId, block.Id)
}
