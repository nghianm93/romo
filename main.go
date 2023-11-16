package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nghianm93/romo/api"
	"github.com/nghianm93/romo/db"
	"github.com/nghianm93/romo/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	errENV := godotenv.Load()
	if errENV != nil {
		log.Fatal(errENV)
	}
	listenAddr := flag.String("listenAddr", ":8000", "Listen address of the server")
	flag.Parse()
	databaseUrl := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(databaseUrl))
	if err != nil {
		log.Fatal(err)
	}
	userStore := *db.NewMongoUserStore(client)
	userHandler := api.NewUserHandler(&userStore)

	hostStore := *db.NewMongoHostStore(client)
	hostHandler := api.NewHostHandler(&hostStore)
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.JSON(types.ValidateMap{"error": err.Error()})
		},
	})
	apiv1 := app.Group("/api/v1")
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("user/:id", userHandler.HandleDeleteUser)

	apiv1.Post("/host", hostHandler.HandlePostUser)
	apiv1.Get("/host", hostHandler.HandleGetHosts)
	apiv1.Get("/host/:id", hostHandler.HandleGetHost)
	app.Listen(*listenAddr)
}
