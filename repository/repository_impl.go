package repository

import (
	"context"
	"database/sql"
	"khaira-catering-user/domain"
)

type RepositoryImpl struct{}

func NewRepositoryImpl() Repository {
	return &RepositoryImpl{}
}

func (repo *RepositoryImpl) GetProducts(ctx context.Context, db *sql.DB) ([]*domain.Products, error) {
	query := "SELECT * FROM products"

	result, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var products []*domain.Products
	for result.Next() {
		var product domain.Products
		var description sql.NullString
		var imageMetadata sql.NullString
		if err := result.Scan(&product.Id, &product.Name, &description, &product.Stock, &product.Price, &imageMetadata, &product.CreatedAt, &product.ModifiedAt); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (repo *RepositoryImpl) Login(ctx context.Context, db *sql.DB, username string, password string) (*domain.User, error) {
	query := "SELECT username, password from users WHERE username = ?"

	result := db.QueryRowContext(ctx, query, username)
	var user domain.User

	if err := result.Scan(&user.Username, &user.Password); err != nil {
		return nil, err
	}

	row := &domain.User{
		Username: user.Username,
	}

	if user.Password == password {
		return row, nil
	}

	return nil, sql.ErrNoRows
}

func (repo *RepositoryImpl) Register(ctx context.Context, db *sql.DB, entity *domain.User) (*domain.User, error) {
	query := "INSERT INTO users(id, username, first_name, last_name, password) VALUES(?, ?, ?, ?, ?)"
	result, err := db.ExecContext(ctx, query, entity.Username, entity.Id, entity.FirstName, entity.LastName, entity.Password)
	if err != nil {
		return nil, err
	}

	rowAff, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowAff == 0 {
		return nil, sql.ErrNoRows
	}

	user := &domain.User{
		Id:       entity.Id,
		Username: entity.Username,
	}

	return user, nil
}
