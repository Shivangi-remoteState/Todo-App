package models

import "time"

type Todos struct {
	Id          string    `json:"id" db:"id"`
	UserId      string    `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description" validate:"required"`
	Complete    string    `json:"complete" db:"complete"`
	ExpiryAt    time.Time `json:"expiry_at" db:"expiry_at" validate:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateTodo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ExpiryAt    time.Time `json:"expiry_at"`
}

type UpdateTodo struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Complete    *bool      `json:"complete"`
	ExpiryAt    *time.Time `json:"expiry_at"`
}
