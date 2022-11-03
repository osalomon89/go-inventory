package domain

import "time"

// swagger:model Book
type Book struct {
	// The ID of a book
	// example: 12345
	ID uint `json:"id"`
	// The Author of a book
	// example: Some name
	Author    string    `json:"author" validate:"required"`
	Title     string    `json:"title"`
	Price     int       `json:"price"`
	Isbn      string    `json:"isbn"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
