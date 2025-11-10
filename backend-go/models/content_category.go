package models

// =========================
// CONTENT_CATEGORIES
// =========================
type ContentCategory struct {
	ContentID  uint `gorm:"primaryKey" json:"content_id"`
	CategoryID uint `form:"primaryKey" json:"category_id"`

	// Relationships
	Content  Content  `gorm:"foreignKey:ContentID" json:"content,omitempty"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
