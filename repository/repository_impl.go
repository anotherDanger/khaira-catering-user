package repository

import (
	"bytes"
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
		if err := result.Scan(&product.Id, &product.Name, &description, &product.Price, &product.Stock, &imageMetadata, &product.CreatedAt, &product.ModifiedAt); err != nil {
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

	docID := fmt.Sprint(entity.Username)

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

func (repo *RepositoryImpl) AddToCart(ctx context.Context, username string, product *domain.Products) error {
	newCartItem := map[string]interface{}{
		"product_id":   product.Id,
		"product_name": product.Name,
		"quantity":     product.Stock,
		"price":        product.Price,
	}

	script := `
		if (ctx._source.cart == null) {
			ctx._source.cart = [params.product];
		} else {
			def item = ctx._source.cart.find(p -> p.product_id == params.product.product_id);
			if (item != null) {
				item.quantity += params.product.quantity;
			} else {
				ctx._source.cart.add(params.product);
			}
		}
	`

	updateBody, err := json.Marshal(map[string]interface{}{
		"script": map[string]interface{}{
			"lang":   "painless",
			"source": script,
			"params": map[string]interface{}{
				"product": newCartItem,
			},
		},
		"upsert": map[string]interface{}{
			"username": username,
			"cart":     []interface{}{newCartItem},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to marshal update request: %w", err)
	}

	res, err := repo.elastic.Update(
		"user_cart",
		username,
		bytes.NewReader(updateBody),
		repo.elastic.Update.WithContext(ctx),
		repo.elastic.Update.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("elasticsearch update error: %s", res.Status())
		}
		return fmt.Errorf("elasticsearch update error: %s: %v", res.Status(), e)
	}

	return nil
}

func (repo *RepositoryImpl) GetCart(ctx context.Context, username string) ([]*domain.CartItem, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"username": username,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	res, err := repo.elastic.Search(
		repo.elastic.Search.WithContext(ctx),
		repo.elastic.Search.WithIndex("user_cart"),
		repo.elastic.Search.WithBody(&buf),
		repo.elastic.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch search error: %s", res.Status())
	}

	var esRes domain.ESResponse
	if err := json.NewDecoder(res.Body).Decode(&esRes); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(esRes.Hits.Hits) == 0 {
		return nil, nil
	}

	var result []*domain.CartItem
	for i := range esRes.Hits.Hits[0].Source.Cart {
		result = append(result, &esRes.Hits.Hits[0].Source.Cart[i])
	}

	return result, nil
}
