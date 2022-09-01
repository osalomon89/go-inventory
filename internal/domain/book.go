package domain

type Book struct {
	ID     uint   `json:"id"`
	Author string `json:"author" validate:"required"`
	Title  string `json:"title"`
	Price  int    `json:"price"`
	Isbn   string `json:"isbn"`
	Stock  int    `json:"stock"`
}
