package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/ahmadirfaan/plantation-engine/model"
)

type BlockRepository interface {
	QueryByEstateIdAndBlockCoordinate(ctx context.Context, estateId string, x int, y int) (*model.Block, error)
}

func NewBlockRepository(db *sql.DB) BlockRepository {
	return &blockRepository{
		DB: db,
	}
}

type blockRepository struct {
	DB *sql.DB
}

func (b *blockRepository) QueryByEstateIdAndBlockCoordinate(ctx context.Context, estateId string, x int, y int) (*model.Block, error) {
	var block model.Block
	err := b.DB.QueryRowContext(ctx, "SELECT id, estate_id, x_coordinate, y_coordinate, ext_info, created_at, updated_at FROM block WHERE estate_id = $1 AND x_coordinate = $2 AND y_coordinate = $3", estateId, x, y).Scan(&block.Id,
		&block.EstateId,
		&block.BlockX,
		&block.BlockY,
		&block.ExtInfo,
		&block.CreatedAt,
		&block.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		slog.Error("failed to query block by coordinate", "estateId", estateId, "x", x, "y", y, "error", err)
		return nil, err
	}
	return &block, nil
}
