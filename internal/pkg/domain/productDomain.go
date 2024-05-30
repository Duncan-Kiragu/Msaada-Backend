package domain

import (
	"context"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/filter"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/validator"
)

const ProductTableName string = "product"

type (
	Product struct {
		Base
		Name string `json:"name" gorm:"column:name;type:varchar(100);unique;index;not null;" validate:"required,min=2"`
	}

	ProductRepository interface {
		CountProducts(context.Context, *filter.Filter) (int64, error)
		GetProductByID(context.Context, uint) (*Product, error)
		GetProducts(context.Context, *filter.Filter) (*[]Product, error)
		CreateProduct(context.Context, *dto.ProductInputDTO) (*Product, error)
		UpdateProduct(context.Context, *Product, *dto.ProductInputDTO) error
		DeleteProduct(context.Context, *Product) error
	}

	ProductService interface {
		GetProductByID(context.Context, uint) (*dto.ProductOutputDTO, error)
		GetProducts(context.Context, *filter.Filter) (*dto.ListItemsOutputDTO, error)
		CreateProduct(context.Context, *dto.ProductInputDTO) (*dto.ProductOutputDTO, error)
		UpdateProduct(context.Context, *Product, *dto.ProductInputDTO) (*dto.ProductOutputDTO, error)
		DeleteProduct(context.Context, *Product) error
	}
)

func (s *Product) TableName() string {
	return ProductTableName
}

func (s *Product) Bind(p *dto.ProductInputDTO) error {
	if p.Name != nil {
		s.Name = *p.Name
	}

	return validator.StructValidator.Validate(s)
}

func (s *Product) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name": s.Name,
	}
}
