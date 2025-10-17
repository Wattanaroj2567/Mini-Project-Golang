package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gamegear/admin-service/internal/clients"
	"github.com/gamegear/admin-service/internal/models"
)

// OrderService coordinates order management through shop-service.
type OrderService interface {
	List(ctx context.Context) ([]models.OrderSummary, error)
	GetByID(ctx context.Context, id uint) (*models.OrderDetail, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
}

type orderService struct {
	shopClient clients.ShopServiceClient
}

// NewOrderService constructs OrderService.
func NewOrderService(shopClient clients.ShopServiceClient) OrderService {
	return &orderService{shopClient: shopClient}
}

func (s *orderService) List(ctx context.Context) ([]models.OrderSummary, error) {
	resp, err := s.shopClient.GetOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("order service: list orders: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("order service: list orders: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to fetch orders")
		return nil, newServiceError(resp.StatusCode, message)
	}

	upstreamOrders, err := decodeOrderSummaries(body)
	if err != nil {
		return nil, fmt.Errorf("order service: decode list response: %w", err)
	}

	summaries := make([]models.OrderSummary, len(upstreamOrders))
	for i, order := range upstreamOrders {
		summaries[i] = mapOrderSummary(order)
	}

	return summaries, nil
}

func (s *orderService) GetByID(ctx context.Context, id uint) (*models.OrderDetail, error) {
	resp, err := s.shopClient.GetOrder(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("order service: get order: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("order service: get order: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to fetch order detail")
		return nil, newServiceError(resp.StatusCode, message)
	}

	order, err := decodeOrderDetail(body)
	if err != nil {
		return nil, fmt.Errorf("order service: decode detail response: %w", err)
	}

	detail := mapOrderDetail(*order)
	return &detail, nil
}

func (s *orderService) UpdateStatus(ctx context.Context, id uint, status string) error {
	resp, err := s.shopClient.UpdateOrderStatus(ctx, id, status)
	if err != nil {
		return fmt.Errorf("order service: update status: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return fmt.Errorf("order service: update status: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to update order status")
		return newServiceError(resp.StatusCode, message)
	}

	return nil
}

type orderSummaryPayload struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	Total         float64 `json:"total"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
}

type orderItemPayload struct {
	OrderItemID uint    `json:"order_item_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Product     *struct {
		Name string `json:"name"`
	} `json:"product"`
}

type orderDetailPayload struct {
	orderSummaryPayload
	ShippingAddress string             `json:"shipping_address"`
	Notes           string             `json:"notes"`
	Items           []orderItemPayload `json:"items"`
}

func decodeOrderSummaries(body []byte) ([]orderSummaryPayload, error) {
	var orders []orderSummaryPayload
	if err := json.Unmarshal(body, &orders); err == nil {
		return orders, nil
	}

	var envelope struct {
		Data []orderSummaryPayload `json:"data"`
	}

	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	return envelope.Data, nil
}

func decodeOrderDetail(body []byte) (*orderDetailPayload, error) {
	var order orderDetailPayload
	if err := json.Unmarshal(body, &order); err == nil {
		return &order, nil
	}

	var envelope struct {
		Data orderDetailPayload `json:"data"`
	}

	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	return &envelope.Data, nil
}

func mapOrderSummary(order orderSummaryPayload) models.OrderSummary {
	return models.OrderSummary{
		ID:            order.ID,
		UserID:        order.UserID,
		Total:         order.Total,
		Status:        order.Status,
		PaymentMethod: order.PaymentMethod,
	}
}

func mapOrderDetail(order orderDetailPayload) models.OrderDetail {
	items := make([]models.OrderItemSummary, len(order.Items))
	for i, item := range order.Items {
		items[i] = models.OrderItemSummary{
			OrderItemID: item.OrderItemID,
			ProductID:   item.ProductID,
			ProductName: extractProductName(item),
			Quantity:    item.Quantity,
			Price:       item.Price,
		}
	}

	return models.OrderDetail{
		OrderSummary: models.OrderSummary{
			ID:            order.ID,
			UserID:        order.UserID,
			Total:         order.Total,
			Status:        order.Status,
			PaymentMethod: order.PaymentMethod,
		},
		ShippingAddress: order.ShippingAddress,
		Notes:           order.Notes,
		Items:           items,
	}
}

func extractProductName(item orderItemPayload) string {
	if strings.TrimSpace(item.ProductName) != "" {
		return item.ProductName
	}

	if item.Product != nil && item.Product.Name != "" {
		return item.Product.Name
	}

	return ""
}
