package domain

import (
	"time"

	"github.com/google/uuid"
)

type Checkout struct {
	Id          uuid.UUID  `json:"id"`
	Nama        string     `json:"nama" validate:"required,max=100"`
	Username    string     `json:"username" validate:"required,max=20"`
	ProductId   string     `json:"product_id" validate:"required"`
	ProductName string     `json:"product_name" validate:"required"`
	Quantity    int        `json:"quantity" validate:"required"`
	Total       float64    `json:"total" validate:"required"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	NoHp        string     `json:"no_hp" validate:"required,max=12"`
	Alamat      string     `json:"alamat" validate:"required"`
	Kecamatan   string     `json:"kecamatan" validate:"required"`
	Desa        string     `json:"desa" validate:"required"`
}
