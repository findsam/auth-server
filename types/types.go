package types

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	Env              string
	Port             string
	MongoURI         string
	JWTSecret        string
	PublicURL        string
	APIKey           string
	ChatGPTSecretKey string
	ChatGPTURL       string
}

type LocationRequest struct {
	Postcode string `json:"postcode"`
}

type RecipesRequest struct {
	List []string `json:"list"`
}

type RegisterRequest struct {
	FirstName string `json:"firstName" bson:"firstName" validate:"required"`
	LastName  string `json:"lastName" bson:"lastName" validate:"required"`
	Email     string `json:"email" bson:"email" validate:"required"`
	Password  string `json:"password" bson:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserStore interface {
	Create(context.Context, RegisterRequest) (primitive.ObjectID, error)
	GetUserByID(context.Context, string) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password"`
}
