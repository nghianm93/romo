package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nghianm93/romo/api"
	"github.com/nghianm93/romo/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const dburi = "mongodb://localhost:27017"

func main() {
	errENV := godotenv.Load()
	if errENV != nil {
		log.Fatal(errENV)
	}

	listenAddr := flag.String("listenAddr", ":8000", "Listen address of the server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {

		log.Fatal(err)
	}

	userStore := db.NewMongoUserStore(client)
	userHandler := api.NewUserHandler(userStore)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.JSON(map[string]string{"error": err.Error()})
		},
	})
	apiv1 := app.Group("/api/v1")

	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("user/:id", userHandler.HandleDeleteUser)

	app.Listen(*listenAddr)
}
