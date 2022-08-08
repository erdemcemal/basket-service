package database

import (
	"errors"
	"github.com/erdemcemal/basket-service/internal/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	lowVatRate    = 1.0
	normalVatRate = 8.0
	highVatRate   = 18.0
)

// MigrateDB - migrate our database and creates our comment table
func MigrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Product{}, &models.ShoppingCart{}, &models.ShoppingCartItem{}, &models.SalesHistory{}, &models.SalesHistoryItem{}); err == nil && db.Migrator().HasTable(&models.Product{}) {
		if err := db.First(&models.Product{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&models.Product{Name: "Phone 9", UnitPrice: decimal.New(549, 0), VatRate: normalVatRate, Quantity: 94}).Error; err != nil {
				return err
			}
			if err := db.Create(&models.Product{Name: "MacBook Pro", UnitPrice: decimal.New(1749, 0), VatRate: highVatRate, Quantity: 83}).Error; err != nil {
				return err
			}
			if err := db.Create(&models.Product{Name: "Key Holder", UnitPrice: decimal.New(30, 0), VatRate: lowVatRate, Quantity: 54}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
