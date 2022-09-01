package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     uint   `json:"id"`
	Author string `json:"author" validate:"required"`
	Title  string `json:"title"`
	Price  int    `json:"price"`
	Isbn   string `json:"isbn"`
	Stock  int    `json:"stock"`
}

var idBook int = 1

var books []Book

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
}

type handler struct{}

func newHandler() Handler {
	return &handler{}
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
	var response Book

	for _, v := range books {
		if uint64(v.ID) == id {
			response = v
		}
	}

	if response == (Book{}) {
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
		Data:   response,
	})

}

func (h *handler) getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	author := r.URL.Query().Get("author")
	if author != "" {
		var sliceLibros []Book
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
	var b Book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	b.ID = uint(idBook)
	idBook++
	books = append(books, b)
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

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error " + idParam,
		})
		return
	}

	var b Book
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	var response Book

	for i, v := range books {
		if uint64(v.ID) == id {
			b.ID = v.ID
			books[i] = b
			response = books[i]
		}
	}

	if response == (Book{}) {
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
		Data:   response,
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
