package service

import (
	"context"
	"khaira-catering-user/domain"
)

type Service interface {
	GetProducts(ctx context.Context) ([]*domain.Products, error)
	Login(ctx context.Context, username string, password string) (*domain.User, error)
}
