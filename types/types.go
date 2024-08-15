package types

import (
	"context"
)

type Config struct {
	Env       string
	Port      string
	MongoURI  string
	JWTSecret string
	PublicURL string
	APIKey    string
}

type RegisterRequest struct {
	FirstName string `json:"firstName" bson:"firstName" validate:"required"`
	LastName  string `json:"lastName" bson:"lastName" validate:"required"`
	Email     string `json:"email" bson:"email" validate:"required"`
	Password  string `json:"password" bson:"password" validate:"required"`
}

type UserStore interface {
	Create(context.Context, string) (string, error)
}
