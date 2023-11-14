package db

import (
	"context"
	"fmt"
	"github.com/nghianm93/romo/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

const userColl = "users"

type UserStore interface {
	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(ctx context.Context) ([]*types.User, error)
	DeleteUser(context.Context, string) error
	InsertUser(context.Context, *types.User) (*types.User, error)
	UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error
	Dropper
}

type Dropper interface {
	Drop(context.Context) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	dbname := os.Getenv(DBNAME)
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbname).Collection(userColl),
	}
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- Dropping mongodb Collection ---")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {
	update := bson.D{
		{
			"$set", params.ToBSON(),
		},
	}
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil

}
func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	fmt.Println(id)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	if _, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	var user types.User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	fmt.Println(ctx)
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	var users []*types.User
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	fmt.Println(ctx)
	return users, nil
}
