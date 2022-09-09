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

type ResponseInfo struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type Handler interface {
	getBookByID(w http.ResponseWriter, r *http.Request)
	//getBooks(w http.ResponseWriter, r *http.Request)
	postBooks(w http.ResponseWriter, r *http.Request)
	putBook(w http.ResponseWriter, r *http.Request)
	patchBook(w http.ResponseWriter, r *http.Request)
	deleteBook(w http.ResponseWriter, r *http.Request)
	getBooksByParams(w http.ResponseWriter, r *http.Request)
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
	id, err := strconv.Atoi(param["id"])

	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   id,
		})
		return
	}
	result, err := h.repo.GetBooksById(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "Libro no encontrado",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   result,
	})
}

/*
	func (h *handler) getBooks(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//author := r.URL.Query().Get("author")
		//if author != "" {
		result, err := h.repo.GetBooks()
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
			Data:   result,
		})
		//return
		//}
		//json.NewEncoder(w).Encode(ResponseInfo{
		//	Status: 200,
		//	Data:   books,
		//})
	}
*/
var decoder = schema.NewDecoder()

type BookRequestQuery struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Isbn      string `json:"isbn"`
	Stock     string `json:"stock"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// handler for get book with query params
func (h *handler) getBooksByParams(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookRequestQuery := new(BookRequestQuery)
	err := decoder.Decode(bookRequestQuery, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "Ee",
		})
		return
	}
	bookQueryString, err := json.Marshal(bookRequestQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   err,
		})
		return
	}
	var params map[string]interface{}
	err = json.Unmarshal(bookQueryString, &params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   err,
		})
		return
	}
	result, err := h.repo.GetBooksByParams(params)

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
		Data:   result,
	})
}
func (h *handler) postBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var b domain.Book
	err := json.NewDecoder(r.Body).Decode(&b) //recibe la info
	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{ //envia la info
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
	_, err := strconv.Atoi(param["id"])
	var updateBook domain.Book
	erro := json.NewDecoder(r.Body).Decode(&updateBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "Error: Libro no encontrado",
		})
		return
	}
	if erro != nil {
		json.NewEncoder(w).Encode(ResponseInfo{ //envia la info
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}
	error := h.repo.UpdateBook(&updateBook)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusInternalServerError,
			Data:   error,
		})
		return
	}

	// for i, v := range books {
	// 	if v.ID == id {
	// 		books = append(books[:i], books[i+1:]...)
	// 		updateBook.ID = id
	// 		books = append(books, updateBook)
	// 	}
	// }

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusAccepted,
		Data:   "Libro actualizado",
	})

}
func (h *handler) patchBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	id, err := strconv.Atoi(param["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "Error: Libro no encontrado",
		})
		return
	}
	//busco el libro
	foundBook, err := h.repo.GetBooksById(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "Error: Libro no encontrado",
		})
		return
	}
	//decodificar el reques
	mapBook := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&mapBook)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{ //envia la info
			Status: http.StatusBadRequest,
			Data:   fmt.Errorf("Error al decodificar el body: %w", err),
		})
		return
	}

	err = h.repo.UpdateBookByParams(mapBook, foundBook)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   fmt.Errorf("Error al actualizar el libro: %w", err),
		})
		return
	}
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusAccepted,
		Data:   foundBook,
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
			Data:   id,
		})
		return
	}
	deletedBook, err := h.repo.GetBooksById(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "Error: Libro no encontrado",
		})
		return
	}
	err = h.repo.DeleteBook(deletedBook.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusInternalServerError,
			Data:   err,
		})
		return
	}
	/*var libro domain.Book
	libroVacio := domain.Book{}
	for i, v := range books {
		if v.ID == id {
			libro = v
			books = append(books[:i], books[i+1:]...)
		}
	}
	if libro == libroVacio {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "Libro no encontrado",
		})
		return
	}*/
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   "Libro eliminado correctamente",
	})
}
