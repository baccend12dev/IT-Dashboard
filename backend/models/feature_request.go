package models

import "time"

type FeatureRequest struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SystemId    uint      `json:"system_id"`
	System      *System   `gorm:"foreignKey:SystemId" json:"system,omitempty"`
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(50);default:'Pending'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateFeatureRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateFeatureRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
