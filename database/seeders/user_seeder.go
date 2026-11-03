package seeders

import (
	"time"

	"elcomware/app/facades"
	"elcomware/app/models"
)

type UserSeeder struct{}

// Signature The unique signature of the seeder.
func (s *UserSeeder) Signature() string {
	return "UserSeeder"
}

// Run executes the seeder logic.
func (s *UserSeeder) Run() error {
	hashedPassword, err := facades.Hash().Make("secret123")
	if err != nil {
		hashedPassword = "plain:secret123"
	}

	users := []models.User{
		{
			ID:          "usr_admin_01",
			Name:        "Soko Global Admin",
			Email:       "admin@soko.global",
			Password:    hashedPassword,
			Role:        models.RoleAdmin,
			CompanyName: "Soko Multi-Region Marketplace Operations Inc.",
			Country:     "United States",
			AccountType: "enterprise",
			Verified:    true,
			Avatar:      "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=200&q=80",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "usr_vendor_01",
			Name:        "Kenyan Industrial Milling Co.",
			Email:       "supplier.ea@soko.global",
			Password:    hashedPassword,
			Role:        models.RoleMerchant,
			CompanyName: "Apex Heavy Agro-Tech & Grain Silos Ltd",
			Country:     "Kenya",
			AccountType: "supplier",
			Verified:    true,
			Avatar:      "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=200&q=80",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "usr_vendor_02",
			Name:        "Shenzhen Nova Microelectronics",
			Email:       "merchant.cn@soko.global",
			Password:    hashedPassword,
			Role:        models.RoleMerchant,
			CompanyName: "Nova Silicon & ASIC Fabrication Co.",
			Country:     "China",
			AccountType: "supplier",
			Verified:    true,
			Avatar:      "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=200&q=80",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "usr_buyer_01",
			Name:        "David Mwangi",
			Email:       "buyer@soko.global",
			Password:    hashedPassword,
			Role:        models.RoleBuyer,
			CompanyName: "Great Rift Infrastructure Consortium",
			Country:     "Kenya",
			AccountType: "enterprise",
			Verified:    true,
			Avatar:      "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=200&q=80",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, u := range users {
		var existing models.User
		if err := facades.Orm().Query().Where("id = ? OR email = ?", u.ID, u.Email).First(&existing); err != nil {
			if err := facades.Orm().Query().Create(&u); err != nil {
				return err
			}
		}
	}

	return nil
}

