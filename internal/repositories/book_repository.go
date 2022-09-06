package repositories

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/osalomon89/go-inventory/internal/domain"
)

type BookRepository interface {
	GetBookByID(id uint) (*domain.Book, error)
	CreateBook(book *domain.Book) error
	UpdateBookByParams(id uint, params map[string]interface{}, book *domain.Book) error
}

type bookRepository struct {
	conn *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) BookRepository {
	return &bookRepository{
		conn: db,
	}
}

func (repo *bookRepository) GetBookByID(id uint) (*domain.Book, error) {
	book := new(domain.Book)
	err := repo.conn.Get(book, "SELECT * FROM books WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting book: %w", err)
	}

	return book, nil
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

	book.ID = uint(id)

	return nil
}

func (repo *bookRepository) UpdateBookByParams(id uint,
	params map[string]interface{}, book *domain.Book) error {
	updateAt := time.Now()

	setParams, setValues := repo.getSetParams(params, updateAt)
	query := fmt.Sprintf("UPDATE books SET %s WHERE id=?", setParams)
	setValues = append(setValues, id)

	_, err := repo.conn.Exec(query, setValues...)

	if err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}

	book.UpdatedAt = updateAt

	return nil
}

func (repo *bookRepository) getSetParams(params map[string]interface{},
	updateAt time.Time) (string, []interface{}) {
	setParams := "updated_at=?"
	var setValues []interface{}
	setValues = append(setValues, updateAt)

	for key, val := range params {
		setParams = fmt.Sprintf("%s, %s=?", setParams, key)
		setValues = append(setValues, val)
	}

	return setParams, setValues
}
