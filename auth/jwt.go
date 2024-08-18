package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/findsam/food-server/config"
	u "github.com/findsam/food-server/util"
	"github.com/golang-jwt/jwt"
)

func permissionDenied(w http.ResponseWriter) {
	u.JSON(w, http.StatusUnauthorized,
		map[string]interface{}{
			"Error": fmt.Errorf("permission denied").Error(),
		})
}

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

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := u.GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}
