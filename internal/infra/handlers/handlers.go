package handlers

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"github.com/Duncan-Kiragu/Msaada-Backend/docs"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/api/handler"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/api/middleware"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/api/service"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/i18n"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/repository"
	httphelper "github.com/Duncan-Kiragu/Msaada-Backend/pkg/http-helper"

	"gorm.io/gorm"
)

var (
	profileRepository domain.ProfileRepository
	userRepository    domain.UserRepository
	productRepository domain.ProductRepository

	profileService domain.ProfileService
	userService    domain.UserService
	authService    domain.AuthService
	productService domain.ProductService
)

func initRepositories(postgresdb *gorm.DB) {
	// Create repositories.
	profileRepository = repository.NewProfileRepository(postgresdb)
	userRepository = repository.NewUserRepository(postgresdb)
	productRepository = repository.NewProductRepository(postgresdb)
}

func initServices() {
	// Create services.
	profileService = service.NewProfileService(profileRepository)
	userService = service.NewUserService(userRepository)
	authService = service.NewAuthService(userRepository)
	productService = service.NewProductService(productRepository)
}

func initHandelrs(app *fiber.App, db *gorm.DB) {
	reqMid := middleware.NewRequesttMiddleware(db)

	// Initialize access middleares
	middleware.MidAccess = middleware.Auth(os.Getenv("ACCESS_TOKEN_PUBLIC"), userRepository)
	middleware.MidRefresh = middleware.Auth(os.Getenv("RFRESH_TOKEN_PUBLIC"), userRepository)

	// Prepare endpoints for the API.
	handler.NewMiscHandler(app.Group(""))
	handler.NewAuthHandler(app.Group("/auth"), authService)
	handler.NewProfileHandler(app.Group("/profile"), profileService, reqMid)
	handler.NewUserHandler(app.Group("/user"), userService, reqMid)
	handler.NewProductHandler(app.Group("/product"), productService, reqMid)

	// Prepare an endpoint for 'Not Found'.
	app.All("*", func(c *fiber.Ctx) error {
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPResponse(c, fiber.StatusNotFound, messages.ErrorNonexistentRoute)
	})
}

func HandleRequests(app *fiber.App, postgresdb *gorm.DB) {
	if strings.ToLower(os.Getenv("API_SWAGGO")) == "true" {
		docs.SwaggerInfo.Version = os.Getenv("SYS_VERSION")

		// 	// Config swagger
		app.Get("/swagger/*", swagger.New(swagger.Config{
			DisplayRequestDuration: true,
			DocExpansion:           "none",
			ValidatorUrl:           "none",
		}))
	}

	initRepositories(postgresdb)
	initServices()
	initHandelrs(app, postgresdb)

	log.Fatal(app.Listen(":" + os.Getenv("API_PORT")))
}
