package models

import (
	"github.com/shopspring/decimal"
)

// Product represents a product in the catalog.
// It includes a unique code and a price.
type Product struct {
	ID         uint            `gorm:"primaryKey"`
	Code       string          `gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	CategoryID uint            `gorm:""`
	Category   Category        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Variants   []Variant       `gorm:"foreignKey:ProductID"`
}

func (p *Product) TableName() string {
	return "products"
}
