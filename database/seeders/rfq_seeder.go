package seeders

import (
	"time"

	"elcomware/app/facades"
	"elcomware/app/models"
)

type RfqSeeder struct{}

// Signature The unique signature of the seeder.
func (s *RfqSeeder) Signature() string {
	return "RfqSeeder"
}

// Run executes the seeder logic.
func (s *RfqSeeder) Run() error {
	rfqs := []struct {
		rfq    models.Rfq
		quotes []models.RfqQuote
	}{
		{
			rfq: models.Rfq{
				ID:                 "rfq_001_solar_bess",
				RfqNumber:          "RFQ-2026-NA-0912",
				BuyerID:            "usr_buyer_01",
				BuyerName:          "David Mwangi",
				BuyerCompany:       "Great Rift Infrastructure Consortium",
				BuyerCountry:       "Kenya",
				Title:              "500kW Commercial Solar BESS Energy Storage Battery Inverter Units",
				Category:           "Energy & Utilities",
				Quantity:           10,
				TargetPrice:        22000.00,
				DestinationCountry: "Kenya (Port of Mombasa DDP)",
				Specifications:     "LiFePO4 containerized battery racks with IP65 outdoor enclosure, PCS bidirectional inverter integration, SCADA MODBUS TCP monitoring, UL9540A fire suppression certificate, and minimum 6000 cycle warranty at 80% DoD.",
				Status:             models.RfqStatusOpen,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
			quotes: []models.RfqQuote{
				{
					ID:            "quote_001_apex",
					RfqID:         "rfq_001_solar_bess",
					VendorID:      "usr_vendor_01",
					VendorName:    "Apex Precision Heavy Machinery Co.",
					VendorCountry: "Kenya",
					UnitPrice:     21500.00,
					TotalAmount:   215000.00,
					LeadTimeDays:  18,
					Notes:         "Includes local assembly in Nairobi, EAC duty-free certificate, and 5-year on-site SLA.",
					Status:        "pending",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				},
				{
					ID:            "quote_002_shenzhen",
					RfqID:         "rfq_001_solar_bess",
					VendorID:      "usr_vendor_02",
					VendorName:    "Shenzhen Nova Microelectronics",
					VendorCountry: "China",
					UnitPrice:     19800.00,
					TotalAmount:   198000.00,
					LeadTimeDays:  30,
					Notes:         "FOB Shenzhen seaport. Tier-1 CATL cells with full UL9540A test report. Sea freight options available.",
					Status:        "pending",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				},
			},
		},
		{
			rfq: models.Rfq{
				ID:                 "rfq_002_telecom_antennas",
				RfqNumber:          "RFQ-2026-EU-4481",
				BuyerID:            "usr_buyer_01",
				BuyerName:          "David Mwangi",
				BuyerCompany:       "Great Rift Infrastructure Consortium",
				BuyerCountry:       "Kenya",
				Title:              "5G Massive MIMO Active Antenna Units (Sub-6GHz Band n78)",
				Category:           "Electronics & Telecom",
				Quantity:           50,
				TargetPrice:        3500.00,
				DestinationCountry: "Germany (Frankfurt Hub)",
				Specifications:     "64T64R beamforming array, 200W total output power, CPRI/eCPRI optical interfaces, operating temperature -40C to +55C.",
				Status:             models.RfqStatusOpen,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
			quotes: []models.RfqQuote{
				{
					ID:            "quote_003_shenzhen",
					RfqID:         "rfq_002_telecom_antennas",
					VendorID:      "usr_vendor_02",
					VendorName:    "Shenzhen Nova Microelectronics",
					VendorCountry: "China",
					UnitPrice:     3400.00,
					TotalAmount:   170000.00,
					LeadTimeDays:  25,
					Notes:         "Full CE/RED certification, 3-year replacement warranty, samples shipped in 3 days.",
					Status:        "pending",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				},
			},
		},
	}

	for _, item := range rfqs {
		var existing models.Rfq
		if err := facades.Orm().Query().Where("id = ?", item.rfq.ID).First(&existing); err != nil {
			if err := facades.Orm().Query().Create(&item.rfq); err != nil {
				return err
			}
			for _, q := range item.quotes {
				if err := facades.Orm().Query().Create(&q); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
