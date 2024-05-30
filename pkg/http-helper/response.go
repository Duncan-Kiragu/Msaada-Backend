package httphelper

import "github.com/gofiber/fiber/v2"

type HTTPResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func NewHTTPResponse(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(&HTTPResponse{
		Code:    status,
		Message: err.Error(),
	})
}
