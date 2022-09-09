package repositories

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/osalomon89/go-inventory/internal/domain"
)

type BookRepository interface {
	GetBooksById(id uint) (*domain.Book, error)
	CreateBook(*domain.Book) error
	GetBooks() ([]domain.Book, error)
	GetBooksByParams(params map[string]interface{}) ([]domain.Book, error)
	UpdateBook(*domain.Book) error
	UpdateBookByParams(map[string]interface{}, *domain.Book) error
	DeleteBook(int) error
}

type bookRepository struct {
	conn *sqlx.DB
}

var mapaCampos = map[string]string{"author": "", "title": "", "price": "", "isbn": "", "stock": ""}

func NewBookRepository(db *sqlx.DB) BookRepository {
	return &bookRepository{
		conn: db,
	}
}

func (repo *bookRepository) GetBooksById(id uint) (*domain.Book, error) {
	book := new(domain.Book)
	err := repo.conn.Get(book, "SELECT * FROM books WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting book: %w", err)
	}
	return book, nil
}
func (repo *bookRepository) GetBooks() ([]domain.Book, error) {
	libros := []domain.Book{}
	err := repo.conn.Select(&libros, "SELECT * FROM books")
	if err != nil {
		return nil, fmt.Errorf("error gettings books: %w", err)
	}
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
func (repo *bookRepository) GetBooksByParams(params map[string]interface{}) ([]domain.Book, error) {
	libros := []domain.Book{}
	setParamsSlice, setValues := repo.getStatementParams(params)
	setParams := strings.Join(setParamsSlice, " AND ")
	query := fmt.Sprintf("SELECT * FROM books WHERE %s", setParams)
	err := repo.conn.Select(&libros, query, setValues...)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los libros: %w", err)
	}
	return libros, nil

}
func (repo *bookRepository) UpdateBook(book *domain.Book) error {
	updateAt := time.Now()
	_, err := repo.conn.Exec(`UPDATE books SET title=?, author=?, price=?, stock=?,update_at=? WHERE id=?`, book.Title, book.Author, book.Price, book.Stock, updateAt, book.ID)
	if err != nil {
		return fmt.Errorf("Error al actualizar el libro: %w", err)
	}
	return nil
}

func (repo *bookRepository) UpdateBookByParams(params map[string]interface{}, book *domain.Book) error {
	updateAt := time.Now()
	params["update_at"] = updateAt
	marshal, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}
	err = json.Unmarshal(marshal, book)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	setParamsSlice, setValues := repo.getStatementParams(params)
	setParams := strings.Join(setParamsSlice, ",")

	query := fmt.Sprintf("UPDATE books SET %s WHERE id=?", setParams)
	setValues = append(setValues, book.ID)

	_, err = repo.conn.Exec(query, setValues...)

	if err != nil {
		return fmt.Errorf("error updating item: %w", err)
	}

	return nil
}

func (repo *bookRepository) getStatementParams(params map[string]interface{}) ([]string, []interface{}) {
	var setParams []string
	var setValues []interface{}

	for key, val := range params {
		_, ok := mapaCampos[key]
		if !ok || val == nil || val == "" || val == 0 {
			continue
		}
		setParams = append(setParams, fmt.Sprintf("%s=?", key))
		setValues = append(setValues, val)
	}

	return setParams, setValues
}

func (repo *bookRepository) DeleteBook(id int) error {
	_, err := repo.conn.Exec(`DELETE FROM books WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("Erro al eliminar el libro: %w", err)
	}
	return nil
}
