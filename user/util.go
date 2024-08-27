package user

import (
	"time"

	"github.com/findsam/food-server/auth"
	t "github.com/findsam/food-server/types"
	u "github.com/findsam/food-server/util"
)

func NewAccount(p t.RegisterRequest) (*t.User, error) {
	hashedPassword, err := auth.HashPassword(p.Password)
	if err != nil {
		return nil, err
	}

	return &t.User{
		Email:     p.Email,
		FirstName: u.CapitalizeFirstLetter(p.FirstName),
		LastName:  u.CapitalizeFirstLetter(p.LastName),
		Password:  string(hashedPassword),
		Meta: t.UserMeta{
			CreatedAt:  time.Now().UTC(),
			LastUpdate: time.Now().UTC(),
		},
		Security: t.UserSecurity{
			EmailVerified: false,
			HasTwoFactor:  false,
			TwoFactorCode: 0,
		},
	}, nil
}
