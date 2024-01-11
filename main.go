package main

import (
	"dictionaryApp/dictionary"
	"dictionaryApp/middleware"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var dict *dictionary.Dictionary

func main() {
	dict = dictionary.New()

	// Load dictionary from file
	err := dict.LoadFromFile()
	if err != nil {
		fmt.Println("Error loading dictionary from file:", err)
		return
	}

	// new router
	r := mux.NewRouter()
	router := r.PathPrefix("/dictionary").Subrouter()
	// logging middleware
	router.Use(middleware.LoggingMiddleware)
	//
	router.HandleFunc("", handleAdd).Methods("POST")
	router.HandleFunc("/{word}", handleDefine).Methods("GET")
	router.HandleFunc("", handleList).Methods("GET")
	router.HandleFunc("/{word}", handleDelete).Methods("DELETE")

	fmt.Println("Listening on port 3000...")
	fmt.Println(http.ListenAndServe(":3000", router))
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Error encoding JSON response")
	}
}

/****************************************
 * HANDLERS
 ***************************************/
func handleAdd(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request body into an Entry struct
	var entry dictionary.Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}

	// Add the entry to the dictionary
	addErr := dict.Add(entry.Word, entry.Definition)
	if addErr != nil {
		if addErr == dictionary.ErrWordAlreadyExists {
			handleError(w, http.StatusConflict, "Word already exists")
		} else {
			handleError(w, http.StatusInternalServerError, "Error adding entry to dictionary")
		}
		return
	}

	// Save the dictionary to file
	err = dict.SaveToFile()
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Error saving dictionary to file")
		return
	}

	// Send success response
	resp := map[string]string{"message": fmt.Sprintf("Word %s added successfully", entry.Word)}
	jsonResponse(w, resp)
}

func handleDefine(w http.ResponseWriter, r *http.Request) {
	// Load the dictionary from file
	dict.LoadFromFile()

	// Extract the word from the request URL
	word := mux.Vars(r)["word"]

	// Get the definition for the word
	entry, err := dict.Get(word)
	if err != nil {
		if err == dictionary.ErrWordNotFound {
			handleError(w, http.StatusNotFound, fmt.Sprintf("Word %s not found", word))
		} else {
			handleError(w, http.StatusInternalServerError, "Error getting definition")
		}
		return
	}

	// Send response with word and definition
	resp := map[string]string{"word": word, "definition": entry.Definition}
	jsonResponse(w, resp)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	// Load the dictionary from file
	dict.LoadFromFile()

	// Extract the word from the request URL
	word := mux.Vars(r)["word"]

	// Delete the word from the dictionary
	err := dict.Remove(word)
	if err != nil {
		if err == dictionary.ErrWordNotFound {
			handleError(w, http.StatusNotFound, fmt.Sprintf("Word %s not found", word))
		} else {
			handleError(w, http.StatusInternalServerError, "Error removing word")
		}
		return
	}

	// Save the dictionary to file
	err = dict.SaveToFile()
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Error saving dictionary to file")
		return
	}

	// Send success response
	resp := map[string]string{"message": fmt.Sprintf("Word %s deleted successfully", word)}
	jsonResponse(w, resp)
}

//handle List
func handleList(w http.ResponseWriter, r *http.Request) {
	dict.LoadFromFile()
	entries := dict.List()
	jsonResponse(w, entries)
}

//handle Error
func handleError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	resp := map[string]string{"error": message}
	json.NewEncoder(w).Encode(resp)
}
