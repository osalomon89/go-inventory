package repositories_test

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/osalomon89/go-inventory/internal/domain"
	"github.com/osalomon89/go-inventory/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func Test_bookRepository_GetBookByID(t *testing.T) {
	assert := assert.New(t)

	mockDB, mock, err := sqlmock.New()
	assert.Nil(err)
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	type behaviourDB struct {
		data *sqlmock.Rows
		err  error
	}

	type args struct {
		id uint
	}

	tests := []struct {
		name      string
		args      args
		dbArgs    behaviourDB
		want      *domain.Book
		wantedErr error
	}{
		{
			name: "Should return a book",
			args: args{id: 1},
			dbArgs: behaviourDB{
				data: sqlmock.
					NewRows([]string{"id", "author", "title", "price", "isbn", "stock"}).
					AddRow(1, "the-author", "the title", 5000, "abcd1234", 20),
				err: nil,
			},
			want: &domain.Book{
				ID:     1,
				Author: "the-author",
				Title:  "the title",
				Price:  5000,
				Isbn:   "abcd1234",
				Stock:  20,
			},
			wantedErr: nil,
		},
		{
			name: "Should return an error when sql query returns an error",
			args: args{id: 1},
			dbArgs: behaviourDB{
				data: sqlmock.NewRows([]string{"id"}).AddRow(1),
				err:  errors.New("sql query error"),
			},
			want:      nil,
			wantedErr: fmt.Errorf("error getting book: %w", errors.New("sql query error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM books WHERE id=?`)).
				WithArgs(tt.args.id).
				WillReturnRows(tt.dbArgs.data).
				WillReturnError(tt.dbArgs.err)

			repo := repositories.NewBookRepository(sqlxDB)
			got, err := repo.GetBookByID(tt.args.id)
			if tt.wantedErr != nil {
				assert.Equal(reflect.TypeOf(tt.wantedErr), reflect.TypeOf(err), "Error type is not the expected")
				assert.Equal(tt.wantedErr.Error(), err.Error(), "Error message is not the expected")
				return
			}
			assert.Equal(err, nil)
			assert.Equal(tt.want, got)
		})
	}
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func Test_bookRepository_CreateBook(t *testing.T) {
	assert := assert.New(t)

	mockDB, mock, err := sqlmock.New()
	assert.Nil(err)
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	type dbArgs struct {
		result driver.Result
		err    error
	}

	type args struct {
		book *domain.Book
	}
	tests := []struct {
		name      string
		dbArgs    dbArgs
		args      args
		wantedErr error
	}{
		{
			name: "Should save a book succesfully",
			args: args{
				book: &domain.Book{
					Author: "the-author",
					Title:  "the title",
					Price:  5000,
					Isbn:   "abcd1234",
					Stock:  20,
				},
			},
			dbArgs: dbArgs{
				result: sqlmock.NewResult(1, 1),
				err:    nil,
			},
			wantedErr: nil,
		},
		{
			name: "Should return an error when LastInsertId function returns an error",
			args: args{
				book: &domain.Book{
					Author: "the-author",
					Title:  "the title",
					Price:  5000,
					Isbn:   "abcd1234",
					Stock:  20,
				},
			},
			dbArgs: dbArgs{
				result: sqlmock.NewErrorResult(errors.New("sql error")),
				err:    nil,
			},
			wantedErr: fmt.Errorf("error getting last insert id: %w", errors.New("sql error")),
		},
		{
			name: "Should return an error when sql statement returns an error",
			args: args{
				book: &domain.Book{
					Author: "the-author",
					Title:  "the title",
					Price:  5000,
					Isbn:   "abcd1234",
					Stock:  20,
				},
			},
			dbArgs: dbArgs{
				result: sqlmock.NewErrorResult(errors.New("sql error")),
				err:    errors.New("sql error"),
			},
			wantedErr: fmt.Errorf("error inserting book: %w", errors.New("sql error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO books 
			(title, author, price, stock, isbn, created_at, updated_at) 
			VALUES(?,?,?,?,?,?, ?)`)).
				WithArgs(tt.args.book.Title, tt.args.book.Author, tt.args.book.Price,
					tt.args.book.Stock, tt.args.book.Isbn, AnyTime{}, AnyTime{}).
				WillReturnResult(tt.dbArgs.result).
				WillReturnError(tt.dbArgs.err)

			repo := repositories.NewBookRepository(sqlxDB)
			err := repo.CreateBook(tt.args.book)
			if tt.wantedErr != nil {
				assert.Equal(reflect.TypeOf(tt.wantedErr), reflect.TypeOf(err), "Error type is not the expected")
				assert.Equal(tt.wantedErr.Error(), err.Error(), "Error message is not the expected")
				return
			}
			assert.Equal(err, nil)
		})
	}
}
