package service

import (
	"context"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/filter"
)

func NewProductService(r domain.ProductRepository) domain.ProductService {
	return &productService{
		productRepository: r,
	}
}

type productService struct {
	productRepository domain.ProductRepository
}

func (s *productService) generateProductOutputDTO(product *domain.Product) *dto.ProductOutputDTO {
	return &dto.ProductOutputDTO{
		Id:   product.Id,
		Name: product.Name,
	}
}

// GetProductByID Implementation of 'GetProductByID'.
func (s *productService) GetProductByID(ctx context.Context, productID uint) (*dto.ProductOutputDTO, error) {
	product, err := s.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return s.generateProductOutputDTO(product), nil
}

// GetProducts Implementation of 'GetProducts'.
func (s *productService) GetProducts(ctx context.Context, filter *filter.Filter) (*dto.ListItemsOutputDTO, error) {
	count, err := s.productRepository.CountProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	products, err := s.productRepository.GetProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	outputProducts := &[]dto.ProductOutputDTO{}
	for _, product := range *products {
		*outputProducts = append(*outputProducts, *s.generateProductOutputDTO(&product))
	}

	return &dto.ListItemsOutputDTO{
		Count: count,
		Items: outputProducts,
	}, nil
}

// CreateProduct Implementation of 'CreateProduct'.
func (s *productService) CreateProduct(ctx context.Context, data *dto.ProductInputDTO) (*dto.ProductOutputDTO, error) {
	product, err := s.productRepository.CreateProduct(ctx, data)
	if err != nil {
		return nil, err
	}

	return s.generateProductOutputDTO(product), nil
}

// UpdateProduct Implementation of 'UpdateProduct'.
func (s *productService) UpdateProduct(ctx context.Context, product *domain.Product, data *dto.ProductInputDTO) (*dto.ProductOutputDTO, error) {
	if err := s.productRepository.UpdateProduct(ctx, product, data); err != nil {
		return nil, err
	}

	return s.generateProductOutputDTO(product), nil
}

// DeleteProduct Implementation of 'DeleteProduct'.
func (s *productService) DeleteProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepository.DeleteProduct(ctx, product)
}
