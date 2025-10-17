package repositories

import (
	"context"

	"github.com/gamegear/shop-service/internal/models"
	"gorm.io/gorm"
)

// OrderRepository persists orders and order items.
type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	FindHistoryByUserID(ctx context.Context, userID uint) ([]models.Order, error)
	FindByID(ctx context.Context, orderID uint) (*models.Order, error)
	FindAll(ctx context.Context) ([]models.Order, error)
	UpdateStatus(ctx context.Context, orderID uint, status string) error
}

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository constructs OrderRepository.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// The Items should be associated with the order and created automatically
		// if the relationship is set up correctly in the GORM model.
		return nil
	})
}

func (r *orderRepository) FindHistoryByUserID(ctx context.Context, userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByID(ctx context.Context, orderID uint) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems.Product").
		First(&order, orderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindAll(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.WithContext(ctx).
		Order("created_at desc").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID uint, status string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
