package user

import (
	"context"
	"fmt"

	t "github.com/findsam/food-server/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "base"
	CollName = "users"
)

type Store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, b t.RegisterRequest) (primitive.ObjectID, error) {
	col := s.db.Database(DbName).Collection(CollName)
	newUser, err := col.InsertOne(ctx, b)

	id := newUser.InsertedID.(primitive.ObjectID)
	return id, err
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*t.User, error) {
	col := s.db.Database(DbName).Collection(CollName)

	u := new(t.User)

	err := col.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(u)

	if primitive.ObjectID.IsZero(u.ID) {
		return nil, nil
	}

	return u, err
}

func (s *Store) GetUserByID(ctx context.Context, id string) (*t.User, error) {
	col := s.db.Database(DbName).Collection(CollName)
	oID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	fmt.Println(oID)
	u := new(t.User)

	err = col.FindOne(ctx, bson.M{
		"_id": oID,
	}).Decode(u)

	if primitive.ObjectID.IsZero(u.ID) {
		return nil, nil
	}

	return u, err
}
