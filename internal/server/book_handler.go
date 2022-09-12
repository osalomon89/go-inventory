package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/osalomon89/go-inventory/internal/domain"
	"github.com/osalomon89/go-inventory/internal/repositories"
)

type ResponseInfo struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

var books []domain.Book

//--------------------------------------------------------------------//-inyeccion de dependencia de handler

func newHandler(bookRepository repositories.BookRepository) Handler {
	return &handler{
		repo: bookRepository,
	}
}

type Handler interface {
	ping(w http.ResponseWriter, r *http.Request)
	getBookById(w http.ResponseWriter, r *http.Request)
	getbooks(w http.ResponseWriter, r *http.Request)
	postbook(w http.ResponseWriter, r *http.Request)
	deletebook(w http.ResponseWriter, r *http.Request)
	putbook(w http.ResponseWriter, r *http.Request)
	updateBook(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	repo repositories.BookRepository
}

//--------------------------------------------------------------------------------------- //funcion ping

func (h *handler) ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   "Pong",
	})
}

//------------------------------------------------------------------------------------//funcion getBookBYId

func (h *handler) getBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	idparam := param["id"]

	id, err := strconv.ParseUint(idparam, 10, 32)

	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	result, err := h.repo.GetBookByID(uint(id))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "El libro no existe",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   result,
	})
}

//-----------------------------------------------------------------------------//funcion getbooks

func (h *handler) getbooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	author := r.URL.Query().Get("author")

	if author != "" {
		var libros []domain.Book
		for _, v := range books {
			if v.Author == author {
				libros = append(libros, v)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(ResponseInfo{
					Status: http.StatusOK,
					Data:   libros,
				})
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ResponseInfo{
					Status: http.StatusBadRequest,
					Data:   "Libro No Encontrado",
				})
				return
			}
		}

	}
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusBadRequest,
		Data:   "Libro No Encontrado",
	})

}

//-----------------------------------------------------------------------------------//POSTBOOK--

func (h *handler) postbook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var libroN domain.Book

	err := json.NewDecoder(r.Body).Decode(&libroN)

	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	err = h.repo.CreateBook(&libroN)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusInternalServerError,
			Data:   err,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   libroN,
	})

}

//------------------------------------------------//Funcion Delete-------------------------------------

func (h *handler) deletebook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param := mux.Vars(r)
	idparam := param["id"]
	id, err := strconv.ParseUint(idparam, 10, 32)

	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error" + idparam,
		})
		return
	}

	for i, v := range books {
		if uint64(v.ID) == id {
			books = append(books[:i], books[i+1:]...)
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   "Se elimino el Libro numero:" + idparam,
	})

}

//----------------------------------Func Put

func (h *handler) putbook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error " + idParam,
		})
		return
	}

	var libro domain.Book
	err = json.NewDecoder(r.Body).Decode(&libro)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	var newbook domain.Book

	for i, v := range books {
		if uint64(v.ID) == id {
			libro.ID = v.ID
			books[i] = libro
			newbook = books[i]
		}
	}

	if newbook == (domain.Book{}) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "El libro no existe",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   newbook,
	})
}

//------------------------------------------------Func Update

func (h *handler) updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error " + idParam,
		})
		return
	}

	foundBook, err := h.repo.GetBookByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "El libro a modificar no existe: " + idParam,
		})
		return
	}

	bookRequestBody := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&bookRequestBody)
	if err != nil {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error decoding request body",
		})
		return
	}

	err = h.repo.UpdateBookByParams(uint(id), bookRequestBody, foundBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "El libro no existe",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   foundBook,
	})
}
