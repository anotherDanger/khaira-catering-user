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
