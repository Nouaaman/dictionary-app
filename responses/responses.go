package responses

import (
	"encoding/json"
	"net/http"
)

// JSONResponse sends a JSON response to the client
func JSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Error encoding JSON response")
	}
}

//handle Error
func HandleError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	resp := map[string]string{"error": message}
	json.NewEncoder(w).Encode(resp)
}
