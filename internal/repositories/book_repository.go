package repositories

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/osalomon89/go-inventory/internal/domain"
)

type BookRepository interface {
	GetBooks() []domain.Book
}

var books []domain.Book
var id_global int = 0

type ResponseInfo struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
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
	var libro domain.Book
	libroVacio := domain.Book{}
	for _, v := range books {
		if v.ID == id {
			libro = v
		}
	}
	if libro == libroVacio {
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
		Data:   libro,
	})
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
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

func PostBook(w http.ResponseWriter, r *http.Request) {
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
	id_global++
	b.ID = id_global
	books = append(books, b)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: 200,
		Data:   b,
	})

}

func PutBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	id, err := strconv.Atoi(param["id"])
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
	for i, v := range books {
		if v.ID == id {
			books = append(books[:i], books[i+1:]...)
			updateBook.ID = id
			books = append(books, updateBook)
		}
	}

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusAccepted,
		Data:   "Libro actualizado",
	})

}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
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
	var libro domain.Book
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
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   "Libro eliminado correctamente",
	})
}
