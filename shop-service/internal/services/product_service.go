package services

import (
	"context"

	"github.com/gamegear/shop-service/internal/models"
	"github.com/gamegear/shop-service/internal/repositories"
)

// ProductService orchestrates catalogue related use cases.
type ProductService interface {
	List(ctx context.Context, filter models.ProductFilter) ([]models.ProductResponse, int64, error)
	GetByID(ctx context.Context, id uint) (*models.ProductResponse, error)
	Create(ctx context.Context, req models.ProductRequest) (*models.ProductResponse, error)
	Update(ctx context.Context, id uint, req models.ProductRequest) (*models.ProductResponse, error)
	Delete(ctx context.Context, id uint) error
}

type productService struct {
	productRepo repositories.ProductRepository
}

// NewProductService constructs ProductService.
func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) List(ctx context.Context, filter models.ProductFilter) ([]models.ProductResponse, int64, error) {
	// Get products from repository
	products, err := s.productRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	totalCount, err := s.productRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Map products to responses
	responses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = mapProductToResponse(&product)
	}

	return responses, totalCount, nil
}

func (s *productService) GetByID(ctx context.Context, id uint) (*models.ProductResponse, error) {
	// Get product from repository
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Map to response
	response := mapProductToResponse(product)
	return &response, nil
}

func (s *productService) Create(ctx context.Context, req models.ProductRequest) (*models.ProductResponse, error) {
	// Map request to product entity
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		ImageURL:    req.ImageURL,
	}

	// Save to repository
	if err := s.productRepo.Create(ctx, &product); err != nil {
		return nil, err
	}

	// Map to response
	response := mapProductToResponse(&product)
	return &response, nil
}

func (s *productService) Update(ctx context.Context, id uint, req models.ProductRequest) (*models.ProductResponse, error) {
	// Find existing product
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.CategoryID = req.CategoryID
	product.ImageURL = req.ImageURL

	// Save to repository
	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	// Map to response
	response := mapProductToResponse(product)
	return &response, nil
}

func (s *productService) Delete(ctx context.Context, id uint) error {
	// Delete from repository
	return s.productRepo.Delete(ctx, id)
}

// mapProductToResponse converts Product entity to ProductResponse DTO.
func mapProductToResponse(product *models.Product) models.ProductResponse {
	return models.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		ImageURL:    product.ImageURL,
	}
}
