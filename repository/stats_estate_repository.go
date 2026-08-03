package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/ahmadirfaan/plantation-engine/model"
)

type StatsEstateRepository interface {
	QueryAllTree(ctx context.Context, estateId string) ([]model.Tree, error)
	SaveStatsEstate(ctx context.Context, estateStats model.EstateStats) error
	QueryById(ctx context.Context, estateId string) (*model.EstateStats, error)
}

func NewStatsEstateRepository(db *sql.DB) StatsEstateRepository {
	return &statsEstateRepository{
		DB: db,
	}
}

type statsEstateRepository struct {
	DB *sql.DB
}

func (s *statsEstateRepository) QueryById(ctx context.Context, estateId string) (*model.EstateStats, error) {
	var statsEstate model.EstateStats
	err := s.DB.QueryRowContext(ctx, "SELECT estate_id, min_height_tree, max_height_tree, median_height_tree, sum_tree, total_distance_drone, ext_info,created_at, updated_at  FROM estate_stats WHERE estate_id = $1", estateId).
		Scan(&statsEstate.EstateId, &statsEstate.MinHeightTree, &statsEstate.MaxHeightTree,
			&statsEstate.MedianHeightTree, &statsEstate.SumTree, &statsEstate.TotalDistanceDrone, &statsEstate.ExtInfo, &statsEstate.CreatedAt, &statsEstate.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		slog.Error("failed to query stats by estateId", "estateId", estateId, "error", err)
		return nil, err
	}
	return &statsEstate, nil
}

func (s *statsEstateRepository) SaveStatsEstate(ctx context.Context, estateStats model.EstateStats) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO estate_stats (
    estate_id,
    sum_tree,
    min_height_tree,
    max_height_tree,
    median_height_tree,
    total_distance_drone,
    ext_info
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (estate_id)
DO UPDATE
SET
    sum_tree = EXCLUDED.sum_tree,
    min_height_tree = EXCLUDED.min_height_tree,
    max_height_tree = EXCLUDED.max_height_tree,
    median_height_tree = EXCLUDED.median_height_tree,
    total_distance_drone = EXCLUDED.total_distance_drone,
    ext_info = EXCLUDED.ext_info,
    updated_at = NOW();`, estateStats.EstateId, estateStats.SumTree, estateStats.MinHeightTree, estateStats.MaxHeightTree, estateStats.MedianHeightTree, estateStats.TotalDistanceDrone, estateStats.ExtInfo)
	return err
}

func (s *statsEstateRepository) QueryAllTree(ctx context.Context, estateId string) ([]model.Tree, error) {
	var trees []model.Tree

	rows, err := s.DB.QueryContext(
		ctx,
		`SELECT t.id, t.estate_id, t.block_id, b.x_coordinate, b.y_coordinate, t.height
         FROM tree t
         JOIN block b ON b.id = t.block_id
         WHERE t.estate_id = $1`, estateId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.Tree
		if err := rows.Scan(&t.Id, &t.EstateId, &t.BlockId, &t.XAxis, &t.YAxis, &t.Height); err != nil {
			return nil, err
		}
		trees = append(trees, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trees, nil
}
