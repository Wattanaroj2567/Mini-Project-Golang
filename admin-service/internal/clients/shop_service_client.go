package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ShopServiceClient interface {
	GetProducts(ctx context.Context) (*http.Response, error)
	GetProduct(ctx context.Context, id uint) (*http.Response, error)
	CreateProduct(ctx context.Context, productData interface{}) (*http.Response, error)
	UpdateProduct(ctx context.Context, id uint, productData interface{}) (*http.Response, error)
	DeleteProduct(ctx context.Context, id uint) (*http.Response, error)
	GetOrders(ctx context.Context) (*http.Response, error)
	GetOrder(ctx context.Context, id uint) (*http.Response, error)
	UpdateOrderStatus(ctx context.Context, id uint, status string) (*http.Response, error)
}

type shopServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewShopServiceClient(baseURL string) ShopServiceClient {
	return &shopServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetProducts - ดึงข้อมูลสินค้าทั้งหมด
func (c *shopServiceClient) GetProducts(ctx context.Context) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, "/api/products", nil)
}

// GetProduct - ดึงข้อมูลสินค้าตาม ID
func (c *shopServiceClient) GetProduct(ctx context.Context, id uint) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, fmt.Sprintf("/api/products/%d", id), nil)
}

// CreateProduct - สร้างสินค้าใหม่
func (c *shopServiceClient) CreateProduct(ctx context.Context, productData interface{}) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, "/api/products", productData)
}

// UpdateProduct - อัปเดตข้อมูลสินค้า
func (c *shopServiceClient) UpdateProduct(ctx context.Context, id uint, productData interface{}) (*http.Response, error) {
	return c.do(ctx, http.MethodPut, fmt.Sprintf("/api/products/%d", id), productData)
}

// DeleteProduct - ลบสินค้า
func (c *shopServiceClient) DeleteProduct(ctx context.Context, id uint) (*http.Response, error) {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/products/%d", id), nil)
}

// GetOrders - ดึงข้อมูลคำสั่งซื้อทั้งหมด
func (c *shopServiceClient) GetOrders(ctx context.Context) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, "/api/orders", nil)
}

// GetOrder - ดึงข้อมูลคำสั่งซื้อตาม ID
func (c *shopServiceClient) GetOrder(ctx context.Context, id uint) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, fmt.Sprintf("/api/orders/%d", id), nil)
}

// UpdateOrderStatus - อัปเดตสถานะคำสั่งซื้อ
func (c *shopServiceClient) UpdateOrderStatus(ctx context.Context, id uint, status string) (*http.Response, error) {
	payload := map[string]string{"status": status}
	return c.do(ctx, http.MethodPut, fmt.Sprintf("/api/orders/%d/status", id), payload)
}

func (c *shopServiceClient) do(ctx context.Context, method, path string, payload interface{}) (*http.Response, error) {
	if strings.TrimSpace(c.baseURL) == "" {
		return nil, fmt.Errorf("shop client: base URL is not configured")
	}

	url := strings.TrimRight(c.baseURL, "/") + path

	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("shop client: marshal payload: %w", err)
		}

		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("shop client: create request: %w", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if authHeader := strings.TrimSpace(AuthTokenFromContext(ctx)); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("shop client: perform request: %w", err)
	}

	return resp, nil
}
