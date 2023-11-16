package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nghianm93/romo/db"
	"github.com/nghianm93/romo/types"
)

type HostHandler struct {
	hostStore db.HostStore
}

func NewHostHandler(hostStore db.HostStore) *HostHandler {
	return &HostHandler{
		hostStore: hostStore,
	}
}

func (h *HostHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateHostParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	host, err := types.NewHostFromParams(params)
	if err != nil {
		return err
	}
	insertedHost, err := h.hostStore.InsertHost(c.Context(), host)
	if err != nil {
		return err
	}
	return c.JSON(insertedHost)
}

func (h *HostHandler) HandleGetHosts(c *fiber.Ctx) error {
	hosts, err := h.hostStore.GetHosts(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hosts)
}

func (h *HostHandler) HandleGetHost(c *fiber.Ctx) error {
	host, err := h.hostStore.GetHostById(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	return c.JSON(host)
}
