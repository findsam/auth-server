package user

import (
	"context"
	"fmt"
	"time"

	"github.com/findsam/food-server/auth"
	t "github.com/findsam/food-server/types"
	u "github.com/findsam/food-server/util"
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
	user, err := NewAccount(b)

	if err != nil {
		return primitive.NilObjectID, err
	}

	col := s.db.Database(DbName).Collection(CollName)
	newUser, err := col.InsertOne(ctx, user)

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

	u := new(t.User)

	err = col.FindOne(ctx, bson.M{
		"_id": oID,
	}).Decode(u)

	if primitive.ObjectID.IsZero(u.ID) {
		return nil, nil
	}

	return u, err
}

func (s *Store) UpdatePassword(ctx context.Context, uid primitive.ObjectID, p string) error {
	col := s.db.Database(DbName).Collection(CollName)
	hashedPassword, err := auth.HashPassword(p)

	if err != nil {
		return err
	}
	_, err = col.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": bson.M{"password": hashedPassword, "meta.lastUpdate": time.Now().UTC()}})

	return err
}

func (s *Store) UpdateUser(ctx context.Context, b t.User) error {
	col := s.db.Database(DbName).Collection(CollName)

	result, err := col.UpdateOne(context.TODO(), bson.M{"_id": b.ID}, bson.M{
		"$set": bson.M{
			"firstName":       u.CapitalizeFirstLetter(b.FirstName),
			"lastName":        u.CapitalizeFirstLetter(b.LastName),
			"email":           b.Email,
			"meta.lastUpdate": time.Now().UTC(),
		},
	})
	fmt.Println(result.MatchedCount, result.ModifiedCount)
	return err
}
