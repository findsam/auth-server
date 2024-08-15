package user

import (
	"fmt"
	"io"
	"net/http"

	t "github.com/findsam/food-server/types"
	"github.com/go-chi/chi/v5"
)

func MakeHTTPHandlerFunc(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type Handler struct {
	store t.UserStore
}

func NewHandler(store t.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Get("/", MakeHTTPHandlerFunc(h.handleCreateUser))
		})
	})
	r.Group(func(r chi.Router) {
		r.Route("/ai", func(r chi.Router) {
			r.Post("/", MakeHTTPHandlerFunc(h.handleRequestDetails))
		})
	})
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	h.store.Create(ctx, fmt.Sprintf("%s", "1"))
	w.Write([]byte("User Created"))
	return nil
}

func (h *Handler) handleRequestDetails(w http.ResponseWriter, r *http.Request) error {
	// Read the body of the POST request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return err
	}
	defer r.Body.Close()

	// Print the body content
	// fmt.Fprintf(w, "Received: %s\n", string(body))
	w.Write([]byte(string(body)))
	return nil
}
