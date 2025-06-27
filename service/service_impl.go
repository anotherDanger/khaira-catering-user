package service

import (
	"context"
	"database/sql"
	"khaira-catering-user/domain"
	"khaira-catering-user/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (svc *ServiceImpl) Login(ctx context.Context, username string, password string) (*domain.User, error) {
	result, err := svc.repo.Login(ctx, svc.db, username, password)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (svc *ServiceImpl) Register(ctx context.Context, entity *domain.User) (*domain.User, error) {
	userId := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Id:        userId,
		Username:  entity.Username,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Password:  string(hashedPassword),
	}

	result, err := svc.repo.Register(ctx, svc.db, user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *ServiceImpl) AddToCart(ctx context.Context, username string, product *domain.Products) error {
	err := svc.repo.AddToCart(ctx, username, product)
	if err != nil {
		return err
	}

	return nil
}

func (svc *ServiceImpl) GetCart(ctx context.Context, username string) ([]*domain.CartItem, error) {
	cart, err := svc.repo.GetCart(ctx, username)
	if err != nil {
		return nil, err
	}

	return cart, nil
}
