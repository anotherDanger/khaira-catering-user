package domain

import (
	"time"

	"github.com/google/uuid"
)

type Checkout struct {
	Id          uuid.UUID  `json:"id"`
	Nama        string     `json:"nama"`
	Username    string     `json:"username"`
	ProductId   string     `json:"product_id"`
	ProductName string     `json:"product_name"`
	Quantity    int        `json:"quantity"`
	Total       float64    `json:"total"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	NoHp        string     `json:"no_hp"`
	Alamat      string     `json:"alamat"`
	Kecamatan   string     `json:"kecamatan"`
	Desa        string     `json:"desa"`
}
