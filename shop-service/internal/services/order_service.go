package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/gamegear/shop-service/internal/models"
	"github.com/gamegear/shop-service/internal/repositories"
	"gorm.io/gorm"
)

var (
	// ErrOrderNotFound indicates the requested order does not exist.
	ErrOrderNotFound = errors.New("order not found")
	// ErrInvalidOrderStatus indicates the supplied status is not supported.
	ErrInvalidOrderStatus = errors.New("invalid order status")
)

// OrderService coordinates checkout lifecycle.
type OrderService interface {
	Create(ctx context.Context, userID uint, req models.CreateOrderRequest) (*models.OrderDetail, error)
	History(ctx context.Context, userID uint) ([]models.OrderSummary, error)
	GetByID(ctx context.Context, userID uint, orderID uint) (*models.OrderDetail, error)
	ListAll(ctx context.Context) ([]models.AdminOrderSummary, error)
	UpdateStatus(ctx context.Context, orderID uint, status string) (*models.OrderDetail, error)
}

type orderService struct {
	orderRepo    repositories.OrderRepository
	cartRepo     repositories.CartRepository
	cartItemRepo repositories.CartItemRepository
}

// NewOrderService constructs OrderService.
func NewOrderService(orderRepo repositories.OrderRepository, cartRepo repositories.CartRepository, cartItemRepo repositories.CartItemRepository) OrderService {
	return &orderService{
		orderRepo:    orderRepo,
		cartRepo:     cartRepo,
		cartItemRepo: cartItemRepo,
	}
}

func (s *orderService) Create(ctx context.Context, userID uint, req models.CreateOrderRequest) (*models.OrderDetail, error) {
	// 1. Load the active cart
	cart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cart == nil || len(cart.CartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	// 2. Create Order and OrderItems from Cart
	order := &models.Order{
		UserID:          userID,
		CartID:          cart.ID,   // Link order to the cart
		Status:          "pending", // Initial status
		Total:           cart.Total,
		ShippingAddress: req.ShippingAddress,
		PaymentMethod:   req.PaymentMethod,
		Notes:           req.Notes,
	}

	for _, item := range cart.CartItems {
		order.OrderItems = append(order.OrderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price, // Price at the time of purchase
		})
	}

	// 3. Persist order in a transaction
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// 4. Clear the cart items
	if err := s.cartItemRepo.RemoveAllItemsByCartID(ctx, cart.ID); err != nil {
		// Log this error, but don't fail the order creation
		fmt.Printf("Warning: could not clear cart items: %v\n", err)
	}

	// 5. Update cart totals to 0
	if err := s.cartRepo.UpdateTotals(ctx, cart.ID); err != nil {
		fmt.Printf("Warning: could not update cart totals: %v\n", err)
	}

	// 6. Return the detailed order
	return s.GetByID(ctx, userID, order.ID)
}

func (s *orderService) History(ctx context.Context, userID uint) ([]models.OrderSummary, error) {
	orders, err := s.orderRepo.FindHistoryByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Map to summary DTO
	summaries := make([]models.OrderSummary, len(orders))
	for i, order := range orders {
		summaries[i] = models.OrderSummary{
			ID:            order.ID,
			Status:        order.Status,
			Total:         order.Total,
			PaymentMethod: order.PaymentMethod,
		}
	}
	return summaries, nil
}

func (s *orderService) GetByID(ctx context.Context, userID uint, orderID uint) (*models.OrderDetail, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Security check: ensure the order belongs to the user
	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return mapOrderToDetail(order), nil
}

func (s *orderService) ListAll(ctx context.Context) ([]models.AdminOrderSummary, error) {
	orders, err := s.orderRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	summaries := make([]models.AdminOrderSummary, len(orders))
	for i, order := range orders {
		summaries[i] = models.AdminOrderSummary{
			ID:              order.ID,
			UserID:          order.UserID,
			Status:          order.Status,
			Total:           order.Total,
			PaymentMethod:   order.PaymentMethod,
			ShippingAddress: order.ShippingAddress,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
		}
	}

	return summaries, nil
}

func (s *orderService) UpdateStatus(ctx context.Context, orderID uint, status string) (*models.OrderDetail, error) {
	if !models.IsValidOrderStatus(status) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidOrderStatus, status)
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, status); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return mapOrderToDetail(order), nil
}

func mapOrderToDetail(order *models.Order) *models.OrderDetail {
	// Map to detail DTO
	items := make([]models.OrderItemDetail, len(order.OrderItems))
	for i, item := range order.OrderItems {
		items[i] = models.OrderItemDetail{
			OrderItemID: item.ID,
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Product: &models.ProductSummary{
				ID:       item.Product.ID,
				Name:     item.Product.Name,
				Price:    item.Product.Price,
				ImageURL: item.Product.ImageURL,
			},
		}
	}

	return &models.OrderDetail{
		OrderSummary: models.OrderSummary{
			ID:            order.ID,
			Status:        order.Status,
			Total:         order.Total,
			PaymentMethod: order.PaymentMethod,
		},
		ShippingAddress: order.ShippingAddress,
		Notes:           order.Notes,
		Items:           items,
	}
}
