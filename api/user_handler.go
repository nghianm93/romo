package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nghianm93/romo/db"
	"github.com/nghianm93/romo/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		userID = c.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.ValidateUpdateUserParams(); len(errors) > 0 {
		return c.JSON(errors)
	}
	filter := bson.M{"_id": oid}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": userID})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": userID})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.ValidateCreateUserParams(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	user, err := h.userStore.GetUserById(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
