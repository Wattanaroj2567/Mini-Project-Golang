package services

import (
	"context"
	"errors"

	"github.com/gamegear/shop-service/internal/models"
	"github.com/gamegear/shop-service/internal/repositories"
)

// CartService manages cart operations for members.
type CartService interface {
	GetCart(ctx context.Context, userID uint) (*models.CartResponse, error)
	AddItem(ctx context.Context, userID uint, req models.AddToCartRequest) (*models.CartResponse, error)
	UpdateItem(ctx context.Context, userID uint, req models.UpdateCartItemRequest) (*models.CartResponse, error)
	RemoveItem(ctx context.Context, userID uint, req models.RemoveCartItemRequest) (*models.CartResponse, error)
}

type cartService struct {
	cartRepo     repositories.CartRepository
	cartItemRepo repositories.CartItemRepository
	productRepo  repositories.ProductRepository
}

// NewCartService constructs CartService.
func NewCartService(
	cartRepo repositories.CartRepository,
	cartItemRepo repositories.CartItemRepository,
	productRepo repositories.ProductRepository,
) CartService {
	return &cartService{
		cartRepo:     cartRepo,
		cartItemRepo: cartItemRepo,
		productRepo:  productRepo,
	}
}

func (s *cartService) GetCart(ctx context.Context, userID uint) (*models.CartResponse, error) {
	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapCartToResponse(cart), nil
}

func (s *cartService) AddItem(ctx context.Context, userID uint, req models.AddToCartRequest) (*models.CartResponse, error) {
	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Validate product stock
	product, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}
	if product.Stock < req.Quantity {
		return nil, errors.New("not enough stock")
	}

	cartItem := &models.CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := s.cartItemRepo.AddItem(ctx, cartItem); err != nil {
		return nil, err
	}

	if err := s.cartRepo.UpdateTotals(ctx, cart.ID); err != nil {
		return nil, err
	}

	updatedCart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapCartToResponse(updatedCart), nil
}

func (s *cartService) UpdateItem(ctx context.Context, userID uint, req models.UpdateCartItemRequest) (*models.CartResponse, error) {
	if err := s.cartItemRepo.UpdateQuantity(ctx, req.CartItemID, req.Quantity); err != nil {
		return nil, err
	}

	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.cartRepo.UpdateTotals(ctx, cart.ID); err != nil {
		return nil, err
	}

	updatedCart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapCartToResponse(updatedCart), nil
}

func (s *cartService) RemoveItem(ctx context.Context, userID uint, req models.RemoveCartItemRequest) (*models.CartResponse, error) {
	if err := s.cartItemRepo.RemoveItem(ctx, req.CartItemID); err != nil {
		return nil, err
	}

	cart, err := s.getOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.cartRepo.UpdateTotals(ctx, cart.ID); err != nil {
		return nil, err
	}

	updatedCart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapCartToResponse(updatedCart), nil
}

// Helper to get or create a cart
func (s *cartService) getOrCreateCart(ctx context.Context, userID uint) (*models.Cart, error) {
	cart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		newCart := &models.Cart{UserID: userID}
		if err := s.cartRepo.Create(ctx, newCart); err != nil {
			return nil, err
		}
		return newCart, nil
	}
	return cart, nil
}

// Helper to map cart to response DTO
func mapCartToResponse(cart *models.Cart) *models.CartResponse {
	items := make([]models.CartItemSummary, len(cart.CartItems))
	for i, item := range cart.CartItems {
		items[i] = models.CartItemSummary{
			CartItemID: item.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			Product: &models.ProductSummary{
				ID:       item.Product.ID,
				Name:     item.Product.Name,
				Price:    item.Product.Price,
				ImageURL: item.Product.ImageURL,
			},
		}
	}

	return &models.CartResponse{
		ID:       cart.ID,
		UserID:   cart.UserID,
		Items:    items,
		Subtotal: cart.Subtotal,
		Total:    cart.Total,
	}
}
