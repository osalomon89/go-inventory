package repositories

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/osalomon89/go-inventory/internal/domain"
)

type BookRepository interface {
	GetBooks(params map[string]interface{}) ([]domain.Book, error)
	GetBookByID(id uint) (*domain.Book, error)
	CreateBook(book *domain.Book) error
	UpdateBook(id uint, book *domain.Book) error
	DeleteBook(id uint) error
	UpdateBookByParams(id uint, book *domain.Book) error
}

type bookRepository struct {
	conn *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) BookRepository {
	return &bookRepository{
		conn: db,
	}
}

func (repo *bookRepository) GetBooks(params map[string]interface{}) ([]domain.Book, error) {
	books := []domain.Book{}
	limit := params["limit"]
	offset := params["offset"]
	delete(params, "limit")
	delete(params, "offset")

	whereQuery := ""
	setParamsSlice, setValues := repo.getStatementParams(params)
	if len(setParamsSlice) > 0 {
		setParams := strings.Join(setParamsSlice, " AND ")
		whereQuery = fmt.Sprintf(" WHERE %s", setParams)
	}

	queryLimit := repo.getLimitOffsetStatement(limit.(float64), offset.(float64))

	query := fmt.Sprintf("SELECT * FROM books%s %s", whereQuery, queryLimit)

	err := db.Select(&books, query, setValues...)
	if err != nil {
		return nil, fmt.Errorf("error getting all books: %w", err)
	}

	return books, nil
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
func (repo *bookRepository) UpdateBook(id uint, book *domain.Book) error {

	updatedAt := time.Now()

	result, err := repo.conn.Exec(`UPDATE books SET
		title=?, author=?, price=?, stock=?, isbn=?, updated_at=? WHERE id=? 
		`, book.Title, book.Author, book.Price, book.Stock, book.Isbn, updatedAt, id)

	if err != nil {
		return fmt.Errorf("error inserting book: %w", err)
	}
	fmt.Println(result)

	return nil
}
func (repo *bookRepository) DeleteBook(id uint) error {

	result, err := repo.conn.Exec(`DELETE FROM books WHERE id=? 
		`, id)

	if err != nil {
		return fmt.Errorf("error al eliminar el libro %w", err)
	}
	fmt.Println(result)

	return nil
}

func (repo *bookRepository) UpdateBookByParams(id uint, book *domain.Book) error {
	updateAt := time.Now()

	query, args := repo.patchMyResource(book)

	result, err := repo.conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error updating book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("the book does not exist: %w", err)
	}

	book.UpdatedAt = updateAt

	return nil
}

func (repo *bookRepository) patchMyResource(book *domain.Book) (string, []interface{}) {
	q := `UPDATE books SET `
	qParts := make([]string, 0, 2)
	args := make([]interface{}, 0, 2)
	updateAt := time.Now()

	qParts = append(qParts, `updated_at=?`)
	args = append(args, updateAt)

	if book.Author != "" {
		qParts = append(qParts, `author=?`)
		args = append(args, book.Author)
	}

	if book.Title != "" {
		qParts = append(qParts, `title = ?`)
		args = append(args, book.Title)
	}

	if book.Price != 0 {
		qParts = append(qParts, `price = ?`)
		args = append(args, book.Price)
	}

	if book.Isbn != "" {
		qParts = append(qParts, `isbn = ?`)
		args = append(args, book.Isbn)
	}

	if book.Stock != 0 {
		qParts = append(qParts, `stock = ?`)
		args = append(args, book.Stock)
	}

	q += strings.Join(qParts, ",") + ` WHERE id = ?`
	args = append(args, book.ID)

	book.UpdatedAt = updateAt

	return q, args
}

func (repo *bookRepository) getStatementParams(params map[string]interface{}) ([]string, []interface{}) {
	var setParams []string
	var setValues []interface{}

	for key, val := range params {
		if val == nil || val == "" || val == 0 {
			continue
		}
		setParams = append(setParams, fmt.Sprintf("%s=?", key))
		setValues = append(setValues, val)
	}

	return setParams, setValues
}

func (repo *bookRepository) getLimitOffsetStatement(limit, offset float64) string {
	queryLimit := 100
	queryOffset := 0
	if limit > 0 && limit < 1000 {
		queryLimit = int(limit)
	}

	if offset > 0 {
		queryOffset = int(offset)
	}

	return fmt.Sprintf("LIMIT %d OFFSET %d", queryLimit, queryOffset)
}
