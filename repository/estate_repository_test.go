package repository

import (
	"errors"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveEstate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO estate").
		WithArgs(sqlmock.AnyArg(), 50, 100).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewEstateRepository(db)
	id, err := repo.SaveEstate(50, 100)

	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.NotEmpty(t, *id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveEstate_InsertFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO estate").
		WithArgs(sqlmock.AnyArg(), 50, 100).
		WillReturnError(errors.New("insert failed"))

	repo := NewEstateRepository(db)
	id, err := repo.SaveEstate(50, 100)

	assert.Nil(t, id)
	assert.EqualError(t, err, "insert failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEstateRepository_QueryByEstateId_QueryFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	estateId := uuid.New().String()
	mock.ExpectQuery("FROM estate").
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("query failed"))

	repo := NewEstateRepository(db)
	estate, err := repo.QueryByEstateId(estateId)

	assert.Nil(t, estate)
	assert.EqualError(t, err, "query failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEstateRepository_QueryByEstateId_QueryNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	estateId := uuid.New().String()
	mock.ExpectQuery("FROM estate").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"})) // empty result s

	repo := NewEstateRepository(db)
	estate, err := repo.QueryByEstateId(estateId)

	assert.Nil(t, estate)
	assert.Nil(t, err)
}

func TestEstateRepository_QueryByEstateId_QueryExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	estateId := uuid.New().String()
	mock.ExpectQuery("FROM estate").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "width", "length", "ext_info", "created_at", "updated_at"}).AddRow(estateId, nil, 50, 50, nil, time.Now(), time.Now()),
		) // empty result s

	repo := NewEstateRepository(db)
	estate, err := repo.QueryByEstateId(estateId)

	assert.Nil(t, err)
	assert.NotNil(t, estate)
	assert.Equal(t, estateId, estate.Id)
}
