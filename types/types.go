package types

import (
	"context"
	"time"

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
	Create(context.Context, RegisterRequest) error
	GetUserByID(context.Context, string) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
	UpdatePassword(context.Context, primitive.ObjectID, string) error
	UpdateUser(context.Context, UpdateUserRequest) error
	ArchiveUser(context.Context, string) error
}

type UserSecurity struct {
	EmailVerified bool  `json:"emailVerified" bson:"emailVerified"`
	HasTwoFactor  bool  `json:"hasTwoFactor" bson:"hasTwoFactor"`
	TwoFactorCode int32 `json:"twoFactorCode" bson:"twoFactorCode"`
}

type UserMeta struct {
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	LastUpdate time.Time `json:"lastUpdate" bson:"lastUpdate"`
	IsArchived bool      `json:"isArchived" bson:"isArchived"`
}

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password"`
	Security  UserSecurity       `json:"security" bson:"security"`
	Meta      UserMeta           `json:"meta" bson:"meta"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type ConfirmResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	ID        string `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Email     string `json:"email" bson:"email"`
}
