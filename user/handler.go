package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/findsam/food-server/auth"
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

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type Handler struct {
	store t.UserStore
}

func NewHandler(store t.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/user", MakeHTTPHandlerFunc(h.handleRegister))
			r.Get("/user/{id}", MakeHTTPHandlerFunc(h.handleGetUser))
		})
	})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) error {
	payload := new(t.RegisterRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return err
	}
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return err
	}
	payload.Password = string(hashedPassword)
	id, err := h.store.Create(r.Context(), *payload)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, id)
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	user, err := h.store.GetUserByID(r.Context(), id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("method %s not allowed", r.Method)
	}

	payload := new(t.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return err
	}

	user, err := h.store.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		return err
	}

	if !auth.ComparePasswords(user.Password, []byte(payload.Password)) {
		return fmt.Errorf("invalid password or user does not exist")
	}

	return WriteJSON(w, http.StatusOK, user)
}
