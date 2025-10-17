package repositories

import (
	"context"

	"github.com/gamegear/shop-service/internal/models"
	"gorm.io/gorm"
)

// CartRepository manages cart aggregate.



type CartRepository interface {

	FindByUserID(ctx context.Context, userID uint) (*models.Cart, error)

	Create(ctx context.Context, cart *models.Cart) error

	UpdateTotals(ctx context.Context, cartID uint) error

}



type cartRepository struct {

	db *gorm.DB

}



// NewCartRepository constructs CartRepository.

func NewCartRepository(db *gorm.DB) CartRepository {

	return &cartRepository{db: db}

}



func (r *cartRepository) FindByUserID(ctx context.Context, userID uint) (*models.Cart, error) {

	var cart models.Cart

	// Eager load items and their associated products

	err := r.db.WithContext(ctx).

		Preload("CartItems.Product").

		Where("user_id = ?", userID).

		First(&cart).

		Error



	if err != nil {

		if err == gorm.ErrRecordNotFound {

			return nil, nil // Return nil if no active cart is found, this is not an error

		}

		return nil, err

	}

	return &cart, nil

}



func (r *cartRepository) Create(ctx context.Context, cart *models.Cart) error {

	return r.db.WithContext(ctx).Create(cart).Error

}



func (r *cartRepository) UpdateTotals(ctx context.Context, cartID uint) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var cart models.Cart

		if err := tx.Preload("CartItems.Product").First(&cart, cartID).Error; err != nil {

			return err

		}



		var subtotal float64

		for _, item := range cart.CartItems {

			subtotal += item.Product.Price * float64(item.Quantity)

		}



		cart.Subtotal = subtotal

		// Assuming Total is Subtotal + Shipping - Discount. For now, we set it to Subtotal.

		cart.Total = subtotal 



		return tx.Save(&cart).Error

	})

}
