package models

import "time"

// =========================
// USERS
// =========================
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:120;not null" json:"name"`
	Email     string    `gorm:"size:191;not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"password_hash"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Interactions    []UserInteraction `gorm:"foreignKey:UserID" json:"interactions,omitempty" swaggerignore:"true"`
	Recommendations []Recommendation  `gorm:"foreignKey:UserID" json:"recommendations,omitempty" swaggerignore:"true"`
}
