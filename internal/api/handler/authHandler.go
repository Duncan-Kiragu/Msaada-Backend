package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/api/middleware"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/i18n"
	httphelper "github.com/Duncan-Kiragu/Msaada-Backend/pkg/http-helper"
)

type AuthHandler struct {
	authService domain.AuthService
}

func (s *AuthHandler) handlerError(c *fiber.Ctx, err error) error {
	messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return httphelper.NewHTTPResponse(c, fiber.StatusUnauthorized, messages.ErrUserNotFound)
	}

	if errors.Is(err, domain.ErrInvalidIpAssociation) {
		return httphelper.NewHTTPResponse(c, fiber.StatusUnauthorized, messages.ErrInvalidIpAssociation)
	}

	switch err.Error() {
	case "invalid password":
		return httphelper.NewHTTPResponse(c, fiber.StatusUnauthorized, messages.ErrIncorrectPassword)
	case "invalid user":
		return httphelper.NewHTTPResponse(c, fiber.StatusUnauthorized, messages.ErrDisabledUser)
	}

	log.Println(err.Error())
	return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, messages.ErrGeneric)
}

// NewAuthHandler Creates a new authenticator handler.
func NewAuthHandler(route fiber.Router, as domain.AuthService) {
	handler := &AuthHandler{
		authService: as,
	}

	route.Post("", handler.login)
	route.Get("", middleware.MidAccess, handler.me)
	route.Put("", middleware.MidRefresh, handler.refresh)
}

// login godoc
// @Summary      User authentication
// @Description  User authentication
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        credentials body dto.AuthInputDTO true "Credentials model"
// @Success      200  {object}  dto.AuthOutputDTO
// @Failure      401  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /auth [post]
func (s *AuthHandler) login(c *fiber.Ctx) error {
	credentials := &dto.AuthInputDTO{}
	if err := c.BodyParser(credentials); err != nil {
		translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrInvalidDatas)
	}

	authResponse, err := s.authService.Login(c.Context(), credentials, c.IP())
	if err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(authResponse)
}

// me godoc
// @Summary      User authenticated
// @Description  User authenticated
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "User token"
// @Param        lang query string false "Language responses"
// @Success      200  {object}  dto.UserOutputDTO
// @Failure      401  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /auth [get]
// @Security	 Bearer
func (s *AuthHandler) me(c *fiber.Ctx) error {
	user := c.Locals(httphelper.LocalUser).(*domain.User)
	return c.Status(fiber.StatusOK).JSON(s.authService.Me(user))
}

// refresh godoc
// @Summary      User refresh
// @Description  User refresh
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "User token"
// @Param        lang query string false "Language responses"
// @Success      200  {object}  dto.AuthOutputDTO
// @Failure      401  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /auth [put]
func (s *AuthHandler) refresh(c *fiber.Ctx) error {
	user := c.Locals(httphelper.LocalUser).(*domain.User)
	return c.Status(fiber.StatusOK).JSON(s.authService.Refresh(user, c.IP()))
}
