package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/i18n"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/filter"
	httphelper "github.com/Duncan-Kiragu/Msaada-Backend/pkg/http-helper"
)

func getQuery(c *fiber.Ctx, data interface{}) error {
	if err := c.QueryParser(data); err != nil {
		log.Println(err.Error())
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrInvalidDatas)
	}

	c.Locals(httphelper.LocalFilter, data)
	return c.Next()
}

func GetGenericFilter(c *fiber.Ctx) error {
	return getQuery(c, filter.NewFilter())
}

func GetUserFilter(c *fiber.Ctx) error {
	return getQuery(c, &filter.UserFilter{
		Filter:    *filter.NewFilter(),
		ProfileID: 0,
	})
}
