package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/findsam/food-server/auth"
	t "github.com/findsam/food-server/types"
	u "github.com/findsam/food-server/util"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store t.UserStore
}

func NewHandler(store t.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/user/sign-up", u.MakeHTTPHandlerFunc(h.handleSignUp))
			r.Post("/user/sign-in", u.MakeHTTPHandlerFunc(h.handleSignIn))
			r.Get("/user/refresh", u.MakeHTTPHandlerFunc(h.handleRefresh))
			r.Get("/user/{id}", auth.WithJWT(u.MakeHTTPHandlerFunc(h.handleGetUser)))
		})
	})
}

func (h *Handler) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	payload := new(t.RegisterRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return u.ERROR(w, http.StatusBadRequest)
	}
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return u.ERROR(w, http.StatusBadRequest)
	}
	payload.Password = string(hashedPassword)
	_, err = h.store.Create(r.Context(), *payload)

	if err != nil {
		return u.ERROR(w, http.StatusUnauthorized)
	}
	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("Successfully created: %s", payload.Email),
		"status":  http.StatusOK,
	})
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	user, err := h.store.GetUserByID(r.Context(), id)

	if err != nil {
		return u.ERROR(w, http.StatusUnauthorized)
	}

	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"results": []*t.User{user},
		"message": fmt.Sprintf("Successfully fetched: %s", id),
		"status":  http.StatusOK,
	})
}

func (h *Handler) handleSignIn(w http.ResponseWriter, r *http.Request) error {
	payload := new(t.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return u.ERROR(w, http.StatusBadRequest)
	}
	user, err := h.store.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		return u.ERROR(w, http.StatusNoContent)
	}
	if !auth.ComparePasswords(user.Password, []byte(payload.Password)) {
		return u.ERROR(w, http.StatusUnauthorized)
	}

	err = createAndSetAuthCookies(user.ID.Hex(), w)
	if err != nil {
		return u.ERROR(w, http.StatusInternalServerError)
	}

	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"results": []*t.User{user},
		"message": fmt.Sprintf("Successfully logged in as: %s", payload.Email),
		"status":  http.StatusOK,
	})
}

func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("Refresh")
	if err != nil {
		return u.ERROR(w, http.StatusBadRequest)

	}
	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"cookie": cookie.Value,
	})
}

func createAndSetAuthCookies(uid string, w http.ResponseWriter) error {
	access, err := auth.CreateJWT(uid, time.Now().Add(time.Minute*15).UTC().Unix())
	if err != nil {
		return u.ERROR(w, http.StatusInternalServerError)
	}
	refresh, err := auth.CreateJWT(uid, time.Now().Add(time.Hour*24*7).UTC().Unix())
	if err != nil {
		return u.ERROR(w, http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Refresh",
		Value:    refresh,
		Path:     "/users/user/refresh",
		Secure:   true,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Path:  "/",
		Value: access,
	})

	return nil
}
