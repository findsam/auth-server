package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "notebase"
	CollName = "users"
)

type Store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, email string) (string, error) {
	// col := s.db.Database(DbName).Collection(CollName)
	fmt.Println(email)
	return "", nil
}
