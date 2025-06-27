package repository

import (
	"context"
	"database/sql"
	"khaira-catering-user/domain"
)

type Repository interface {
	GetProducts(ctx context.Context, db *sql.DB) ([]*domain.Products, error)
	Login(ctx context.Context, db *sql.DB, username string, password string) (*domain.User, error)
	Register(ctx context.Context, db *sql.DB, entity *domain.User) (*domain.User, error)
	AddToCart(ctx context.Context, username string, product *domain.Products) error
	GetCart(ctx context.Context, username string) ([]*domain.CartItem, error)
}
