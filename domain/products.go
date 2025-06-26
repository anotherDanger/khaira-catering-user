package domain

import "time"

type Products struct {
	Id            string     `json:"product_id" validate:"required"`
	Name          string     `json:"product_name" validate:"required,min=5,max=50"`
	Description   string     `json:"description" validate:"alphanum"`
	Stock         int        `json:"product_stock" validate:"required,number"`
	Price         int        `json:"product_price" validate:"required,number"`
	ImageMetadata string     `json:"image_metadata" validate:"max=255"`
	CreatedAt     *time.Time `json:"created_at" validate:"required"`
	ModifiedAt    *time.Time `json:"modified_at" validate:"required"`
}
