package models

import "time"

// StorefrontContent stores the editable, plain-text copy for the marketplace.
type StorefrontContent struct {
	ID          string    `gorm:"primaryKey;type:varchar(32)" json:"id"`
	ContentJSON string    `gorm:"column:content_json;type:text" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (StorefrontContent) TableName() string { return "storefront_contents" }
