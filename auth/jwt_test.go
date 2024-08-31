package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/findsam/food-server/config"
)

func TestCreateJWT(t *testing.T) {
	config.Envs.JWTSecret = "testsecret"

	uid := "12345"
	exp := time.Now().Add(time.Hour).Unix()

	tokenString, err := CreateJWT(uid, exp)
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if tokenString == "" {
		t.Error("expected a valid JWT token string, got an empty string")
	}
}

func TestWithJWT(t *testing.T) {
	config.Envs.JWTSecret = "testsecret"

	uid := "12345"
	exp := time.Now().Add(time.Hour).Unix()
	tokenString, _ := CreateJWT(uid, exp)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uidFromCtx := r.Context().Value("uid")
		if uidFromCtx == nil || uidFromCtx != uid {
			t.Errorf("expected uid %s in context, got %v", uid, uidFromCtx)
		}
	})

	rr := httptest.NewRecorder()
	WithJWT(handler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestValidateJWT(t *testing.T) {
	config.Envs.JWTSecret = "testsecret"

	uid := "12345"
	exp := time.Now().Add(time.Hour).Unix()
	validToken, _ := CreateJWT(uid, exp)

	token, err := ValidateJWT(validToken)
	if err != nil || !token.Valid {
		t.Errorf("expected valid token, got error: %v", err)
	}

	invalidToken := validToken + "invalid"
	_, err = ValidateJWT(invalidToken)
	if err == nil {
		t.Error("expected an error for invalid token, but got none")
	}
}

func TestValidateJWT_InvalidTokens(t *testing.T) {
	config.Envs.JWTSecret = "testsecret"

	uid := "12345"
	exp := time.Now().Add(time.Hour).Unix()
	validToken, _ := CreateJWT(uid, exp)

	manipulatedToken := validToken + "manipulated"
	_, err := ValidateJWT(manipulatedToken)
	if err == nil {
		t.Error("expected an error for manipulated token, but got none")
	}

	originalSecret := "testsecret"
	invalidSecret := "differentsecret"
	config.Envs.JWTSecret = invalidSecret

	_, err = ValidateJWT(validToken)
	if err == nil {
		t.Error("expected an error for token with different secret, but got none")
	}

	config.Envs.JWTSecret = originalSecret

	exp = time.Now().Add(-time.Hour).Unix()
	expiredToken, _ := CreateJWT(uid, exp)

	_, err = ValidateJWT(expiredToken)
	if err == nil {
		t.Error("expected an error for expired token, but got none")
	}
}
