package util

import (
	"encoding/json"
	"net/http"

	errors "github.com/findsam/food-server/error"
)

func GetTokenFromRequest(r *http.Request) string {
	authCookie, err := r.Cookie("Authorization")
	if err != nil {
		return ""
	}
	return authCookie.Value
}

func MakeHTTPHandlerFunc(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
func JSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ERROR(w http.ResponseWriter, status int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  errors.Message(status),
		"status": status,
	})
}
