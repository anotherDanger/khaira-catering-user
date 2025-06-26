package repository

import (
	"context"
	"database/sql"
	"khaira-catering-user/domain"
)

type Repository interface {
	GetProducts(ctx context.Context, db *sql.DB) ([]*domain.Products, error)
}
