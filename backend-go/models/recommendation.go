package models

import (
	"time"

	"gorm.io/datatypes"
)

// =========================
// RECOMMENDATIONS
// =========================
type Recommendation struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	UserID                uint           `json:"user_id"`
	RecommendedContentIDs datatypes.JSON `json:"recommended_content_ids" swaggertype:"string"`
	ModelVersion          string         `json:"model_version"`
	CreatedAt             time.Time      `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty" swaggerignore:"true"`
}
