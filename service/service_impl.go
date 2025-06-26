package service

import (
	"context"
	"database/sql"
	"khaira-catering-user/domain"
	"khaira-catering-user/repository"
)

type ServiceImpl struct {
	db   *sql.DB
	repo repository.Repository
}

func NewServiceImpl(db *sql.DB, repo repository.Repository) Service {
	return &ServiceImpl{
		db:   db,
		repo: repo,
	}
}

func (svc *ServiceImpl) GetProducts(ctx context.Context) ([]*domain.Products, error) {
	result, err := svc.repo.GetProducts(ctx, svc.db)
	if err != nil {
		return nil, err
	}

	return result, nil
}
