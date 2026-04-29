package models

import "time"

type Todos struct {
	Id          string    `json:"id" db:"id"`
	UserId      string    `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name" binding:"required"`
	Description string    `json:"description" db:"description" binding:"required"`
	CompleteAt  string    `json:"completed_at" db:"completed_at"`
	ExpiryAt    time.Time `json:"expiry_at" db:"expiry_at" binding:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateTodo struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	ExpiryAt    time.Time `json:"expiry_at" binding:"required"`
}

type UpdateTodo struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Complete    *bool      `json:"complete"`
	ExpiryAt    *time.Time `json:"expiry_at"`
}
