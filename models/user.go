package models

import "time"

type RegisterUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Email       string     `json:"email" db:"email"`
	Role        string     `json:"role" db:"role"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ArchivedAt  *time.Time `json:"archived_at" db:"archived_at"`
	SuspendedAt *time.Time `json:"suspended_at" db:"suspended_at"`
}
