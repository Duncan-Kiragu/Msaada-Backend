package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/api/middleware"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/i18n"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/filter"
	httphelper "github.com/Duncan-Kiragu/Msaada-Backend/pkg/http-helper"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/pgerror"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/validator"
)

type ProductHandler struct {
	productService domain.ProductService
}

func (h *ProductHandler) foreignKeyViolatedMethod(c *fiber.Ctx, translation *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrProductNotFound)
	case fiber.MethodDelete:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrProductUsed)
	default:
		return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}
}

func (h *ProductHandler) handlerError(c *fiber.Ctx, err error) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch err := pgerror.HandlerError(err); {
	case errors.Is(err, pgerror.ErrDuplicatedKey):
		return httphelper.NewHTTPResponse(c, fiber.StatusConflict, translation.ErrProductRegistered)
	case errors.Is(err, pgerror.ErrForeignKeyViolated):
		return h.foreignKeyViolatedMethod(c, translation)
	case errors.Is(err, pgerror.ErrUndefinedColumn):
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrUndefinedColumn)
	}

	if errors.As(err, &validator.ErrValidator) {
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, err)
	}

	log.Println(err.Error())
	return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
}

// NewProductHandler Creates a new product handler.
func NewProductHandler(route fiber.Router, ps domain.ProductService, mid *middleware.RequesttMiddleware) {
	handler := &ProductHandler{
		productService: ps,
	}

	route.Use(middleware.MidAccess)

	route.Get("", middleware.GetGenericFilter, handler.getProducts)
	route.Post("", middleware.GetProductDTO, handler.createProduct)
	route.Get("/:"+httphelper.ParamID, mid.ProductByID, handler.getProductBydID)
	route.Put("/:"+httphelper.ParamID, mid.ProductByID, middleware.GetProductDTO, handler.updateProduct)
	route.Delete("/:"+httphelper.ParamID, mid.ProductByID, handler.deleteProduct)
}

// getProducts godoc
// @Summary      Get products
// @Description  Get products
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query filter.Filter false "Optional Filter"
// @Success      200  {array}   dto.ListItemsOutputDTO
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /product [get]
// @Security	 Bearer
func (h *ProductHandler) getProducts(c *fiber.Ctx) error {
	response, err := h.productService.GetProducts(c.Context(), c.Locals(httphelper.LocalFilter).(*filter.Filter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// getProductBydID godoc
// @Summary      Get product by ID
// @Description  Get product by ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id   path			int			true        "Product ID"
// @Success      200  {object}  dto.ProductOutputDTO
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /product/{id} [get]
// @Security	 Bearer
func (h *ProductHandler) getProductBydID(c *fiber.Ctx) error {
	product := c.Locals(httphelper.LocalObject).(*domain.Product)
	return c.Status(fiber.StatusOK).JSON(&dto.ProductOutputDTO{
		Id:   product.Id,
		Name: product.Name,
	})
}

// createProduct godoc
// @Summary      Insert product
// @Description  Insert product
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        product body dto.ProductInputDTO true "Product model"
// @Success      201  {object}  dto.ProductOutputDTO
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      409  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /product [post]
// @Security	 Bearer
func (h *ProductHandler) createProduct(c *fiber.Ctx) error {
	productDTO := c.Locals(httphelper.LocalDTO).(*dto.ProductInputDTO)
	product, err := h.productService.CreateProduct(c.Context(), productDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// updateProduct godoc
// @Summary      Update product by ID
// @Description  Update product by ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Product ID"
// @Param        product body dto.ProductInputDTO true "Product model"
// @Success      200  {object}  dto.ProductOutputDTO
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /product/{id} [put]
// @Security	 Bearer
func (h *ProductHandler) updateProduct(c *fiber.Ctx) error {
	productDTO := c.Locals(httphelper.LocalDTO).(*dto.ProductInputDTO)
	oldProduct := c.Locals(httphelper.LocalObject).(*domain.Product)
	newProduct, err := h.productService.UpdateProduct(c.Context(), oldProduct, productDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(newProduct)
}

// deleteProduct godoc
// @Summary      Delete product by ID
// @Description  Delete product by ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Product ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /product/{id} [delete]
// @Security	 Bearer
func (h *ProductHandler) deleteProduct(c *fiber.Ctx) error {
	product := c.Locals(httphelper.LocalObject).(*domain.Product)
	if err := h.productService.DeleteProduct(c.Context(), product); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
