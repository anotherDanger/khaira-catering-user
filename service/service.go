package service

import (
	"context"
	"khaira-catering-user/domain"
)

type Service interface {
	GetProducts(ctx context.Context) ([]*domain.Products, error)
	Login(ctx context.Context, username string, password string) (*domain.User, error)
	Register(ctx context.Context, entity *domain.User) (*domain.User, error)
	AddToCart(ctx context.Context, username string, product *domain.Products, quantity int) error
	GetCart(ctx context.Context, username string) ([]*domain.CartItem, error)
	DeleteCartItem(ctx context.Context, username string, productID string) error
	DeleteCartItemByQuantity(ctx context.Context, username string, productId string, quantity int) error
	CreateOrder(ctx context.Context, orderDetails *domain.Checkout) error
	GetOrderHistory(ctx context.Context, username string) ([]*domain.Checkout, error)
}
