package main

import (
	"dictionaryApp/dictionary"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var dict *dictionary.Dictionary

func main() {
	dict = dictionary.New()
	dict.LoadFromFile()

	r := mux.NewRouter()
	router := r.PathPrefix("/dictionary").Subrouter()
	router.HandleFunc("", handleAdd).Methods("POST")
	router.HandleFunc("/{word}", handleDefine).Methods("GET")
	router.HandleFunc("", handleList).Methods("GET")
	router.HandleFunc("/{word}", handleDelete).Methods("DELETE")

	fmt.Println("Listening on port 3000...")
	fmt.Println(http.ListenAndServe(":3000", router))
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	var entry dictionary.Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, "plz provide a Word and his Definition", http.StatusBadRequest)
		return
	}

	dict.Add(entry.Word, entry.Definition)
	dict.SaveToFile()

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": entry.Word + " added."}
	jsonResponse(w, response)
}

func handleDefine(w http.ResponseWriter, r *http.Request) {

	dict.LoadFromFile()

	params := mux.Vars(r)
	word := params["word"]

	entry, err := dict.Get(word)
	if err != nil {
		http.Error(w, word+" doesn't exist", http.StatusNotFound)
		return
	}

	response := map[string]string{"word": word, "definition": entry.Definition}
	jsonResponse(w, response)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	word := params["word"]

	dict.Remove(word)
	dict.SaveToFile()

	response := map[string]string{"message": fmt.Sprintf("%s deleted.", word)}
	jsonResponse(w, response)
}

func handleList(w http.ResponseWriter, r *http.Request) {
	dict.LoadFromFile()
	entries := dict.List()
	jsonResponse(w, entries)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
	}
}
