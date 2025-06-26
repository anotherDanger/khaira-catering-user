package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"khaira-catering-user/domain"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
)

type RepositoryImpl struct {
	elastic *elasticsearch.Client
}

func NewRepositoryImpl(elastic *elasticsearch.Client) Repository {
	return &RepositoryImpl{
		elastic: elastic,
	}
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
		Password: user.Password,
	}

	return row, nil
}

func (repo *RepositoryImpl) Register(ctx context.Context, db *sql.DB, entity *domain.User) (*domain.User, error) {
	query := "INSERT INTO users(id, username, first_name, last_name, password) VALUES(?, ?, ?, ?, ?)"
	result, err := db.ExecContext(ctx, query, entity.Id, entity.Username, entity.FirstName, entity.LastName, entity.Password)
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

	doc := map[string]interface{}{
		"id":       entity.Id,
		"username": entity.Username,
	}

	body, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	docID := fmt.Sprint(entity.Id)

	res, err := repo.elastic.Index(
		"user_cart",
		strings.NewReader(string(body)),
		repo.elastic.Index.WithDocumentID(docID),
		repo.elastic.Index.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.Status())
	}

	user := &domain.User{
		Id:       entity.Id,
		Username: entity.Username,
	}

	return user, nil
}
