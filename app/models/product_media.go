package models

import "time"

type ProductMedia struct {
	ID          string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	OwnerID     string    `gorm:"index;type:varchar(64)" json:"-"`
	URL         string    `gorm:"type:text" json:"url"`
	Kind        string    `gorm:"type:varchar(24)" json:"kind"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	OptionsJSON string    `gorm:"type:text" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ProductMedia) TableName() string { return "product_media" }
