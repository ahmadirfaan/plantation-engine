package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ErrBlockHasTree is returned when a block already has a tree and the insert
// violates the unique constraint on (estate_id, x_coordinate, y_coordinate).
var ErrBlockHasTree = errors.New("block already has a tree")

type TreeRepository interface {
	// SaveBlockAndTree atomically inserts the block and its tree in a single
	// statement so a failed tree insert cannot leave an orphan block.
	SaveBlockAndTree(ctx context.Context, estateId string, x int, y int, height int) (*string, error)
}

func NewTreeRepository(db *sql.DB) TreeRepository {
	return &treeRepository{
		DB: db,
	}
}

type treeRepository struct {
	DB *sql.DB
}

func (t *treeRepository) SaveBlockAndTree(ctx context.Context, estateId string, x int, y int, height int) (*string, error) {
	blockId := uuid.New().String()
	treeId := uuid.New().String()
	err := t.DB.QueryRowContext(ctx, `
WITH new_block AS (
    INSERT INTO block (id, estate_id, x_coordinate, y_coordinate)
    VALUES ($1, $2, $3, $4)
    RETURNING id
)
INSERT INTO tree (id, block_id, estate_id, height)
VALUES ($5, (SELECT id FROM new_block), $2, $6)
RETURNING id`, blockId, estateId, x, y, treeId, height).Scan(&treeId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, ErrBlockHasTree
		}
		slog.Error("failed to insert block and tree", "estateId", estateId, "x", x, "y", y, "error", err)
		return nil, err
	}
	return &treeId, nil
}
