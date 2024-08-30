package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/findsam/food-server/auth"
	ge "github.com/findsam/food-server/error"
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
			//security requests
			r.Put("/user/confirm-reset-password", u.MakeHTTPHandlerFunc(h.handleConfirmResetPassword))
			r.Put("/user/reset-password", u.MakeHTTPHandlerFunc(h.handlePreResetPassword))
			//generation requests
			r.Post("/user/sign-up", u.MakeHTTPHandlerFunc(h.handleSignUp))
			r.Post("/user/sign-in", u.MakeHTTPHandlerFunc(h.handleSignIn))
			//token required requests
			r.Get("/user", auth.WithJWT(u.MakeHTTPHandlerFunc(h.handleSelf)))
			r.Put("/user", auth.WithJWT(u.MakeHTTPHandlerFunc(h.handleUpdateUser)))
			r.Delete("/user", auth.WithJWT(u.MakeHTTPHandlerFunc(h.handleArchiveUser)))
			//token generation requests
			r.Get("/user/refresh", u.MakeHTTPHandlerFunc(h.handleRefresh))
		})
	})
}

func (h *Handler) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	payload := new(t.RegisterRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return u.ERROR(w, ge.Internal)
	}

	existingUser, err := h.store.GetUserByEmail(r.Context(), payload.Email)

	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	if existingUser != nil {
		return u.ERROR(w, ge.EmailExists)
	}

	err = h.store.Create(r.Context(), *payload)

	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("Please verify your email address: %s", payload.Email),
	})
}

func (h *Handler) handleSelf(w http.ResponseWriter, r *http.Request) error {
	uid := r.Context().Value("uid").(string)
	user, err := h.store.GetUserByID(r.Context(), uid)

	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"results": []*t.User{user},
		"message": fmt.Sprintf("Successfully fetched: %s", uid),
	})
}

func (h *Handler) handleSignIn(w http.ResponseWriter, r *http.Request) error {
	payload := new(t.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return u.ERROR(w, ge.Internal)
	}

	user, err := h.store.GetUserByEmail(r.Context(), payload.Email)

	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	if user == nil || user.Meta.IsArchived {
		return u.ERROR(w, ge.UserNotFound)
	}

	if !auth.ComparePasswords(user.Password, []byte(payload.Password)) {
		return u.ERROR(w, ge.IncorrectCredentials)
	}

	access, err := createAndSetAuthCookies(user.ID.Hex(), w)

	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"results": []*t.User{user},
		"token":   access,
		"message": fmt.Sprintf("Successfully logged in as: %s", payload.Email),
	})
}

func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("refresh")
	if err != nil {
		return u.ERROR(w, ge.Unauthorized)
	}

	refresh, err := auth.ValidateJWT(cookie.Value)
	if err != nil || !refresh.Valid {
		return u.ERROR(w, ge.Internal)
	}

	uid := auth.ReadJWT(refresh)
	access, err := createAndSetAuthCookies(uid, w)
	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"token": access,
	})
}

func (h *Handler) handlePreResetPassword(w http.ResponseWriter, r *http.Request) error {
	/*********************************
	TODO: Send email via sendgrid/mailgun with url to reset via the token geneated in this request.
	*********************************/
	payload := new(t.ResetPasswordRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return u.ERROR(w, ge.Internal)
	}

	user, err := h.store.GetUserByEmail(r.Context(), payload.Email)

	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	if user == nil {
		return u.ERROR(w, ge.UserNotFound)
	}

	token, err := auth.CreateJWT(payload.Email, time.Now().Add(time.Minute*5).UTC().Unix())
	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	// for development purposes we'll log the token and put it in manually to save credits on sendgrid.
	fmt.Println("Token: ", token)

	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("Password reset email sent to %s", payload.Email),
	})
}

func (h *Handler) handleConfirmResetPassword(w http.ResponseWriter, r *http.Request) error {
	payload := new(t.ConfirmResetPasswordRequest)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return u.ERROR(w, ge.Internal)
	}

	token, err := auth.ValidateJWT(payload.Token)
	if err != nil || !token.Valid {
		return u.ERROR(w, ge.ResetExpired)
	}

	email := auth.ReadJWT(token)
	user, err := h.store.GetUserByEmail(r.Context(), email)

	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	if user == nil {
		return u.ERROR(w, ge.UserNotFound)
	}

	err = h.store.UpdatePassword(r.Context(), user.ID, payload.Password)

	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Password successfully changed",
	})
}

func (h *Handler) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	payload := new(t.User)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return u.ERROR(w, ge.Internal)
	}

	err := h.store.UpdateUser(r.Context(), *payload)

	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Sucessfully updated user",
	})
}

func (h *Handler) handleArchiveUser(w http.ResponseWriter, r *http.Request) error {
	uid := r.Context().Value("uid").(string)
	err := h.store.ArchiveUser(r.Context(), uid)

	if err != nil {
		return u.ERROR(w, ge.Internal)
	}

	return u.JSON(w, http.StatusOK, map[string]interface{}{
		"message":    "Your account has been archieved",
		"isArchived": true,
	})
}

func createAndSetAuthCookies(uid string, w http.ResponseWriter) (string, error) {
	access, err := auth.CreateJWT(uid, time.Now().Add(time.Minute*5).UTC().Unix())
	if err != nil {
		return "", u.ERROR(w, ge.Internal)
	}
	refresh, err := auth.CreateJWT(uid, time.Now().Add(time.Hour*24*7).UTC().Unix())
	if err != nil {
		return "", u.ERROR(w, ge.Internal)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh",
		Value:    refresh,
		Path:     "/users/user/refresh",
		Secure:   true,
		HttpOnly: true,
	})

	return access, nil
}
