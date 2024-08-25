package auth

import (
	"time"

	t "github.com/findsam/food-server/types"
	u "github.com/findsam/food-server/util"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePasswords(hashed string, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	return err == nil
}

func NewAccount(p t.RegisterRequest) (*t.User, error) {
	hashedPassword, err := HashPassword(p.Password)
	if err != nil {
		return nil, err
	}

	return &t.User{
		Email:     p.Email,
		FirstName: u.CapitalizeFirstLetter(p.FirstName),
		LastName:  u.CapitalizeFirstLetter(p.LastName),
		Password:  string(hashedPassword),
		Meta: t.UserMeta{
			CreatedAt: time.Now().UTC(),
		},
		Security: t.UserSecurity{
			EmailVerified: false,
			HasTwoFactor:  false,
			TwoFactorCode: 0,
		},
	}, nil
}
