package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type MiscHandler struct{}

func NewMiscHandler(miscRoute fiber.Router) {
	handler := &MiscHandler{}

	miscRoute.Get("", handler.healthCheck).Name("Root")
	//miscRoute.Get("/monitor", monitor.New(monitor.Config{
	//	Title:   "Server Monitor",
	//	Refresh: 5 * time.Second,
	//	APIOnly: false,
	//}))
}

// healthCheck godoc
// @Summary      Ping Pong
// @Description  Ping Pong
// @Tags         Ping
// @Produce      json
// @Success      200  {object}   map[string]string
// @Router       / [get]
func (h *MiscHandler) healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"time": time.Now(),
	})
}
