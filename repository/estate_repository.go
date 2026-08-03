package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/google/uuid"
)

type EstateRepository interface {
	SaveEstate(ctx context.Context, width int, length int) (*string, error)
	QueryByEstateId(ctx context.Context, id string) (*model.Estate, error)
}

func NewEstateRepository(db *sql.DB) EstateRepository {
	return &estateRepository{
		DB: db,
	}
}

type estateRepository struct {
	DB *sql.DB
}

func (e *estateRepository) SaveEstate(ctx context.Context, width int, length int) (*string, error) {
	uuidString := uuid.New().String()
	_, err := e.DB.ExecContext(ctx, "INSERT INTO estate (id, width, length) VALUES ($1, $2, $3)", uuidString, width, length)
	if err != nil {
		slog.Error("failed to insert estate", "error", err)
		return nil, err
	}
	return &uuidString, nil
}

func (e *estateRepository) QueryByEstateId(ctx context.Context, id string) (*model.Estate, error) {
	var estate model.Estate
	err := e.DB.QueryRowContext(ctx, "SELECT id, name, width, length, ext_info, created_at, updated_at FROM estate WHERE id = $1", id).Scan(&estate.Id,
		&estate.Name,
		&estate.Width,
		&estate.Length,
		&estate.ExtInfo,
		&estate.CreatedAt,
		&estate.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		slog.Error("failed to query estate by id", "estateId", id, "error", err)
		return nil, err
	}
	return &estate, nil
}
