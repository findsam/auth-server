package user

import (
	"fmt"
	"net/http"

	t "github.com/findsam/food-server/types"
	"github.com/go-chi/chi/v5"
)

func WrapHandler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
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
		r.Route("/api/v1", func(r chi.Router) {
			r.Get("/", WrapHandler(h.handleCreateUser))
		})
	})
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("method %s not allowed", r.Method)
	}
	ctx := r.Context()

	h.store.Create(ctx, fmt.Sprintf("%s", "1"))
	w.Write([]byte("User Created"))
	return nil
}
