package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/book/{id}", getBookByID).Methods("GET")
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books", postBook).Methods("POST")
	router.HandleFunc("/book/{id}", putBook).Methods("PUT")
	router.HandleFunc("/book/{id}", patchBook).Methods("PATCH")
	router.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")

	return router
}

func Run(port string, router http.Handler) error {
	return http.ListenAndServe(port, router)
}
