package models

import (
	"time"
)

type Note struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SystemId  uint      `json:"system_id"`
	Title     string    `gorm:"type:varchar(255)" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Note represents the structure of a note in the database.
// It includes fields such as ID, SystemId, Title, Content, CreatedAt, and UpdatedAt.
type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
