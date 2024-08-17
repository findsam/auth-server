package auth

import (
	"github.com/findsam/food-server/config"
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
