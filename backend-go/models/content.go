package models

import "time"

// =========================
// CONTENTS
// =========================
type Content struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	ReleaseDate time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	Interactions []UserInteraction `gorm:"foreignKey:ContentID" json:"user_interactions,omitempty"`
	Categories   []Category        `gorm:"many2many:content_categories" json:"categories,omitempty"`
}
