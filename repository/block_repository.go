package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/google/uuid"
)

type BlockRepository interface {
	QueryByEstateIdAndBlockCoordinate(estateId string, x int, y int) (*model.Block, error)
	SaveBlock(estateId string, x int, y int) (*string, error)
}

func NewBlockRepository(db *sql.DB) BlockRepository {
	return &blockRepository{
		DB: db,
	}
}

type blockRepository struct {
	DB *sql.DB
}

func (b *blockRepository) SaveBlock(estateId string, x int, y int) (*string, error) {
	id := uuid.New().String()
	_, err := b.DB.ExecContext(context.Background(), "INSERT INTO block (id, estate_id, x_coordinate, y_coordinate) VALUES ($1, $2, $3, $4)", id, estateId, x, y)
	if err != nil {
		log.Println("Failed to insert block into database")
		log.Println(err.Error())
		return nil, err
	}
	return &id, nil
}

func (b *blockRepository) QueryByEstateIdAndBlockCoordinate(estateId string, x int, y int) (*model.Block, error) {
	var block model.Block
	err := b.DB.QueryRowContext(context.Background(), "SELECT id, estate_id, x_coordinate, y_coordinate, ext_info, created_at, updated_at FROM block WHERE estate_id = $1 AND x_coordinate = $2 AND y_coordinate = $3", estateId, x, y).Scan(&block.Id,
		&block.EstateId,
		&block.BlockX,
		&block.BlockY,
		&block.ExtInfo,
		&block.CreatedAt,
		&block.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println("Failed to query block by id")
		log.Println(err.Error())
		return nil, err
	}
	return &block, nil
}
