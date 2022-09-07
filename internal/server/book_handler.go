package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/osalomon89/go-inventory/internal/domain"
	"github.com/osalomon89/go-inventory/internal/repositories"
)

var books []domain.Book

type ResponseInfo struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type Handler interface {
	getBookByID(w http.ResponseWriter, r *http.Request)
	getBooks(w http.ResponseWriter, r *http.Request)
	postBook(w http.ResponseWriter, r *http.Request)
	putBook(w http.ResponseWriter, r *http.Request)
	deleteBook(w http.ResponseWriter, r *http.Request)
	updateBook(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	repo repositories.BookRepository
}

func newHandler(bookRepository repositories.BookRepository) Handler {
	return &handler{
		repo: bookRepository,
	}
}

func (h *handler) getBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error " + idParam,
		})
		return
	}

	result, err := h.repo.GetBookByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "El libro no existe",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   result,
	})
}

func (h *handler) getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	author := r.URL.Query().Get("author")
	if author != "" {
		var sliceLibros []domain.Book
		for _, v := range books {
			if v.Author == author {
				sliceLibros = append(sliceLibros, v)
			}
		}
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: 200,
			Data:   sliceLibros,
		})
		return
	}
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: 200,
		Data:   books,
	})
}

func (h *handler) postBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var b domain.Book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	err = h.repo.CreateBook(&b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusInternalServerError,
			Data:   err,
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

	id, _ := strconv.ParseUint(idParam, 10, 32)
	var b domain.Book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}
	b.ID = uint(id)
	err = h.repo.UpdateBook(uint(id), &b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusInternalServerError,
			Data:   err,
		})
		return
	}

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: 200,
		Data:   b,
	})
}

func (h *handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
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
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   "Libro eliminado. ID: " + idParam,
	})
}

func (h *handler) updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error " + idParam,
		})
		return
	}

	bookRequestBody := new(domain.Book)
	err = json.NewDecoder(r.Body).Decode(bookRequestBody)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error decoding request body",
		})
		return
	}

	bookRequestBody.ID = uint(id)

	err = h.repo.UpdateBookByParams(uint(id), bookRequestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "El libro no existe",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   bookRequestBody,
	})
}
