package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/google/uuid"
)

type EstateRepository interface {
	SaveEstate(width int, length int) (*string, error)
	QueryByEstateId(id string) (*model.Estate, error)
}

func NewEstateRepository(db *sql.DB) EstateRepository {
	return &estateRepository{
		DB: db,
	}
}

type estateRepository struct {
	DB *sql.DB
}

func (e *estateRepository) SaveEstate(width int, length int) (*string, error) {
	uuidString := uuid.New().String()
	_, err := e.DB.ExecContext(context.Background(), "INSERT INTO estate (id, width, length) VALUES ($1, $2, $3)", uuidString, width, length)
	if err != nil {
		log.Println("Failed to insert into database")
		log.Println(err.Error())
		return nil, err
	}
	return &uuidString, nil
}

func (e *estateRepository) QueryByEstateId(id string) (*model.Estate, error) {
	var estate model.Estate
	err := e.DB.QueryRowContext(context.Background(), "SELECT id, name, width, length, ext_info, created_at, updated_at FROM estate WHERE id = $1", id).Scan(&estate.Id,
		&estate.Name,
		&estate.Width,
		&estate.Length,
		&estate.ExtInfo,
		&estate.CreatedAt,
		&estate.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		log.Println("Failed to query estate by id")
		log.Println(err.Error())
		return nil, err
	}
	return &estate, nil
}
