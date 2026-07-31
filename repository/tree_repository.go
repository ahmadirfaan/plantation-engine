package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"log"
)

type TreeRepository interface {
	SaveTree(blockId string, estateId string, height int) (*string, error)
}

func NewTreeRepository(db *sql.DB) TreeRepository {
	return &treeRepository{
		DB: db,
	}
}

type treeRepository struct {
	DB *sql.DB
}

func (t *treeRepository) SaveTree(blockId string, estateId string, height int) (*string, error) {
	id := uuid.New().String()
	_, err := t.DB.ExecContext(context.Background(), "INSERT INTO tree (id, block_id, estate_id, height) VALUES ($1, $2, $3, $4)", id, blockId, estateId, height)
	if err != nil {
		log.Println("Failed to insert tree into database")
		log.Println(err.Error())
		return nil, err
	}
	return &id, nil
}
