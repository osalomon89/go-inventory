package server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/osalomon89/go-inventory/internal/domain"
	"github.com/osalomon89/go-inventory/internal/server"
	"github.com/osalomon89/go-inventory/internal/test/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_handler_PostBook(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		response *httptest.ResponseRecorder
		jsonBody string
	}

	type repoArgs struct {
		book  *domain.Book
		err   error
		times int
	}

	tests := []struct {
		name         string
		args         args
		repoArgs     repoArgs
		wantedStatus int
	}{
		{
			name: "Should get a success status",
			args: args{
				response: httptest.NewRecorder(),
				jsonBody: `{
					"id": 1,
					"author": "the-author",
					"title": "the title",
					"price": 5000,
					"isbn": "abcd1234",
					"stock": 20
				}`,
			},
			repoArgs: repoArgs{
				book: &domain.Book{
					ID:     1,
					Author: "the-author",
					Title:  "the title",
					Price:  5000,
					Isbn:   "abcd1234",
					Stock:  20,
				},
				err:   nil,
				times: 1,
			},
			wantedStatus: http.StatusOK,
		},
		{
			name: "Should get a bad request status",
			args: args{
				response: httptest.NewRecorder(),
				jsonBody: `{author:""`,
			},
			repoArgs: repoArgs{
				err:   nil,
				times: 0,
			},
			wantedStatus: http.StatusBadRequest,
		},
		{
			name: "Should get an internal server error status",
			args: args{
				response: httptest.NewRecorder(),
				jsonBody: `{
					"id": 1,
					"author": "the-author",
					"title": "the title",
					"price": 5000,
					"isbn": "abcd1234",
					"stock": 20
				}`,
			},
			repoArgs: repoArgs{
				book: &domain.Book{
					ID:     1,
					Author: "the-author",
					Title:  "the title",
					Price:  5000,
					Isbn:   "abcd1234",
					Stock:  20,
				},
				err:   errors.New("the repository error"),
				times: 1,
			},
			wantedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			bookRepositoryMock := mocks.NewMockBookRepository(mockCtrl)

			bookRepositoryMock.EXPECT().
				CreateBook(tt.repoArgs.book).
				Return(tt.repoArgs.err).
				Times(tt.repoArgs.times)

			handler := server.NewHandler(bookRepositoryMock)

			request, err := http.NewRequest("POST", "/books", strings.NewReader(tt.args.jsonBody))
			assert.Nil(err)

			handler.PostBook(tt.args.response, request)

			assert.Equal(tt.wantedStatus, tt.args.response.Code)
		})
	}
}
