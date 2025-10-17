package repositories

import (
	"context"

	"github.com/gamegear/shop-service/internal/models"
	"gorm.io/gorm"
)

// ProductRepository handles persistence for catalogue data.
type ProductRepository interface {
	List(ctx context.Context, filter models.ProductFilter) ([]models.Product, error)
	Count(ctx context.Context, filter models.ProductFilter) (int64, error)
	FindByID(ctx context.Context, id uint) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository wires a GORM-backed repository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) List(ctx context.Context, filter models.ProductFilter) ([]models.Product, error) {
	var products []models.Product

	query := r.db.WithContext(ctx)

	// Apply category filter
	if filter.Category != "" {
		query = query.Where("category_id = ?", filter.Category)
	}

	// Apply search filter (search in name and description)
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit)

	// Order by ID descending (newest first)
	query = query.Order("id DESC")

	// Execute query
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) Count(ctx context.Context, filter models.ProductFilter) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&models.Product{})

	// Apply category filter
	if filter.Category != "" {
		query = query.Where("category_id = ?", filter.Category)
	}

	// Apply search filter (search in name and description)
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Count records
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product

	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	result := r.db.WithContext(ctx).Create(product)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrInvalidData
	}

	return nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	result := r.db.WithContext(ctx).Save(product)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Product{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
