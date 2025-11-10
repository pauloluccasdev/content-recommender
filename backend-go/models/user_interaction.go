package models

import "time"

// =========================
// USER_INTERACTIONS
// =========================
type UserInteraction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `json:"user_id"`
	ContentID       uint      `json:"content_id"`
	InteractionType string    `json:"interaction_type"`
	Rating          *float64  `json:"rating,omitempty"`
	CreatedAt       time.Time `json:"created_at"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content Content `gorm:"foreignKey:ContentID" json:"content,omitempty"`
}
