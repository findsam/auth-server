package auth

import (
	"fmt"
	"net/http"

	"github.com/findsam/food-server/config"
	ge "github.com/findsam/food-server/error"
	u "github.com/findsam/food-server/util"
	"github.com/golang-jwt/jwt"
)

func CreateJWT(uid string, exp int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uid,
		"exp": exp,
	})

	str, err := token.SignedString([]byte(config.Envs.JWTSecret))

	if err != nil {
		return "", err
	}

	return str, err
}

func WithJWT(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := u.GetTokenFromRequest(r)
		token, err := ValidateJWT(tokenString)
		if err != nil {
			u.ERROR(w, ge.Internal)
			return
		}

		if !token.Valid {
			u.ERROR(w, ge.Internal)
			return
		}

		handlerFunc(w, r)
	}
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}
