package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPRouter interface {
	SetupRouter() *mux.Router
	Run(router http.Handler) error
}

type httpRouter struct {
	port string
}

func NewHTTPRouter(port string) HTTPRouter {
	return &httpRouter{
		port: port,
	}
}

func (r *httpRouter) SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", ping).Methods("GET")

	bookHandler := newHandler()

	router.HandleFunc("/books", bookHandler.getBooks).Methods("GET")
	router.HandleFunc("/books", bookHandler.postBook).Methods("POST")
	router.HandleFunc("/books/{id}", bookHandler.getBookByID).Methods("GET")
	router.HandleFunc("/books/{id}", bookHandler.putBook).Methods("PUT")
	router.HandleFunc("/books/{id}", bookHandler.deleteBook).Methods("DELETE")

	return router
}

func (r *httpRouter) Run(router http.Handler) error {
	return http.ListenAndServe(r.port, router)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   "pong",
	})
}
