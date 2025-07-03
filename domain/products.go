package domain

import "time"

type Products struct {
	Id            string     `json:"product_id"`
	Name          string     `json:"product_name"`
	Description   string     `json:"description"`
	Stock         int        `json:"product_stock"`
	Price         int        `json:"product_price"`
	ImageMetadata string     `json:"image_metadata"`
	CreatedAt     *time.Time `json:"created_at"`
	ModifiedAt    *time.Time `json:"modified_at"`
}
