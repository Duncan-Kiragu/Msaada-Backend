package handler

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/api/middleware"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/i18n"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/filter"
	httphelper "github.com/Duncan-Kiragu/Msaada-Backend/pkg/http-helper"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/pgerror"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/validator"
)

type UserHandler struct {
	userService domain.UserService
}

func (h *UserHandler) foreignKeyViolatedFrom(c *fiber.Ctx, messages *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrProfileNotFound)
	case fiber.MethodDelete:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrUserUsed)
	default:
		return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, messages.ErrGeneric)
	}
}

func (h *UserHandler) handlerError(c *fiber.Ctx, err error) error {
	messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch err := pgerror.HandlerError(err); {
	case errors.Is(err, pgerror.ErrDuplicatedKey):
		return httphelper.NewHTTPResponse(c, fiber.StatusConflict, messages.ErrUserRegistered)
	case errors.Is(err, pgerror.ErrForeignKeyViolated):
		return h.foreignKeyViolatedFrom(c, messages)
	case errors.Is(err, pgerror.ErrUndefinedColumn):
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrUndefinedColumn)
	}

	if errors.As(err, &validator.ErrValidator) {
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, err)
	}

	log.Println(err.Error())
	return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, messages.ErrGeneric)
}

func (h *UserHandler) getUserByEmail(c *fiber.Ctx) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	mail := strings.ReplaceAll(c.Params(httphelper.ParamMail), "%40", "@")
	user, err := h.userService.GetUserByMail(c.Context(), mail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httphelper.NewHTTPResponse(c, fiber.StatusNotFound, translation.ErrUserNotFound)
		}
		log.Println(err.Error())
		return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}

	c.Locals(httphelper.LocalObject, user)
	return c.Next()
}

// NewUserHandler Creates a new user handler.
func NewUserHandler(route fiber.Router, us domain.UserService, mid *middleware.RequesttMiddleware) {
	handler := &UserHandler{
		userService: us,
	}

	route.Patch("/:"+httphelper.ParamMail+"/passw", handler.getUserByEmail, middleware.GetPasswordInputDTO, handler.setUserPassword)

	route.Use(middleware.MidAccess)

	route.Get("", middleware.GetUserFilter, handler.getUsers)
	route.Post("", middleware.GetUserDTO, handler.createUser)
	route.Get("/:"+httphelper.ParamID, mid.UserByID, handler.getUser)
	route.Put("/:"+httphelper.ParamID, mid.UserByID, middleware.GetUserDTO, handler.updateUser)
	route.Delete("/:"+httphelper.ParamID, mid.UserByID, handler.deleteUser)
	route.Patch("/:"+httphelper.ParamID+"/reset", mid.UserByID, handler.resetUserPassword)
}

// getUsers godoc
// @Summary      Get users
// @Description  Get all users
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query filter.UserFilter false "Optional Filter"
// @Success      200  {array}   dto.ListItemsOutputDTO
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user [get]
// @Security	 Bearer
func (h *UserHandler) getUsers(c *fiber.Ctx) error {
	response, err := h.userService.GetUsers(c.Context(), c.Locals(httphelper.LocalFilter).(*filter.UserFilter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// createUser godoc
// @Summary      Insert user
// @Description  Insert user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        user body dto.UserInputDTO true "User model"
// @Success      201  {object}  dto.UserOutputDTO
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      409  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user [post]
// @Security	 Bearer
func (h *UserHandler) createUser(c *fiber.Ctx) error {
	userDTO := c.Locals(httphelper.LocalDTO).(*dto.UserInputDTO)
	user, err := h.userService.CreateUser(c.Context(), userDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// getUser godoc
// @Summary      Get user
// @Description  Get user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "User ID"
// @Success      200  {object}  dto.UserOutputDTO
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user/{id} [get]
// @Security	 Bearer
func (h *UserHandler) getUser(c *fiber.Ctx) error {
	user := c.Locals(httphelper.LocalObject).(*domain.User)

	return c.Status(fiber.StatusOK).JSON(&dto.UserOutputDTO{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Name,
		Status: user.Status,
		Profile: dto.ProfileOutputDTO{
			Id:   user.ProfileID,
			Name: user.Profile.Name,
		},
	})
}

// updateUser godoc
// @Summary      Update user
// @Description  Update user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "User ID"
// @Param        user body dto.UserInputDTO true "User model"
// @Success      200  {object}  dto.UserOutputDTO
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user/{id} [put]
// @Security	 Bearer
func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	userDTO := c.Locals(httphelper.LocalDTO).(*dto.UserInputDTO)
	oldUser := c.Locals(httphelper.LocalObject).(*domain.User)
	newUser, err := h.userService.UpdateUser(c.Context(), oldUser, userDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(newUser)
}

// deleteUser godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "User ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user/{id} [delete]
// @Security	 Bearer
func (h *UserHandler) deleteUser(c *fiber.Ctx) error {
	if err := h.userService.DeleteUser(c.Context(), c.Locals(httphelper.LocalObject).(*domain.User)); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// resetUser godoc
// @Summary      Reset user password
// @Description  Reset user password by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "User ID"
// @Success      200  {object}  nil
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user/{id}/reset [patch]
// @Security	 Bearer
func (h *UserHandler) resetUserPassword(c *fiber.Ctx) error {
	user := c.Locals(httphelper.LocalObject).(*domain.User)

	if !user.New {
		if err := h.userService.ResetUserPassword(c.Context(), user); err != nil {
			return h.handlerError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(nil)
	}

	return c.Status(fiber.StatusOK).JSON(nil)
}

// passwordUser godoc
// @Summary      Set user password
// @Description  Set user password by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        email     path    string     true        "User email"
// @Param        password body dto.PasswordInputDTO true "Password model"
// @Success      200  {object}  nil
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user/{email}/passw [patch]
func (h *UserHandler) setUserPassword(c *fiber.Ctx) error {
	pass := c.Locals(httphelper.LocalDTO).(*dto.PasswordInputDTO)
	if !pass.IsValid() {
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrPassUnmatch)
	}

	user := c.Locals(httphelper.LocalObject).(*domain.User)
	if !user.New {
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrUserHasPass)
	}

	if err := h.userService.SetUserPassword(c.Context(), user, pass); err != nil {
		return h.handlerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(nil)
}
