package service

import (
	"context"
	"database/sql"
	"khaira-catering-user/domain"
	"khaira-catering-user/helper"
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
		return nil, helper.NewError(err.Error())
	}
	return result, nil
}

func (svc *ServiceImpl) Login(ctx context.Context, username string, password string) (*domain.User, error) {
	result, err := svc.repo.Login(ctx, svc.db, username, password)
	if err != nil {
		return nil, helper.NewError(err.Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
	if err != nil {
		return nil, helper.NewError("password tidak cocok")
	}
	return result, nil
}

func (svc *ServiceImpl) Register(ctx context.Context, entity *domain.User) (*domain.User, error) {
	userId := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, helper.NewError("gagal mengenkripsi password")
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
		return nil, helper.NewError(err.Error())
	}
	return result, nil
}

func (svc *ServiceImpl) AddToCart(ctx context.Context, username string, product *domain.Products, quantity int) error {
	err := svc.repo.AddToCart(ctx, username, product, quantity, svc.db)
	if err != nil {
		return helper.NewError(err.Error())
	}
	return nil
}

func (svc *ServiceImpl) GetCart(ctx context.Context, username string) ([]*domain.CartItem, error) {
	cart, err := svc.repo.GetCart(ctx, username)
	if err != nil {
		return nil, helper.NewError(err.Error())
	}
	return cart, nil
}

func (svc *ServiceImpl) DeleteCartItem(ctx context.Context, username string, productID string) error {
	err := svc.repo.DeleteCartItem(ctx, username, productID)
	if err != nil {
		return helper.NewError(err.Error())
	}
	return nil
}

func (svc *ServiceImpl) DeleteCartItemByQuantity(ctx context.Context, username string, productId string, quantity int) error {
	err := svc.repo.DeleteCartItemByQuantity(ctx, username, productId, quantity)
	if err != nil {
		return helper.NewError(err.Error())
	}
	return nil
}

func (svc *ServiceImpl) CreateOrder(ctx context.Context, orderDetails *domain.Checkout) error {
	tx, err := svc.db.Begin()
	if err != nil {
		return helper.NewError("gagal memulai transaksi")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	id := uuid.New()
	err = svc.repo.CreateOrder(ctx, tx, orderDetails, id)
	if err != nil {
		return helper.NewError(err.Error())
	}
	return nil
}

func (svc *ServiceImpl) GetOrderHistory(ctx context.Context, username string) ([]*domain.Checkout, error) {
	res, err := svc.repo.GetOrderHistory(ctx, svc.db, username)
	if err != nil {
		return nil, helper.NewError(err.Error())
	}
	return res, nil
}
