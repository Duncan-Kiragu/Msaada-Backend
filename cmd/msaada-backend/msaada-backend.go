package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	_ "github.com/Duncan-Kiragu/Msaada-Backend/configs"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/api/middleware"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/infra/database"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/infra/handlers"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/i18n"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/helpers"
	httphelper "github.com/Duncan-Kiragu/Msaada-Backend/pkg/http-helper"
)

// @title 							Go - Template API
// @description 					Template API.

// @contact.name					Raul del Aguila
// @contact.email					email@email.com

// @BasePath						/

// @securityDefinitions.apiKey		Bearer
// @in								header
// @name							Authorization
// @description 					Type "Bearer" followed by a space and JWT token.
func main() {
	postgresdb, err := database.ConnectPostgresDB()
	helpers.PanicIfErr(err)

	app := fiber.New(fiber.Config{
		EnablePrintRoutes:     false,
		Prefork:               os.Getenv("SYS_PREFORK") == "true",
		CaseSensitive:         true,
		StrictRouting:         true,
		DisableStartupMessage: false,
		AppName:               "Go - Template API",
		ReduceMemoryUsage:     false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, err)
		},
	})

	app.Use(
		recover.New(),
		middleware.GetRequestLanguage,
		requestid.New(),
	)

	if strings.ToLower(os.Getenv("API_LOGGER")) == "true" {
		app.Use(logger.New(logger.Config{
			CustomTags: map[string]logger.LogFunc{
				"xip": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
					return output.WriteString(fmt.Sprintf("%15s", c.IP()))
				},
				"fullPath": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
					return output.WriteString(string(c.Request().RequestURI()))
				},
			},
			Format:     "[FIBER:${magenta}${pid}${reset}] ${time} | ${status} | ${latency} | ${xip} | ${method} ${fullPath} ${magenta}${error}${reset}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   time.Local.String(),
		}))
	}

	app.Use(
		cors.New(cors.Config{
			AllowOrigins:  "*",
			AllowMethods:  strings.Join([]string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete, fiber.MethodOptions}, ","),
			AllowHeaders:  "*",
			ExposeHeaders: "*",
			MaxAge:        1,
			// AllowCredentials: true,
		}),
		limiter.New(limiter.Config{
			Max:        200,
			Expiration: time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
				return httphelper.NewHTTPResponse(c, fiber.StatusTooManyRequests, messages.ErrManyRequest)
			},
		}),
	)

	handlers.HandleRequests(app, postgresdb)
}
