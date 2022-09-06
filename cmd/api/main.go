package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/osalomon89/go-inventory/internal/repositories"
)

func main() {
	router := mux.NewRouter()

	const port string = ":8888"

	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/books/{id}", repositories.GetBookByID).Methods("GET")
	router.HandleFunc("/books", repositories.GetBooks).Methods("GET")
	router.HandleFunc("/books", repositories.PostBook).Methods("POST")
	router.HandleFunc("/books/{id}", repositories.PutBook).Methods("PUT")
	router.HandleFunc("/books/{id}", repositories.DeleteBook).Methods("DELETE")

	log.Println("Server listening on port", port)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalln(err)
	}
}
func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(repositories.ResponseInfo{
		Status: http.StatusOK,
		Data:   "pong",
	})
}
