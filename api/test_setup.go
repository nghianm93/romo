package api

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/nghianm93/romo/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

const (
	dburl = "MONGO_DB_URL_TEST"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) dropDown(t *testing.T) {
	dbname := os.Getenv(db.DBNAME)
	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	databaseURI := os.Getenv(dburl)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(databaseURI))
	if err != nil {
		t.Error(err)
	}
	return &testdb{
		client: client,
		Store: &db.Store{
			User: db.NewMongoUserStore(client),
		},
	}
}
