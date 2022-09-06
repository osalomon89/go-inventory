package repositories

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/osalomon89/go-inventory/internal/domain"
)

type BookRepository interface {
	GetBookByID(id int) (*domain.Book, error)
	CreateBook(book *domain.Book) error
	GetBook() ([]domain.Book, error)
	UpdateBook(id int, book *domain.Book) error
	DeleteBook(id int) error
}

type bookRepository struct {
	conn *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) BookRepository {
	return &bookRepository{
		conn: db,
	}
}

func (repo *bookRepository) GetBookByID(id int) (*domain.Book, error) {
	book := new(domain.Book)
	err := repo.conn.Get(book, "SELECT * FROM books WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting books: %w", err)
	}
	return book, nil
}

func (repo *bookRepository) CreateBook(book *domain.Book) error {
	createdAt := time.Now()

	result, err := repo.conn.Exec(`INSERT INTO books 
		(title, author, price, stock, isbn, created_at, updated_at) 
		VALUES(?,?,?,?,?,?,?)`, book.Title, book.Author, book.Price, book.Stock, book.Isbn, createdAt, createdAt)

	if err != nil {
		return fmt.Errorf("error inserting book: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error saving book: %w", err)
	}

	book.Id = int(id)

	return nil
}

func (repo *bookRepository) GetBook() ([]domain.Book, error) {
	books := []domain.Book{}

	err := repo.conn.Select(&books, "SELECT * FROM books")
	if err != nil {
		return nil, fmt.Errorf("error getting books: %w", err)
	}
	return books, nil

}

func (repo *bookRepository) UpdateBook(id int, book *domain.Book) error {
	updatedAt := time.Now()

	result, err := repo.conn.Exec(`UPDATE books SET title=?, author=?, price=?, stock=?, isbn=?, 
	updated_at=? WHERE id=?`, book.Title, book.Author, book.Price, book.Stock, book.Isbn,
		updatedAt, id)

	if err != nil {
		return fmt.Errorf("error insering book: %w", err)
	}

	fmt.Println(result)

	return nil
}

func (repo *bookRepository) DeleteBook(id int) error {
	result, err := repo.conn.Exec(`DELETE FROM books WHERE id=?`, id)

	if err != nil {
		return fmt.Errorf("error deleting book: %w", err)
	}

	fmt.Println(result)

	return nil
}
