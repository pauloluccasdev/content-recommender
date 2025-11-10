package models

// =========================
// CATEGORIES
// =========================
type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`

	// Relationships
	Contents []Content `gorm:"many2many:content_categories" json:"contents,omitempty"`
}
