package middleware

import (
	"dictionaryApp/responses"
	"net/http"
)

const authToken = "XdkmKQ79Q979XxEVHcUS1liWoDUhPACF1DuAGWUqPSxlWRr3JHr86l2okqOh"

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if token == "" {
			responses.HandleError(w, http.StatusUnauthorized, "Missing authentication token")
			return
		}

		if token != "Bearer "+authToken {
			responses.HandleError(w, http.StatusUnauthorized, "Invalid authentication token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
