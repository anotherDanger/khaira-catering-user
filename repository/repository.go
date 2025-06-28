package repository

import (
	"context"
	"database/sql"
	"khaira-catering-user/domain"

	"github.com/google/uuid"
)

type Repository interface {
	GetProducts(ctx context.Context, db *sql.DB) ([]*domain.Products, error)
	Login(ctx context.Context, db *sql.DB, username string, password string) (*domain.User, error)
	Register(ctx context.Context, db *sql.DB, entity *domain.User) (*domain.User, error)
	AddToCart(ctx context.Context, username string, product *domain.Products, quantity int, sql *sql.DB) error
	GetCart(ctx context.Context, username string) ([]*domain.CartItem, error)
	DeleteCartItem(ctx context.Context, username string, productID string) error
	DeleteCartItemByQuantity(ctx context.Context, username, productId string, quantity int) error
	CreateOrder(ctx context.Context, tx *sql.Tx, orderDetails *domain.Checkout, id uuid.UUID) error
}
