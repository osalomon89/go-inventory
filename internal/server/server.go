package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/osalomon89/go-inventory/internal/repositories"
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
		port: port}
}

func (r *httpRouter) SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/books/{id}", repositories.GetBookByID).Methods("GET")
	router.HandleFunc("/books", repositories.GetBooks).Methods("GET")
	router.HandleFunc("/books", repositories.PostBook).Methods("POST")
	router.HandleFunc("/books/{id}", repositories.PutBook).Methods("PUT")
	router.HandleFunc("/books/{id}", repositories.DeleteBook).Methods("DELETE")

	return router
}

func (r *httpRouter) Run(router http.Handler) error {
	return http.ListenAndServe(r.port, router)
}
func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(repositories.ResponseInfo{
		Status: http.StatusOK,
		Data:   "pong",
	})
}
