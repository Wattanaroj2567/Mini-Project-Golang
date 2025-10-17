package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gamegear/admin-service/internal/clients"
	"github.com/gamegear/admin-service/internal/models"
)

// ProductService coordinates admin product management through shop-service.
type ProductService interface {
	List(ctx context.Context) ([]models.ProductSummary, error)
	Create(ctx context.Context, req models.ProductUpsertRequest) (*models.ProductSummary, error)
	Update(ctx context.Context, id uint, req models.ProductUpsertRequest) (*models.ProductSummary, error)
	Delete(ctx context.Context, id uint) error
}

type productService struct {
	shopClient clients.ShopServiceClient
}

// NewProductService constructs ProductService.
func NewProductService(shopClient clients.ShopServiceClient) ProductService {
	return &productService{shopClient: shopClient}
}

func (s *productService) List(ctx context.Context) ([]models.ProductSummary, error) {
	resp, err := s.shopClient.GetProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("product service: list products: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("product service: list products: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to fetch products")
		return nil, newServiceError(resp.StatusCode, message)
	}

	// Decode using flexible helper to support multiple upstream shapes
	products, err := decodeProductList(body)
	if err != nil {
		return nil, fmt.Errorf("product service: decode list response: %w", err)
	}

	return products, nil
}

func (s *productService) Create(ctx context.Context, req models.ProductUpsertRequest) (*models.ProductSummary, error) {
	resp, err := s.shopClient.CreateProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("product service: create product: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("product service: create product: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to create product")
		return nil, newServiceError(resp.StatusCode, message)
	}

	product, err := decodeProduct(body)
	if err != nil {
		return nil, fmt.Errorf("product service: decode create response: %w", err)
	}

	// Ensure we never return an empty product
	if product == nil || isEmptyProduct(*product) {
		return nil, newServiceError(http.StatusBadGateway, "upstream returned empty product")
	}

	return product, nil
}

func (s *productService) Update(ctx context.Context, id uint, req models.ProductUpsertRequest) (*models.ProductSummary, error) {
	resp, err := s.shopClient.UpdateProduct(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("product service: update product: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("product service: update product: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to update product")
		return nil, newServiceError(resp.StatusCode, message)
	}

	product, err := decodeProduct(body)
	if err != nil {
		return nil, fmt.Errorf("product service: decode update response: %w", err)
	}

	if product == nil || isEmptyProduct(*product) {
		return nil, newServiceError(http.StatusBadGateway, "upstream returned empty product")
	}

	return product, nil
}

func (s *productService) Delete(ctx context.Context, id uint) error {
	resp, err := s.shopClient.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("product service: delete product: %w", err)
	}
	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		return fmt.Errorf("product service: delete product: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		message := parseErrorMessage(body, "failed to delete product")
		return newServiceError(resp.StatusCode, message)
	}

	return nil
}

func decodeProductList(body []byte) ([]models.ProductSummary, error) {
	var products []models.ProductSummary
	if err := json.Unmarshal(body, &products); err == nil {
		if products == nil {
			return []models.ProductSummary{}, nil
		}
		return products, nil
	}

	var envelope struct {
		Data  []models.ProductSummary `json:"data"`
		Items []models.ProductSummary `json:"items"`
	}

	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	switch {
	case envelope.Items != nil:
		return envelope.Items, nil
	case envelope.Data != nil:
		return envelope.Data, nil
	default:
		return []models.ProductSummary{}, nil
	}
}

func decodeProduct(body []byte) (*models.ProductSummary, error) {
	// Try direct unmarshal first
	if prod, ok := tryUnmarshalProduct(body); ok {
		return prod, nil
	}

	// Try with data wrapper
	var envelope struct {
		Data *models.ProductSummary `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil && envelope.Data != nil {
		if !isEmptyProduct(*envelope.Data) {
			return envelope.Data, nil
		}
	}

	// Try recursive search as last resort
	if prod := searchProductRecursively(body); prod != nil {
		return prod, nil
	}

	return nil, fmt.Errorf("decode product: unexpected payload: %s", string(bytes.TrimSpace(body)))
}

func isEmptyProduct(p models.ProductSummary) bool {
	return p.ID == 0 &&
		p.Name == "" &&
		p.Description == "" &&
		p.Price == 0 &&
		p.Stock == 0 &&
		p.CategoryID == 0 &&
		p.ImageURL == ""
}

func tryUnmarshalProduct(data []byte) (*models.ProductSummary, bool) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil, false
	}

	var product models.ProductSummary
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, false
	}

	if isEmptyProduct(product) {
		return nil, false
	}
	return &product, true
}

func searchProductRecursively(data []byte) *models.ProductSummary {
	if prod, ok := tryUnmarshalProduct(data); ok {
		return prod
	}

	if prod, ok := tryCoerceProductFromMap(data); ok {
		return prod
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err == nil {
		for _, value := range object {
			if prod := searchProductRecursively(value); prod != nil {
				return prod
			}
		}
		return nil
	}

	var array []json.RawMessage
	if err := json.Unmarshal(data, &array); err == nil {
		for _, value := range array {
			if prod := searchProductRecursively(value); prod != nil {
				return prod
			}
		}
	}

	return nil
}

func tryCoerceProductFromMap(data []byte) (*models.ProductSummary, bool) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, false
	}

	if len(obj) == 0 {
		return nil, false
	}

	product := models.ProductSummary{
		ID:          coerceUint(obj["id"]),
		Name:        coerceString(obj["name"]),
		Description: coerceString(obj["description"]),
		Price:       coerceFloat(obj["price"]),
		Stock:       coerceInt(obj["stock"]),
		CategoryID:  coerceUint(obj["category_id"]),
		ImageURL:    coerceString(obj["image_url"]),
	}

	if isEmptyProduct(product) {
		return nil, false
	}

	return &product, true
}

func coerceString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case fmt.Stringer:
		return strings.TrimSpace(v.String())
	case float64:
		if v == 0 {
			return ""
		}
		return strings.TrimSpace(strconv.FormatFloat(v, 'f', -1, 64))
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		return ""
	}
}

func coerceUint(value interface{}) uint {
	switch v := value.(type) {
	case float64:
		if v < 0 {
			return 0
		}
		return uint(v)
	case string:
		v = strings.TrimSpace(v)
		if v == "" {
			return 0
		}
		if parsed, err := strconv.ParseUint(v, 10, 64); err == nil {
			return uint(parsed)
		}
	case json.Number:
		if parsed, err := v.Int64(); err == nil {
			if parsed < 0 {
				return 0
			}
			return uint(parsed)
		}
	case int:
		if v < 0 {
			return 0
		}
		return uint(v)
	case int64:
		if v < 0 {
			return 0
		}
		return uint(v)
	}
	return 0
}

func coerceInt(value interface{}) int {
	switch v := value.(type) {
	case float64:
		return int(v)
	case string:
		v = strings.TrimSpace(v)
		if v == "" {
			return 0
		}
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			return int(parsed)
		}
	case json.Number:
		if parsed, err := v.Int64(); err == nil {
			return int(parsed)
		}
	case int:
		return v
	case int64:
		return int(v)
	}
	return 0
}

func coerceFloat(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case string:
		v = strings.TrimSpace(v)
		if v == "" {
			return 0
		}
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			return parsed
		}
	case json.Number:
		if parsed, err := v.Float64(); err == nil {
			return parsed
		}
	case int:
		return float64(v)
	case int64:
		return float64(v)
	}
	return 0
}
