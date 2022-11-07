package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
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
		port: port,
	}
}

func (r *httpRouter) SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", ping).Methods("GET")
	router.Handle("/swagger.yaml", http.FileServer(http.Dir("../../")))

	// documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "../../swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handle("/docs", sh)

	// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "../../swagger.yaml", Path: "docs1"}
	sh1 := middleware.Redoc(opts1, nil)
	router.Handle("/docs1", sh1)

	dbConn, err := repositories.GetConnectionDB()
	if err != nil {
		panic("error db")
	}

	bookRepository := repositories.NewBookRepository(dbConn)

	bookHandler := newHandler(bookRepository)

	router.HandleFunc("/books", bookHandler.getBooks).Methods("GET")
	router.HandleFunc("/books", bookHandler.createBook).Methods("POST")
	router.HandleFunc("/books/{id}", bookHandler.getBookByID).Methods("GET")
	router.HandleFunc("/books/{id}", bookHandler.putBook).Methods("PUT")
	router.HandleFunc("/books/{id}", bookHandler.updateBook).Methods("PATCH")
	router.HandleFunc("/books/{id}", bookHandler.deleteBook).Methods("DELETE")

	return router
}

func (r *httpRouter) Run(router http.Handler) error {
	return http.ListenAndServe(r.port, router)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode("pong")
}
