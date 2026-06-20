package models

import "time"

type Documentation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SystemId  uint      `json:"system_id"`
	Title     string    `gorm:"type:varchar(255)" json:"title"`
	Category  string    `gorm:"type:varchar(100)" json:"category"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateDocumentationRequest struct {
	Title    string `json:"title" binding:"required"`
	Category string `json:"category" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type UpdateDocumentationRequest struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

// IsValidDocumentationCategory checks if the given category string is one of the allowed categories
func IsValidDocumentationCategory(category string) bool {
	switch category {
	case "Business Flow",
		"Technical Flow",
		"API Documentation",
		"Database Documentation",
		"Deployment Guide",
		"User Manual":
		return true
	}
	return false
}
