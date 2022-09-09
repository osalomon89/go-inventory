package repositories

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/osalomon89/go-inventory/internal/domain"
)

type BookRepository interface {
	GetBookByID(id int) (*domain.Book, error)
	CreateBook(book *domain.Book) error
	GetBook(params map[string]interface{}) ([]domain.Book, error)
	GetBookByAuthor(author string) ([]domain.Book, error)
	UpdateBook(id int, book *domain.Book) error
	UpdateBookByParams(id int, params map[string]interface{}, book *domain.Book) error
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

func (repo *bookRepository) GetBook(params map[string]interface{}) ([]domain.Book, error) {
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

func (repo *bookRepository) GetBookByAuthor(author string) ([]domain.Book, error) {
	books := []domain.Book{}

	/*if author == "" {
		err := repo.conn.Select(&books, "SELECT * FROM books")
		if err != nil {
			return nil, fmt.Errorf("error getting books: %w", err)
		}
	} else {*/
	err := repo.conn.Select(&books, "SELECT * FROM books WHERE author=? ORDER BY title ASC LIMIT 5 ", author)
	if err != nil {
		return nil, fmt.Errorf("error getting books: %w", err)
	}
	//}

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

func (repo *bookRepository) UpdateBookByParams(id int, params map[string]interface{}, book *domain.Book) error {
	updatedAt := time.Now()

	params["updated_at"] = updatedAt
	//paso los datos a bytes
	ja, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("error marshalling book: %w", err)
	}
	//los convierte en struct
	err = json.Unmarshal(ja, book)
	if err != nil {
		return fmt.Errorf("error unmarshalling book: %w", err)
	}

	setParams, setValues := repo.getStatementParams(params)

	query := fmt.Sprintf("UPDATE books SET %s WHERE id=?", setParams)
	setValues = append(setValues, id)

	result, err := repo.conn.Exec(query, setValues...)

	if err != nil {
		return fmt.Errorf("error updating item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error updating item: %w%d", err, rowsAffected)
	}

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
