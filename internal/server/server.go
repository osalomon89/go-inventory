package server

import (
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
		port: port,
	}
}

func (r *httpRouter) SetupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", ping).Methods("GET")

	dbConn, err := repositories.GetConnectionDB()
	if err != nil {
		panic("error db")
	}

	bookRepository := repositories.NewBookRepository(dbConn)

	bookHandler := newHandler(bookRepository)

	router.HandleFunc("/book/{id}", bookHandler.getBookByID).Methods("GET")
	router.HandleFunc("/books", bookHandler.getBooks).Methods("GET")
	router.HandleFunc("/book", bookHandler.getBooksByParam).Methods("GET")
	router.HandleFunc("/books", bookHandler.postBook).Methods("POST")
	router.HandleFunc("/book/{id}", bookHandler.putBook).Methods("PUT")
	router.HandleFunc("/book/{id}", bookHandler.patchBook).Methods("PATCH")
	router.HandleFunc("/book/{id}", bookHandler.deleteBook).Methods("DELETE")

	return router
}

func (r *httpRouter) Run(router http.Handler) error {
	return http.ListenAndServe(r.port, router)
}
