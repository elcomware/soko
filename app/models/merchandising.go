package models

type MerchandisingRecord struct {
	ID           string `gorm:"primaryKey;type:varchar(32)"`
	SettingsJSON string
	Version      int
}

func (MerchandisingRecord) TableName() string { return "storefront_merchandising" }

type MerchandisingPlacement struct {
	Enabled    bool     `json:"enabled"`
	Title      string   `json:"title"`
	Strategy   string   `json:"strategy"`
	ProductIDs []string `json:"product_ids"`
	Category   string   `json:"category"`
	Fallback   string   `json:"fallback"`
	Limit      int      `json:"limit"`
}
type CategoryNavigation struct {
	Enabled bool   `json:"enabled"`
	Title   string `json:"title"`
}
type MerchandisingSettings struct {
	SectionOrder       []string               `json:"section_order"`
	Announcement       string                 `json:"announcement"`
	Version            int                    `json:"version"`
	Hero               MerchandisingPlacement `json:"hero"`
	HeroSide           MerchandisingPlacement `json:"hero_side"`
	Discovery          MerchandisingPlacement `json:"discovery"`
	CategoryNavigation CategoryNavigation     `json:"category_navigation"`
}
