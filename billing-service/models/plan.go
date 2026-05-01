package models

import "time"

type Plan struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"size:120;not null"`
	Description  string    `json:"description" gorm:"size:500"`
	Price        float64   `json:"price" gorm:"not null;check:price >= 0"`
	Currency     string    `json:"currency" gorm:"size:10;not null;default:USD"`
	BillingCycle string    `json:"billing_cycle" gorm:"size:30;not null;default:monthly"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
