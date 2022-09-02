package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    `json:"id"`
	Author string `json:"author" validate:"required"`
	Title  string `json:"title"`
	Price  int    `json:"price"`
	Isbn   string `json:"isbn"`
	Stock  int    `json:"stock"`
}

var books []Book
var id_global int = 0

type ResponseInfo struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func main() {
	router := mux.NewRouter()

	const port string = ":8888"

	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/books/{id}", getBookByID).Methods("GET")
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books", postBook).Methods("POST")
	router.HandleFunc("/books/{id}", putBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	log.Println("Server listening on port", port)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalln(err)
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
func getBookByID(w http.ResponseWriter, r *http.Request) {
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
	var libro Book
	libroVacio := Book{}
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

func getBooks(w http.ResponseWriter, r *http.Request) {
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

func postBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var b Book
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

func putBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	id, err := strconv.Atoi(param["id"])
	var updateBook Book
	erro := json.NewDecoder(r.Body).Decode(&updateBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "Error: Libro no encontrado",
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

	if erro != nil {
		json.NewEncoder(w).Encode(ResponseInfo{ //envia la info
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusAccepted,
		Data:   "Libro actualizado",
	})

}
func deleteBook(w http.ResponseWriter, r *http.Request) {
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
	var libro Book
	libroVacio := Book{}
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
