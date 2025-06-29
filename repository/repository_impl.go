package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"khaira-catering-user/domain"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/google/uuid"
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
		return nil, errors.New("gagal mengambil data produk")
	}

	var products []*domain.Products
	for result.Next() {
		var product domain.Products
		var description sql.NullString
		var imageMetadata sql.NullString
		if err := result.Scan(&product.Id, &product.Name, &description, &product.Price, &product.Stock, &imageMetadata, &product.CreatedAt, &product.ModifiedAt); err != nil {
			return nil, errors.New("gagal membaca data produk")
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
		return nil, errors.New("user tidak ditemukan")
	}

	return &user, nil
}

func (repo *RepositoryImpl) Register(ctx context.Context, db *sql.DB, entity *domain.User) (*domain.User, error) {
	query := "INSERT INTO users(id, username, first_name, last_name, password) VALUES(?, ?, ?, ?, ?)"
	result, err := db.ExecContext(ctx, query, entity.Id, entity.Username, entity.FirstName, entity.LastName, entity.Password)
	if err != nil {
		return nil, errors.New("registrasi gagal")
	}

	rowAff, err := result.RowsAffected()
	if err != nil {
		return nil, errors.New("registrasi gagal")
	}

	if rowAff == 0 {
		return nil, errors.New("registrasi gagal")
	}

	doc := map[string]interface{}{
		"id":       entity.Id,
		"username": entity.Username,
	}

	body, err := json.Marshal(doc)
	if err != nil {
		return nil, errors.New("gagal menyiapkan data pengguna")
	}

	docID := fmt.Sprint(entity.Username)

	res, err := repo.elastic.Index(
		"user_cart",
		strings.NewReader(string(body)),
		repo.elastic.Index.WithDocumentID(docID),
		repo.elastic.Index.WithContext(ctx),
	)
	if err != nil {
		return nil, errors.New("gagal membuat cart pengguna")
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New("gagal menyimpan cart pengguna")
	}

	return &domain.User{
		Id:       entity.Id,
		Username: entity.Username,
	}, nil
}

func (repo *RepositoryImpl) AddToCart(ctx context.Context, username string, product *domain.Products, quantity int, db *sql.DB) error {
	var productStock int
	err := db.QueryRowContext(ctx, "SELECT stock FROM products WHERE id = ?", product.Id).Scan(&productStock)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	getRes, err := repo.elastic.Get("user_cart", username)
	if err != nil {
		return errors.New("gagal mengambil data cart")
	}
	defer getRes.Body.Close()

	if getRes.StatusCode != 404 {
		if getRes.IsError() {
			return errors.New("gagal mengambil data cart")
		}
		var cartData struct {
			Source struct {
				Cart []struct {
					ProductID string `json:"product_id"`
					Quantity  int    `json:"quantity"`
				} `json:"cart"`
			} `json:"_source"`
		}
		if err := json.NewDecoder(getRes.Body).Decode(&cartData); err != nil {
			return errors.New("gagal membaca data cart")
		}
		currentQty := 0
		for _, item := range cartData.Source.Cart {
			if item.ProductID == product.Id {
				currentQty = item.Quantity
				break
			}
		}
		if currentQty+quantity > productStock {
			return errors.New("melebihi stok produk")
		}
	}

	newCartItem := map[string]interface{}{
		"product_id":   product.Id,
		"product_name": product.Name,
		"quantity":     quantity,
		"price":        product.Price,
	}

	script := `
		if (ctx._source.cart == null) {
			ctx._source.cart = [params.product];
		} else {
			def item = ctx._source.cart.find(p -> p.product_id == params.product.product_id);
			if (item != null) {
				item.quantity += params.quantity;
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
				"product":  newCartItem,
				"quantity": quantity,
			},
		},
		"upsert": map[string]interface{}{
			"username": username,
			"cart":     []interface{}{newCartItem},
		},
	})
	if err != nil {
		return errors.New("gagal menyiapkan data pembaruan cart")
	}

	res, err := repo.elastic.Update(
		"user_cart",
		username,
		bytes.NewReader(updateBody),
		repo.elastic.Update.WithContext(ctx),
		repo.elastic.Update.WithRefresh("true"),
	)
	if err != nil {
		return errors.New("gagal memperbarui cart")
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New("gagal memperbarui cart")
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
		return nil, errors.New("gagal menyiapkan permintaan pencarian")
	}

	res, err := repo.elastic.Search(
		repo.elastic.Search.WithContext(ctx),
		repo.elastic.Search.WithIndex("user_cart"),
		repo.elastic.Search.WithBody(&buf),
		repo.elastic.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, errors.New("pencarian cart gagal")
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New("pencarian cart gagal")
	}

	var esRes domain.ESResponse
	if err := json.NewDecoder(res.Body).Decode(&esRes); err != nil {
		return nil, errors.New("gagal membaca hasil pencarian")
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

func (repo *RepositoryImpl) DeleteCartItem(ctx context.Context, username string, productID string) error {
	updateScript := map[string]interface{}{
		"script": map[string]interface{}{
			"source": `
				if (ctx._source.cart != null) {
					ctx._source.cart.removeIf(item -> item.product_id == params.product_id);
				}
			`,
			"lang": "painless",
			"params": map[string]interface{}{
				"product_id": productID,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(updateScript); err != nil {
		return errors.New("gagal menyiapkan data penghapusan")
	}

	res, err := repo.elastic.Update(
		"user_cart",
		username,
		&buf,
		repo.elastic.Update.WithContext(ctx),
		repo.elastic.Update.WithRefresh("true"),
	)
	if err != nil {
		return errors.New("gagal menghapus item dari cart")
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New("gagal menghapus item dari cart")
	}

	return nil
}

func (repo *RepositoryImpl) DeleteCartItemByQuantity(ctx context.Context, username, productId string, quantity int) error {
	res, err := repo.elastic.Get(
		"user_cart",
		username,
		repo.elastic.Get.WithContext(ctx),
	)
	if err != nil {
		return errors.New("gagal mengambil data cart")
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return errors.New("cart tidak ditemukan")
	}

	var data map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return errors.New("gagal membaca cart")
	}

	source := data["_source"].(map[string]interface{})
	cart := source["cart"].([]interface{})
	newCart := make([]interface{}, 0)

	for _, item := range cart {
		cartItem := item.(map[string]interface{})
		if cartItem["product_id"] == productId {
			qty := 0
			switch v := cartItem["quantity"].(type) {
			case float64:
				qty = int(v)
			case int:
				qty = v
			}
			remaining := qty - quantity
			if remaining > 0 {
				cartItem["quantity"] = remaining
				newCart = append(newCart, cartItem)
			}
		} else {
			newCart = append(newCart, cartItem)
		}
	}

	source["cart"] = newCart

	body, err := json.Marshal(source)
	if err != nil {
		return errors.New("gagal menyiapkan data cart terbaru")
	}

	indexRes, err := repo.elastic.Index(
		"user_cart",
		strings.NewReader(string(body)),
		repo.elastic.Index.WithDocumentID(username),
		repo.elastic.Index.WithContext(ctx),
		repo.elastic.Index.WithRefresh("true"),
	)
	if err != nil {
		return errors.New("gagal memperbarui cart")
	}
	defer indexRes.Body.Close()

	if indexRes.IsError() {
		return errors.New("gagal memperbarui cart")
	}

	return nil
}

func (repo *RepositoryImpl) CreateOrder(ctx context.Context, tx *sql.Tx, orderDetails *domain.Checkout, id uuid.UUID) error {
	orderId := id
	query := "INSERT INTO orders(id, product_id, product_name, nama, no_hp, alamat, kecamatan, desa, username, quantity, total) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query,
		orderId,
		orderDetails.ProductId,
		orderDetails.ProductName,
		orderDetails.Nama,
		orderDetails.NoHp,
		orderDetails.Alamat,
		orderDetails.Kecamatan,
		orderDetails.Desa,
		orderDetails.Total,
		orderDetails.Username,
		orderDetails.Quantity,
		orderDetails.Total,
	)
	if err != nil {
		return err
	}

	updateStockQuery := "UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?"
	result, err := tx.ExecContext(ctx, updateStockQuery, orderDetails.Quantity, orderDetails.ProductId, orderDetails.Quantity)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("stok tidak mencukupi untuk produk %s", orderDetails.ProductId)
	}

	return nil
}

func (repo *RepositoryImpl) GetOrderHistory(ctx context.Context, db *sql.DB, username string) ([]*domain.Checkout, error) {
	query := "SELECT id, product_id, product_name, nama, no_hp, alamat, kecamatan, desa, jumlah, username, quantity, created_at, status FROM orders WHERE username = ?"
	res, err := db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}

	var history []*domain.Checkout

	for res.Next() {
		var item domain.Checkout
		if err := res.Scan(&item.Id, &item.ProductId, &item.ProductName, &item.Nama, &item.NoHp, &item.Alamat, &item.Kecamatan, &item.Desa, &item.Total, &item.Username, &item.Quantity, &item.CreatedAt, &item.Status); err != nil {
			return nil, err
		}
		history = append(history, &item)
	}

	defer res.Close()

	return history, nil
}
