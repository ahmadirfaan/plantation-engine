package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestQueryAllTree_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := statsEstateRepository{DB: db}
	estateID := uuid.New().String()

	treeIdOne := uuid.New().String()
	treeIdTwo := uuid.New().String()
	mockRows := sqlmock.NewRows([]string{
		"id", "estate_id", "block_id", "x_coordinate", "y_coordinate", "height",
	}).
		AddRow(treeIdOne, estateID, uuid.New().String(), 10, 20, 5).
		AddRow(treeIdTwo, estateID, uuid.New().String(), 11, 21, 6)

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT t.id, t.estate_id, t.block_id, b.x_coordinate, b.y_coordinate, t.height
         FROM tree t
         JOIN block b ON b.id = t.block_id
         WHERE t.estate_id = $1`)).
		WithArgs(estateID).
		WillReturnRows(mockRows)

	result, err := repo.QueryAllTree(context.Background(), estateID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	assert.Equal(t, treeIdOne, result[0].Id)
	assert.Equal(t, 10, result[0].XAxis)
	assert.Equal(t, treeIdTwo, result[1].Id)
	assert.Equal(t, 21, result[1].YAxis)

	// Ensure expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryAllTree_Fail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := statsEstateRepository{DB: db}

	estateID := uuid.New().String()

	mock.ExpectQuery("SELECT .* FROM tree .*").
		WithArgs(estateID).
		WillReturnError(errors.New("connection timeout"))

	result, err := repo.QueryAllTree(context.Background(), estateID)
	assert.Nil(t, result)
	assert.ErrorContains(t, err, "timeout")
}

func TestQueryAllTree_FailWhenScan(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := statsEstateRepository{DB: db}

	estateID := uuid.New().String()
	mockRows := sqlmock.NewRows([]string{
		"id", "estate_id", "block_id", "x_coordinate", "y_coordinate", "height",
	}).AddRow(uuid.New().String(), estateID, uuid.New().String(), "x_axis_not_integer", 20, 5) // Invalid x_coordinate

	mock.ExpectQuery("SELECT .* FROM tree .*").
		WithArgs(estateID).
		WillReturnRows(mockRows)

	result, err := repo.QueryAllTree(context.Background(), estateID)
	assert.Nil(t, result)
	assert.ErrorContains(t, err, "Scan error")
}

func TestSaveStatsEstate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO estate_stats").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewStatsEstateRepository(db)
	err = repo.SaveStatsEstate(context.Background(), model.EstateStats{EstateId: uuid.New().String(), SumTree: 5, MinHeightTree: 5, MaxHeightTree: 5, MedianHeightTree: 15.6, TotalDistanceDrone: 200})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveStatsEstate_InsertFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO estate_stats").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert failed"))

	repo := NewStatsEstateRepository(db)
	err = repo.SaveStatsEstate(context.Background(), model.EstateStats{EstateId: uuid.New().String(), SumTree: 5, MinHeightTree: 5, MaxHeightTree: 5, MedianHeightTree: 15.6, TotalDistanceDrone: 200})

	assert.EqualError(t, err, "insert failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveStatsEstate_QueryById_Failed(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM estate_stats .*").
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("queryById failed"))

	repo := NewStatsEstateRepository(db)
	_, err = repo.QueryById(context.Background(), uuid.New().String())

	assert.EqualError(t, err, "queryById failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveStatsEstate_QueryById_ErrNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM estate_stats .*").
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(sql.ErrNoRows)

	repo := NewStatsEstateRepository(db)
	statsEstate, err := repo.QueryById(context.Background(), uuid.New().String())

	assert.Nil(t, statsEstate)
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveStatsEstate_QueryById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	estateId := uuid.New().String()

	mock.ExpectQuery("SELECT .* FROM estate_stats .*").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows(
			[]string{"estate_id", "min_height_tree", "max_height_tree", "median_height_tree", "sum_tree", "total_distance_drone", "ext_info", "created_at", "updated_at"}).
			AddRow(estateId, 3, 5, 4, 3, 54, nil, time.Now(), time.Now()))

	repo := NewStatsEstateRepository(db)
	statsEstate, err := repo.QueryById(context.Background(), estateId)

	assert.Nil(t, err)
	assert.NotNil(t, statsEstate)
	assert.Equal(t, estateId, statsEstate.EstateId)
	assert.Equal(t, 4.0, statsEstate.MedianHeightTree)
	assert.Equal(t, 54, statsEstate.TotalDistanceDrone)
	assert.NoError(t, mock.ExpectationsWereMet())
}
