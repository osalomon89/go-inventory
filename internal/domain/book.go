package domain

import "time"

type Book struct {
	ID        int       `json:"id"`
	Author    string    `json:"author" validate:"required"`
	Title     string    `json:"title"`
	Price     int       `json:"price"`
	Isbn      string    `json:"isbn"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
