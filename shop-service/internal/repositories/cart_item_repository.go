package repositories

import (
	"context"

	"github.com/gamegear/shop-service/internal/models"
	"gorm.io/gorm"
)

// CartItemRepository manages individual cart items.
type CartItemRepository interface {
	AddItem(ctx context.Context, item *models.CartItem) error
	UpdateQuantity(ctx context.Context, itemID uint, quantity int) error
	RemoveItem(ctx context.Context, itemID uint) error
	RemoveAllItemsByCartID(ctx context.Context, cartID uint) error
}

type cartItemRepository struct {
	db *gorm.DB
}

// NewCartItemRepository constructs CartItemRepository.
func NewCartItemRepository(db *gorm.DB) CartItemRepository {
	return &cartItemRepository{db: db}
}

func (r *cartItemRepository) AddItem(ctx context.Context, item *models.CartItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if the item already exists in the cart
		var existingItem models.CartItem
		err := tx.Where("cart_id = ? AND product_id = ?", item.CartID, item.ProductID).First(&existingItem).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if existingItem.ID > 0 {
			// If item exists, update the quantity
			existingItem.Quantity += item.Quantity
			return tx.Save(&existingItem).Error
		} else {
			// Otherwise, create a new cart item
			return tx.Create(item).Error
		}
	})
}

func (r *cartItemRepository) UpdateQuantity(ctx context.Context, itemID uint, quantity int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item models.CartItem
		if err := tx.First(&item, itemID).Error; err != nil {
			return err
		}

		item.Quantity = quantity
		return tx.Save(&item).Error
	})
}

func (r *cartItemRepository) RemoveItem(ctx context.Context, itemID uint) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, itemID).Error
}

func (r *cartItemRepository) RemoveAllItemsByCartID(ctx context.Context, cartID uint) error {
	return r.db.WithContext(ctx).Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}
