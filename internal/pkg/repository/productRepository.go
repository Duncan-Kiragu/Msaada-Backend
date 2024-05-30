package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/domain"
	"github.com/Duncan-Kiragu/Msaada-Backend/internal/pkg/dto"
	"github.com/Duncan-Kiragu/Msaada-Backend/pkg/filter"
)

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

type productRepository struct {
	db *gorm.DB
}

func (s *productRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	db := s.db.WithContext(ctx)
	db = filter.ApplySearchLike(db, "name")

	return filter.ApplyOrder(db)
}

func (s *productRepository) CountProducts(ctx context.Context, filter *filter.Filter) (int64, error) {
	var count int64
	db := s.applyFilter(ctx, filter)
	return count, db.Model(&domain.Product{}).Count(&count).Error
}

func (s *productRepository) GetProducts(ctx context.Context, filter *filter.Filter) (*[]domain.Product, error) {
	db := filter.ApplyPagination(s.applyFilter(ctx, filter))

	products := &[]domain.Product{}
	return products, db.Find(products).Error
}

func (s *productRepository) GetProductByID(ctx context.Context, productID uint) (*domain.Product, error) {
	product := &domain.Product{}
	return product, s.db.WithContext(ctx).First(product, productID).Error
}

func (s *productRepository) CreateProduct(ctx context.Context, data *dto.ProductInputDTO) (*domain.Product, error) {
	product := &domain.Product{}
	if err := product.Bind(data); err != nil {
		return nil, err
	}

	return product, s.db.WithContext(ctx).Create(product).Error
}

func (s *productRepository) UpdateProduct(ctx context.Context, product *domain.Product, data *dto.ProductInputDTO) error {
	if err := product.Bind(data); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(product).Updates(product.ToMap()).Error
}

func (s *productRepository) DeleteProduct(ctx context.Context, product *domain.Product) error {
	return s.db.WithContext(ctx).Delete(product).Error
}
