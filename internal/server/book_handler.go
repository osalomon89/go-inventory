package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/osalomon89/go-inventory/internal/domain"
	"github.com/osalomon89/go-inventory/internal/repositories"
)

var decoder = schema.NewDecoder()
var books []domain.Book

type Handler interface {
	getBookByID(w http.ResponseWriter, r *http.Request)
	getBooks(w http.ResponseWriter, r *http.Request)
	createBook(w http.ResponseWriter, r *http.Request)
	putBook(w http.ResponseWriter, r *http.Request)
	updateBook(w http.ResponseWriter, r *http.Request)
	deleteBook(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	repo repositories.BookRepository
}

func newHandler(bookRepository repositories.BookRepository) Handler {
	return &handler{
		repo: bookRepository,
	}
}

// swagger:route GET /books/{id} getBookByID getBookByID
// This is the description for getting a book by its ID. Which can be longer.
// Responses:
// - 200: ResponseInfo
// - 400: ResponseError
func (h *handler) getBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   "error " + idParam,
		})
		return
	}

	result, err := h.repo.GetBookByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   "El libro no existe",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   *result,
	})
}

// swagger:route GET /books getBooks getBooks
//
// Lists books filtered by some parameters.
//
// This will show all available books by default.
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		Schemes: http, https
//
//		Parameters:
//		  + name: limit
//		    in: query
//		    description: maximum numnber of results to return
//		    required: false
//		    type: integer
//		    format: int32
//		  + name: offset
//		    in: query
//		    description: number of results to skip
//		    required: false
//		    type: integer
//		    format: int32
//		  + name: isbn
//		    in: query
//		    description: isbn number
//		    required: false
//		    type: string
//
//	    Responses:
//	      200: ResponseAllInfo
//	      400: ResponseError
func (h *handler) getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookRequestQuery := new(BookRequestQuery)
	err := decoder.Decode(bookRequestQuery, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   "El libro no existe",
		})
		return
	}

	bookRequestQueryString, err := json.Marshal(bookRequestQuery)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusInternalServerError,
			Data:   err.Error(),
		})
		return
	}

	var params map[string]interface{}
	err = json.Unmarshal(bookRequestQueryString, &params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusInternalServerError,
			Data:   err.Error(),
		})
		return
	}

	booksResult, err := h.repo.GetBooks(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   "El libro no existe",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseAllInfo{
		Status: 200,
		Data:   booksResult,
	})
}

// swagger:route POST /books createBook createBook
//
// Create book.
//
// This will create a book.
//
//	Consumes:
//		- application/json
//
//	Produces:
//		- application/json
//
//	Schemes: http, https
//
//	Parameters:
//		+ name: Body
//		in: body
//		description: body with book parameters
//		required: true
//		schema:
//	   		"$ref": "#/definitions/Body"
//
//	Responses:
//		200: ResponseInfo
//		400: ResponseError
func (h *handler) createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var b domain.Book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	err = h.repo.CreateBook(&b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusInternalServerError,
			Data:   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: 200,
		Data:   b,
	})
}

func (h *handler) putBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   "error " + idParam,
		})
		return
	}

	var b domain.Book
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	var response domain.Book

	for i, v := range books {
		if uint64(v.ID) == id {
			b.ID = v.ID
			books[i] = b
			response = books[i]
		}
	}

	if response == (domain.Book{}) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   "El libro no existe",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   response,
	})
}

func (h *handler) updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   "error " + idParam,
		})
		return
	}

	foundBook, err := h.repo.GetBookByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   fmt.Sprintf("the book you are trying to modify does not exist. ID: %s", idParam),
		})
		return
	}

	bookRequestBody := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&bookRequestBody)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   "error decoding request body",
		})
		return
	}

	err = h.repo.UpdateBookByParams(bookRequestBody, foundBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   *foundBook,
	})
}

func (h *handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseError{
			Status: http.StatusBadRequest,
			Data:   "error " + idParam,
		})
		return
	}

	for i, v := range books {
		if uint64(v.ID) == id {
			books = append(books[:i], books[i+1:]...)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseDeleteInfo{
		Status: http.StatusOK,
		Data:   "Libro eliminado. ID: " + idParam,
	})
}
