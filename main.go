package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAdr := flag.String("listenAdr", ":8000", "Listen address of the server")
	app := fiber.New()
	apiv1 := app.Group("/api/v1")
	apiv1.Get("user", handleUser)
	app.Get("/temp", handleTemp)
	app.Listen(*listenAdr)
}

func handleTemp(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "hello"})
}

func handleUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"user": "nghianm93"})
}
