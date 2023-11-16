package db

import (
	"context"
	"github.com/nghianm93/romo/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type HostStore interface {
	InsertHost(ctx context.Context, host *types.Host) (*types.Host, error)
	GetHosts(ctx context.Context) ([]*types.Host, error)
	GetHostById(ctx context.Context, id string) (*types.Host, error)
}

type MongoHostStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHostStore(client *mongo.Client) *MongoHostStore {
	dbname := os.Getenv(DBNAME)
	return &MongoHostStore{
		client: client,
		coll:   client.Database(dbname).Collection(hostColl),
	}
}

func (s *MongoHostStore) InsertHost(ctx context.Context, host *types.Host) (*types.Host, error) {
	res, err := s.coll.InsertOne(ctx, host)
	if err != nil {
		return nil, err
	}
	host.ID = res.InsertedID.(primitive.ObjectID)
	return host, err
}

func (s *MongoHostStore) GetHosts(ctx context.Context) ([]*types.Host, error) {
	var hosts []*types.Host
	cur, err := s.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	if err := cur.All(ctx, &hosts); err != nil {
		return nil, err
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return hosts, nil
}

func (s *MongoHostStore) GetHostById(ctx context.Context, id string) (*types.Host, error) {
	var host types.Host
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&host); err != nil {
		return nil, err
	}
	return &host, nil

}
