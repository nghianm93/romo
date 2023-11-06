package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nghianm93/romo/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Nghia",
		LastName:  "Minh",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "hello"})
}
