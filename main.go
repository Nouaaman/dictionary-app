// main.go

package main

import (
	"dictionaryApp/dictionary"
	"dictionaryApp/middleware"
	"dictionaryApp/responses"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var dict *dictionary.Dictionary

func main() {
	var err error

	dict, err = dictionary.New()
	if err != nil {
		fmt.Println("Error initializing dictionary DB:", err)
		return
	}

	r := setupRouter()

	fmt.Println("Listening on port 3000...")
	fmt.Println(http.ListenAndServe(":3000", r))
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	router := r.PathPrefix("/dictionary").Subrouter()

	// middlewares
	router.Use(middleware.AuthenticationMiddleware)
	router.Use(middleware.LoggingMiddleware)

	// routes
	router.HandleFunc("", handleAdd).Methods("POST")
	router.HandleFunc("/{word}", handleDefine).Methods("GET")
	router.HandleFunc("", handleList).Methods("GET")
	router.HandleFunc("/{word}", handleDelete).Methods("DELETE")

	return r
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	var entry dictionary.Entry

	// decode the JSON request body into  Entry struct
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		responses.HandleError(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}

	if err := dict.Add(entry.Word, entry.Definition); err != nil {
		if err == dictionary.ErrWordAlreadyExists {
			responses.HandleError(w, http.StatusConflict, "Word already exists")
		} else {
			responses.HandleError(w, http.StatusInternalServerError, "Error adding entry to dictionary")
		}
		return
	}

	// send success response
	resp := map[string]string{"message": fmt.Sprintf("Word %s added successfully", entry.Word)}
	responses.JSONResponse(w, resp)
}

func handleDefine(w http.ResponseWriter, r *http.Request) {
	word := mux.Vars(r)["word"]

	// get the definition for the word
	entry, err := dict.Get(word)
	if err != nil {
		if err == dictionary.ErrWordNotFound {
			responses.HandleError(w, http.StatusNotFound, fmt.Sprintf("Word %s not found", word))
		} else {
			responses.HandleError(w, http.StatusInternalServerError, "Error getting definition")
		}
		return
	}

	// send response with data
	resp := map[string]string{"word": word, "definition": entry.Definition}
	responses.JSONResponse(w, resp)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	word := mux.Vars(r)["word"]

	if err := dict.Remove(word); err != nil {
		if err == dictionary.ErrWordNotFound {
			responses.HandleError(w, http.StatusNotFound, fmt.Sprintf("Word %s not found", word))
		} else {
			responses.HandleError(w, http.StatusInternalServerError, "Error removing word")
		}
		return
	}

	// send success response
	resp := map[string]string{"message": fmt.Sprintf("Word %s deleted successfully", word)}
	responses.JSONResponse(w, resp)
}

func handleList(w http.ResponseWriter, r *http.Request) {
	entries, err := dict.List()
	if err != nil {
		responses.HandleError(w, http.StatusInternalServerError, "Error retrieving word list")
		return
	}

	responses.JSONResponse(w, entries)
}
