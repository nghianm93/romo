package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/nghianm93/romo/db"
	"github.com/nghianm93/romo/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	user, err := h.userStore.GetUserById(context.Background(), c.Params("id"))
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Nghia",
		LastName:  "Minh",
	}
	return c.JSON(u)
}
