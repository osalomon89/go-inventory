package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	const port string = ":8888"

	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/book/{id}", getBookByID).Methods("GET")
	//router.HandleFunc("/books", getBookByID).Methods("POST")
	router.HandleFunc("/books", getBooks).Methods("GET")

	log.Println("Server listening on port", port)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalln(err)
	}
}

type ResponseInfo struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseInfo{
		Status: http.StatusOK,
		Data:   "id: " + idParam,
	})
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	code := r.URL.Query().Get("code")
	name := r.URL.Query().Get("name")

	if code == "" && name == "" {
		json.NewEncoder(w).Encode(ResponseInfo{
			Status: http.StatusBadRequest,
			Data:   "error",
		})
		return
	}

	json.NewEncoder(w).Encode(ResponseInfo{
		Status: 200,
		Data:   "code: " + code + ". name: " + name,
	})
}
