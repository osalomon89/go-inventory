package repositories

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/osalomon89/go-inventory/internal/domain"
)

type BookRepository interface {
	GetBooksById(id uint) (*domain.Book, error)
	CreateBook(*domain.Book) error
	GetBooks() ([]domain.Book, error)
	UpdateBook(int, *domain.Book) error
}

type bookRepository struct {
	conn *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) BookRepository {
	return &bookRepository{
		conn: db,
	}
}

func (repo *bookRepository) GetBooksById(id uint) (*domain.Book, error) {
	book := new(domain.Book)
	err := repo.conn.Get(book, "SELECT * FROM books WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting book: %w", &err)
	}
	return book, nil
}
func (repo *bookRepository) GetBooks() ([]domain.Book, error) {
	var libros []domain.Book
	err := repo.conn.Select(libros, "SELECT * FROM books")
	if err != nil {
		return nil, fmt.Errorf("error gettings books: %w", &err)
	}
	return libros, nil

	return libros, nil
}
func (repo *bookRepository) CreateBook(book *domain.Book) error {
	createdAt := time.Now()
	result, err := repo.conn.Exec(`INSERT INTO books 
	(title, author, price, stock, isbn, created_at, updated_at) 
	VALUES(?,?,?,?,?,?, ?)`, book.Title, book.Author, book.Price, book.Stock, book.Isbn, createdAt, createdAt)

	if err != nil {
		return fmt.Errorf("error inserting book: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error saving book: %w", err)
	}
	book.ID = int(id)
	return nil
}

func (repo *bookRepository) UpdateBook(id int, book *domain.Book) error {
	return nil
}
