package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:120;not null"`
	Email     string    `json:"email" gorm:"size:255;not null;uniqueIndex"`
	Password  string    `json:"-" gorm:"size:255;not null;default:''"`
	PlanID    *uint     `json:"plan_id"`
	Plan      *Plan     `json:"plan,omitempty"`
	IsActive  bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
