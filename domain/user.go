package domain

import "github.com/google/uuid"

type User struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Username  string    `json:"username,omitempty" validate:"required,min=5,max=20"`
	FirstName string    `json:"first_name,omitempty" validate:"required,max=20,alpha"`
	LastName  string    `json:"last_name,omitempty" validate:"required,max=20,alpha"`
	Password  string    `json:"password,omitempty" validate:"required,max=20"`
}
