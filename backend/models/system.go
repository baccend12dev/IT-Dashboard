package models

import (
	"time"
)

type System struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	Name            string           `gorm:"type:varchar(255)" json:"name"`
	Type            string           `gorm:"type:varchar(255)" json:"type"`
	Links           string           `gorm:"type:varchar(255)" json:"links"`
	ServerId        uint             `json:"server_id"`
	Server          Server           `gorm:"foreignKey:ServerId" json:"server,omitempty"`
	Status          string           `gorm:"type:varchar(255)" json:"status"`
	Description     string           `gorm:"type:varchar(255)" json:"description"`
	Notes           []Note           `gorm:"foreignKey:SystemId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"notes,omitempty"`
	FeatureRequests []FeatureRequest `gorm:"foreignKey:SystemId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"feature_requests,omitempty"`
	Documentations  []Documentation  `gorm:"foreignKey:SystemId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"documentations,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// System represents the structure of a system in the database.
// It includes fields such as ID, Name, Type, Links, ServerId, Status, Description, CreatedAt, and UpdatedAt.
type CreateSystemRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Links       string `json:"links" binding:"required"`
	ServerId    uint   `json:"server_id" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreateSystemRequest represents the expected structure of the request body when creating a new system.
// It includes fields such as Name, Type, Links, ServerId, Status, and Description, all of which are required.
