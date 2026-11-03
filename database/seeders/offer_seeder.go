package seeders

import (
	"time"

	"elcomware/app/facades"
	"elcomware/app/models"
)

type OfferSeeder struct{}

// Signature The unique signature of the seeder.
func (s *OfferSeeder) Signature() string {
	return "OfferSeeder"
}

// Run executes the seeder logic.
func (s *OfferSeeder) Run() error {
	offers := []models.Offer{
		{
			ID:           "off_001_workstation",
			ProductID:    "prod_openbox_workstation_04",
			ProductTitle: "Precision Dual Xeon Enterprise Rack Workstation 512GB ECC [Open Box]",
			BuyerID:      "usr_buyer_01",
			BuyerName:    "David Mwangi",
			BuyerEmail:   "buyer@soko.global",
			BuyerCompany: "Great Rift Infrastructure Consortium",
			OfferedPrice: 9800.00,
			Quantity:     2,
			TotalOffered: 19600.00,
			Message:      "We can execute immediate wire payment and provide commercial freight collection if price is acceptable.",
			Status:       models.OfferStatusPending,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           "off_002_cnc",
			ProductID:    "prod_b2b_cnc_01",
			ProductTitle: "Industrial High-Precision 5-Axis CNC Milling Center V850",
			BuyerID:      "usr_buyer_01",
			BuyerName:    "David Mwangi",
			BuyerEmail:   "buyer@soko.global",
			BuyerCompany: "Great Rift Infrastructure Consortium",
			OfferedPrice: 41000.00,
			Quantity:     1,
			TotalOffered: 41000.00,
			Message:      "Requesting quote with delivery to Nairobi Industrial Area without overseas installation.",
			Status:       models.OfferStatusAccepted,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, o := range offers {
		exists, _ := facades.Orm().Query().Model(&models.Offer{}).Where("id = ?", o.ID).Exists()
		if !exists {
			if err := facades.Orm().Query().Create(&o); err != nil {
				return err
			}
		}
	}

	return nil
}

