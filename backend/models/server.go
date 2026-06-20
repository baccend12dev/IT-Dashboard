package models

import (
	"time"
)

type Server struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	IP        string    `gorm:"type:varchar(255)" json:"ip"`
	OS        string    `gorm:"type:varchar(255)" json:"os"`
	Location  string    `gorm:"type:varchar(255)" json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Server represents the structure of a server in the database.
type CreateServerRequest struct {
	Name     string `json:"name" binding:"required"`
	IP       string `json:"ip" binding:"required"`
	OS       string `json:"os" binding:"required"`
	Location string `json:"location" binding:"required"`
}
