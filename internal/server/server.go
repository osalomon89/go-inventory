package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/osalomon89/go-inventory/internal/repositories"
)

// inyeccion de dependencia
// ---------------------------Router-----------------------------------

func NewHTTPRouter(port string) HTTPRouter {
	return &httpRouter{
		port: port,
	}
}

type HTTPRouter interface {
	Run(router http.Handler) error
	SetupRouter() *mux.Router
}

type httpRouter struct {
	port string
}

//------------------------------------------------------------------------

func (r *httpRouter) SetupRouter() *mux.Router {

	router := mux.NewRouter()

	dbConn, err := repositories.GetConnectionDB()
	if err != nil {
		panic("error db")
	}

	bookRepository := repositories.NewBookRepository(dbConn)

	bookHandler := newHandler(bookRepository)

	router.HandleFunc("/ping", bookHandler.ping).Methods("GET")
	router.HandleFunc("/books/{id}", bookHandler.getBookById).Methods("GET")
	router.HandleFunc("/books", bookHandler.getbooks).Methods("GET")
	router.HandleFunc("/books", bookHandler.postbook).Methods("POST")
	router.HandleFunc("/books/{id}", bookHandler.deletebook).Methods("DELETE")
	router.HandleFunc("/boooks/{id", bookHandler.putbook).Methods("PUT")

	return router
}

func (r *httpRouter) Run(router http.Handler) error {
	return http.ListenAndServe(r.port, router)
}
