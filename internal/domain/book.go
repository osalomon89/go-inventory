package domain

import "time"

// swagger:model Book
type Book struct {
	// The ID of a book
	// example: 12345
	ID uint `json:"id"`
	// The Author of a book
	// example: Some name
	Author string `json:"author" validate:"required"`
	// The Title of a book
	// example: The tittle book
	Title string `json:"title"`
	// The Price of a book
	// example: 5000
	Price int `json:"price"`
	// The ISBN of a book
	// example: 123errtr5789
	Isbn string `json:"isbn"`
	// The number of books
	// example: 50
	Stock int `json:"stock"`
	// The date of creation
	// example: 2022-10-01T22:11:33
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// The date of update
	// example: 2022-10-01T22:11:33
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
