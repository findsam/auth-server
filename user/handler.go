package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
			r.Post("/user/sign-up", MakeHTTPHandlerFunc(h.handleSignUp))
			r.Post("/user/sign-in", MakeHTTPHandlerFunc(h.handleSignIn))
			r.Get("/user/{id}", MakeHTTPHandlerFunc(h.handleGetUser))
		})
	})
}

func (h *Handler) handleSignUp(w http.ResponseWriter, r *http.Request) error {
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

func (h *Handler) handleSignIn(w http.ResponseWriter, r *http.Request) error {
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

	err = createAndSetAuthCookies(user.ID.Hex(), w)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, payload)
}

func createAndSetAuthCookies(uid string, w http.ResponseWriter) error {
	accessExpiry := time.Now().Add(time.Minute * 15)
	refreshExpiry := time.Now().Add(time.Hour * 24 * 7)

	access, err := auth.CreateJWT(uid, accessExpiry.UTC().Unix())
	if err != nil {
		return err
	}
	refresh, err := auth.CreateJWT(uid, refreshExpiry.UTC().Unix())
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Refresh",
		Value:    refresh,
		Path:     "/auth/refresh",
		Secure:   true,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: access,
	})

	return nil
}
