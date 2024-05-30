package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/i18n"
	httphelper "github.com/Duncan-Kiragu/Msaada-Backend/pkg/http-helper"
)

func getDTO(c *fiber.Ctx, dto interface{}) error {
	if err := c.BodyParser(dto); err != nil {
		log.Println(err.Error())
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrInvalidDatas)
	}

	c.Locals(httphelper.LocalDTO, dto)
	return c.Next()
}

func GetProductDTO(c *fiber.Ctx) error {
	return getDTO(c, &dto.ProductInputDTO{})
}

func GetProfileDTO(c *fiber.Ctx) error {
	return getDTO(c, &dto.ProfileInputDTO{})
}

func GetUserDTO(c *fiber.Ctx) error {
	return getDTO(c, &dto.UserInputDTO{})
}

func GetPasswordInputDTO(c *fiber.Ctx) error {
	return getDTO(c, &dto.PasswordInputDTO{})
}
