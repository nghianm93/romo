package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/nghianm93/romo/api"
)

func main() {
	listenAdr := flag.String("listenAdr", ":8000", "Listen address of the server")

	app := fiber.New()
	apiv1 := app.Group("/api/v1")
	apiv1.Get("/user", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)

	app.Listen(*listenAdr)
}
