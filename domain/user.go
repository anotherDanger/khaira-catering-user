package domain

import "github.com/google/uuid"

type User struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password,omitempty"`
}
