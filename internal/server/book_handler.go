package server

import (
	"encoding/json"
	"fmt"
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
	getBooksByParam(w http.ResponseWriter, r *http.Request)
	postBook(w http.ResponseWriter, r *http.Request)
	patchBook(w http.ResponseWriter, r *http.Request)
	putBook(w http.ResponseWriter, r *http.Request)
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

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   "pong",
	})
}

func (h *handler) getBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	id, err := strconv.Atoi(param["id"])

	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	/*for _, v := range books {
		if v.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ResponseInfo{
				Status: http.StatusOK,
				Data:   v,
			})
			return
		}
	}*/

	result, err := h.repo.GetBookByID(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "libro no encontrado",
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
	/*author := r.URL.Query().Get("author")

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
	}*/

	result, err := h.repo.GetBook()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "libro no encontrado",
		})
		return
	}

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: 200,
		Data:   result,
	})
}

func (h *handler) getBooksByParam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	author := r.URL.Query().Get("author")

	result, err := h.repo.GetBookByParams(author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "autor no encontrado",
		})
		return
	}

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: 200,
		Data:   result,
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
	/*b.Id = len(books) + 1
	books = append(books, b)*/

	error := h.repo.CreateBook(&b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusInternalServerError,
			Data:   error,
		})
		return
	}

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: 200,
		Data:   b,
	})
}

func (h *handler) patchBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	id, err := strconv.Atoi(param["id"])

	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	foundBook, err := h.repo.GetBookByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   fmt.Sprintf("the book you are trying to modify does not exist. ID: %v", id),
		})
		return
	}

	bookRequestBody := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&bookRequestBody)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error decoding request body",
		})
		return
	}

	err = h.repo.UpdateBookByParams(id, bookRequestBody, foundBook)
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
		Data:   foundBook,
	})
}

func (h *handler) putBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	id, err := strconv.Atoi(param["id"])

	var newAtrib domain.Book

	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	error := json.NewDecoder(r.Body).Decode(&newAtrib)
	if error != nil {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	erro := h.repo.UpdateBook(id, &newAtrib)
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusInternalServerError,
			Data:   erro,
		})
		return
	}

	for i, v := range books {
		if v.Id == id {
			books = append(books[:i], books[i+1:]...)
			newAtrib.Id = id
			books = append(books, newAtrib)

			json.NewEncoder(w).Encode(ResponseInfo{
				Status: http.StatusOK,
				Data:   "actualización completa",
			})
			return
		}
	}

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   "actualización completa",
	})
}

func (h *handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	id, err := strconv.Atoi(param["id"])

	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	/*for i, v := range books {
		if v.Id == id {
			books = append(books[:i], books[i+1:]...)

			json.NewEncoder(w).Encode(ResponseInfo{
				Status: http.StatusOK,
				Data:   "libro eliminado",
			})
			return
		}
	}*/

	error := h.repo.DeleteBook(id)
	if error != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "libro no encontrado",
		})
		return
	}

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   "libro eliminado",
	})
}
