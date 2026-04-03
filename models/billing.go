package models

import "time"

const (
	BillingStatusPending = "pending"
	BillingStatusPaid    = "paid"
	BillingStatusFailed  = "failed"
)

type Billing struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	UserID      uint       `json:"user_id" gorm:"not null;index"`
	User        User       `json:"user,omitempty"`
	PlanID      uint       `json:"plan_id" gorm:"not null;index"`
	Plan        Plan       `json:"plan,omitempty"`
	Amount      float64    `json:"amount" gorm:"not null;check:amount >= 0"`
	Status      string     `json:"status" gorm:"size:20;not null;default:pending"`
	DueDate     time.Time  `json:"due_date" gorm:"not null"`
	PaidAt      *time.Time `json:"paid_at"`
	Description string     `json:"description" gorm:"size:500"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
