package seeders

import (
	"encoding/json"
	"time"

	"elcomware/app/facades"
	"elcomware/app/models"
)

type ProductSeeder struct{}

// Signature The unique signature of the seeder.
func (s *ProductSeeder) Signature() string {
	return "ProductSeeder"
}

// Run executes the seeder logic.
func (s *ProductSeeder) Run() error {
	type seedProductDef struct {
		product     models.Product
		priceTiers  []models.PriceTier
		galleryUrls []string
		tags        []string
	}

	definitions := []seedProductDef{
		{
			product: models.Product{
				ID:                   "prod_b2b_cnc_01",
				Title:                "Industrial High-Precision 5-Axis CNC Milling Center V850",
				Slug:                 "industrial-high-precision-5-axis-cnc-milling-center-v850",
				Description:          "Heavy duty industrial CNC milling center engineered for titanium, aeronautical alloys, and precision mold casting. Features Siemens 828D CNC controller, 12,000 RPM high torque ceramic spindle, 24-station automatic tool changer (ATC), and full CE safety enclosure. Direct factory FOB/CIF export with on-site technician deployment options worldwide.",
				ShortDesc:            "Tier-1 5-Axis Heavy Machining Center for Aerospace & Precision Mold Tooling.",
				Category:             "Industrial & Machinery",
				SubCategory:          "Milling Machines",
				Type:                 models.TypePhysical,
				Format:               models.FormatWholesale,
				Condition:            models.ConditionNew,
				Price:                48500.00,
				WholesalePrice:       42000.00,
				MOQ:                  1,
				Stock:                15,
				Sku:                  "IND-CNC-5AX-850",
				VendorID:             "usr_vendor_01",
				VendorName:           "Apex Precision Heavy Machinery Co.",
				VendorCountry:        "Kenya",
				VendorRegion:         "East Africa",
				VendorVerified:       true,
				VendorBadge:          "Gold Supplier • Verified OEM",
				VendorPositiveRate:   "99.4%",
				TradeAssuranceCover:  1000000.00,
				AcceptsOffers:        true,
				Rating:               4.92,
				ReviewsCount:         128,
				ImageUrl:             "https://images.unsplash.com/photo-1581092160607-ee22621dd758?w=800&q=80",
				Featured:             true,
				LeadTimeDays:         21,
				Incoterm:             "FOB, CIF, DDP Available",
				LocalPickupAvailable: true,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
			priceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 2, Price: 48500.00},
				{MinQty: 3, MaxQty: 5, Price: 44000.00},
				{MinQty: 6, MaxQty: 99, Price: 39500.00},
			},
			galleryUrls: []string{
				"https://images.unsplash.com/photo-1581092160607-ee22621dd758?w=800&q=80",
				"https://images.unsplash.com/photo-1581092335397-9583fe92d232?w=800&q=80",
			},
			tags: []string{"CNC", "Machining", "Wholesale", "Industrial", "Heavy Equipment"},
		},
		{
			product: models.Product{
				ID:                   "prod_saas_erp_02",
				Title:                "Enterprise Multi-Tenant ERP & Global Supply Chain Platform",
				Slug:                 "enterprise-multi-tenant-erp-and-global-supply-chain-platform",
				Description:          "Full-featured cloud ERP suite with Docker microservice deployment, automated customs tariff compliance engine, multi-currency ledger, real-time inventory synchronization across multi-region distribution hubs, and native B2B EDI 850/855 connector. Includes perpetual enterprise commercial source code license + 1-year security updates.",
				ShortDesc:            "Turnkey Global ERP & Logistics Engine with Perpetual Commercial License.",
				Category:             "Digital Software",
				SubCategory:          "Enterprise Software",
				Type:                 models.TypeDigitalSoftware,
				Format:               models.FormatDigital,
				Condition:            models.ConditionNew,
				Price:                2499.00,
				WholesalePrice:       1899.00,
				MOQ:                  1,
				Stock:                -1,
				Sku:                  "SFT-ERP-ENT-2026",
				VendorID:             "usr_admin_01",
				VendorName:           "Soko Cloud Core Systems",
				VendorCountry:        "United States",
				VendorRegion:         "Global",
				VendorVerified:       true,
				VendorBadge:          "Certified Core Publisher",
				VendorPositiveRate:   "99.9%",
				TradeAssuranceCover:  500000.00,
				AcceptsOffers:        false,
				Rating:               4.98,
				ReviewsCount:         340,
				ImageUrl:             "https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800&q=80",
				Featured:             true,
				LeadTimeDays:         0,
				Incoterm:             "Instant Digital Delivery",
				LocalPickupAvailable: false,
				DigitalDownloadUrl:   "https://soko-storage.local/downloads/soko-enterprise-erp-v4.8.zip",
				DigitalLicenseFormat: "CRYPT-ECDSA-XXXX-YYYY-ZZZZ",
				DigitalFileSize:      "1.42 GB (Source + Containers)",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
			priceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 4, Price: 2499.00},
				{MinQty: 5, MaxQty: 19, Price: 1999.00},
				{MinQty: 20, MaxQty: 999, Price: 1499.00},
			},
			galleryUrls: []string{
				"https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800&q=80",
				"https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800&q=80",
			},
			tags: []string{"ERP", "Software", "Digital", "License", "Source Code", "Instant Delivery"},
		},
		{
			product: models.Product{
				ID:                   "prod_telecom_switch_03",
				Title:                "Enterprise Core 100G QSFP28 Data Center Switch (48-Port)",
				Slug:                 "enterprise-core-100g-qsfp28-data-center-switch-48-port",
				Description:          "High density 100GbE spine switch with 48x 100G QSFP28 and 8x 400G QSFP-DD uplink ports. Ultra-low 450ns cut-through latency, dual hot-swap redundant titanium PSUs, ONIE open network install environment supporting SONiC and Cumulus Linux.",
				ShortDesc:            "Carrier-grade 100GbE / 400GbE Low-Latency Data Center Switch.",
				Category:             "Electronics & Telecom",
				SubCategory:          "Networking Hardware",
				Type:                 models.TypePhysical,
				Format:               models.FormatRetail,
				Condition:            models.ConditionNew,
				Price:                7850.00,
				WholesalePrice:       6900.00,
				MOQ:                  1,
				Stock:                42,
				Sku:                  "NET-SW-100G-48P",
				VendorID:             "usr_vendor_02",
				VendorName:           "Shenzhen Nova Microelectronics",
				VendorCountry:        "China",
				VendorRegion:         "Asia-Pacific",
				VendorVerified:       true,
				VendorBadge:          "Gold Supplier • Tier 1",
				VendorPositiveRate:   "98.8%",
				TradeAssuranceCover:  250000.00,
				AcceptsOffers:        true,
				Rating:               4.85,
				ReviewsCount:         64,
				ImageUrl:             "https://images.unsplash.com/photo-1544197150-b99a580bb7a8?w=800&q=80",
				Featured:             true,
				LeadTimeDays:         5,
				Incoterm:             "Express Air Courier / DDP",
				LocalPickupAvailable: false,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
			priceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 4, Price: 7850.00},
				{MinQty: 5, MaxQty: 15, Price: 7200.00},
				{MinQty: 16, MaxQty: 100, Price: 6500.00},
			},
			galleryUrls: []string{
				"https://images.unsplash.com/photo-1544197150-b99a580bb7a8?w=800&q=80",
			},
			tags: []string{"Networking", "Datacenter", "100G", "Switch", "Hardware"},
		},
		{
			product: models.Product{
				ID:                   "prod_openbox_workstation_04",
				Title:                "Precision Dual Xeon Enterprise Rack Workstation 512GB ECC [Open Box]",
				Slug:                 "precision-dual-xeon-enterprise-rack-workstation-512gb-ecc-open-box",
				Description:          "Off-lease datacenter workstation in mint cosmetic and functional condition. Configured with dual Intel Xeon Platinum 8280 CPUs (56 cores total), 512GB DDR4 ECC Registered RAM, 4x 3.84TB Enterprise NVMe U.2 SSDs in hardware RAID 10, and dual NVIDIA RTX 6000 Ada 48GB GPUs. Seller accepts price negotiation offers.",
				ShortDesc:            "Open-Box High-Density Compute Workstation for AI Fine-Tuning & Simulation.",
				Category:             "Electronics & Telecom",
				SubCategory:          "Compute Servers",
				Type:                 models.TypePhysical,
				Format:               models.FormatMakeOffer,
				Condition:            models.ConditionOpenBox,
				Price:                11500.00,
				WholesalePrice:       10200.00,
				MOQ:                  1,
				Stock:                3,
				Sku:                  "SRV-XEON-512-OPNBX",
				VendorID:             "usr_vendor_01",
				VendorName:           "Apex Precision Heavy Machinery Co.",
				VendorCountry:        "Kenya",
				VendorRegion:         "East Africa",
				VendorVerified:       true,
				VendorBadge:          "Top Rated Seller",
				VendorPositiveRate:   "99.7%",
				TradeAssuranceCover:  150000.00,
				AcceptsOffers:        true,
				Rating:               4.95,
				ReviewsCount:         31,
				ImageUrl:             "https://images.unsplash.com/photo-1591488320449-011701bb6704?w=800&q=80",
				Featured:             false,
				LeadTimeDays:         3,
				Incoterm:             "DDP Insured Air Freight",
				LocalPickupAvailable: true,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
			priceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 3, Price: 11500.00},
			},
			galleryUrls: []string{
				"https://images.unsplash.com/photo-1591488320449-011701bb6704?w=800&q=80",
			},
			tags: []string{"Server", "Workstation", "Xeon", "eBay Format", "Make Offer", "Open Box"},
		},
		{
			product: models.Product{
				ID:                   "prod_agri_sorter_05",
				Title:                "Optical Multi-Spectrum Coffee & Grain Color Sorter 10T/H",
				Slug:                 "optical-multi-spectrum-coffee-grain-color-sorter-10th",
				Description:          "Industrial AI-powered optical sorter equipped with 4K RGB + InGaAs infrared cameras for automated sorting of specialty arabica coffee beans, macadamia nuts, sesame, and cereals. Eliminates defects, mold, and foreign objects with 99.98% purity. Installed across coffee cooperatives in Kenya, Ethiopia, and Uganda.",
				ShortDesc:            "10-Ton/Hour Industrial AI Optical Sorter for Commercial Grain Processing.",
				Category:             "Agriculture & Processing",
				SubCategory:          "Post-Harvest Machinery",
				Type:                 models.TypePhysical,
				Format:               models.FormatWholesale,
				Condition:            models.ConditionNew,
				Price:                32000.00,
				WholesalePrice:       27500.00,
				MOQ:                  1,
				Stock:                8,
				Sku:                  "AGR-OPT-SRT-10T",
				VendorID:             "usr_vendor_01",
				VendorName:           "Apex Precision Heavy Machinery Co.",
				VendorCountry:        "Kenya",
				VendorRegion:         "East Africa",
				VendorVerified:       true,
				VendorBadge:          "Gold Supplier • Verified OEM",
				VendorPositiveRate:   "99.4%",
				TradeAssuranceCover:  750000.00,
				AcceptsOffers:        true,
				Rating:               4.88,
				ReviewsCount:         49,
				ImageUrl:             "https://images.unsplash.com/photo-1500937386664-56d1dfef3854?w=800&q=80",
				Featured:             true,
				LeadTimeDays:         14,
				Incoterm:             "FOB Mombasa / CIF Worldwide",
				LocalPickupAvailable: true,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
			priceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 2, Price: 32000.00},
				{MinQty: 3, MaxQty: 5, Price: 29000.00},
				{MinQty: 6, MaxQty: 20, Price: 26000.00},
			},
			galleryUrls: []string{
				"https://images.unsplash.com/photo-1500937386664-56d1dfef3854?w=800&q=80",
			},
			tags: []string{"Agri-Tech", "Coffee Sorter", "Processing", "Wholesale", "East Africa"},
		},
		{
			product: models.Product{
				ID:                   "prod_ai_model_weights_06",
				Title:                "Precision Vision AI weights & Commercial Inspection Engine",
				Slug:                 "precision-vision-ai-weights-commercial-inspection-engine",
				Description:          "Trained on over 25 million industrial quality control images. Zero-shot defect detection SDK with C++, Python, and Go bindings. Includes pre-quantized TensorRT models for Jetson Orin edge devices and full PyTorch checkpoint weights for enterprise fine-tuning.",
				ShortDesc:            "Industrial Defect Computer Vision Model Weights with Commercial Royalty-Free License.",
				Category:             "Digital Assets",
				SubCategory:          "AI & Machine Learning",
				Type:                 models.TypeDigitalAsset,
				Format:               models.FormatDigital,
				Condition:            models.ConditionNew,
				Price:                1200.00,
				WholesalePrice:       850.00,
				MOQ:                  1,
				Stock:                -1,
				Sku:                  "AI-VIS-DEFECT-V3",
				VendorID:             "usr_admin_01",
				VendorName:           "Soko Cloud Core Systems",
				VendorCountry:        "United States",
				VendorRegion:         "Global",
				VendorVerified:       true,
				VendorBadge:          "Certified Core Publisher",
				VendorPositiveRate:   "99.9%",
				TradeAssuranceCover:  300000.00,
				AcceptsOffers:        false,
				Rating:               4.96,
				ReviewsCount:         115,
				ImageUrl:             "https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?w=800&q=80",
				Featured:             false,
				LeadTimeDays:         0,
				Incoterm:             "Instant Digital Delivery",
				LocalPickupAvailable: false,
				DigitalDownloadUrl:   "https://soko-storage.local/downloads/soko-vision-ai-weights-v3.tar.gz",
				DigitalLicenseFormat: "AI-WEIGHTS-XXXX-YYYY",
				DigitalFileSize:      "8.65 GB",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
			priceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 99, Price: 1200.00},
			},
			galleryUrls: []string{
				"https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?w=800&q=80",
			},
			tags: []string{"AI", "Computer Vision", "Digital Asset", "Instant Delivery", "Weights"},
		},
	}

	for _, def := range definitions {
		tiersBytes, _ := json.Marshal(def.priceTiers)
		galleryBytes, _ := json.Marshal(def.galleryUrls)
		tagsBytes, _ := json.Marshal(def.tags)

		p := def.product
		p.PriceTiersJSON = string(tiersBytes)
		p.GalleryUrlsJSON = string(galleryBytes)
		p.TagsJSON = string(tagsBytes)

		var existing models.Product
		if err := facades.Orm().Query().Where("id = ?", p.ID).First(&existing); err != nil {
			if err := facades.Orm().Query().Create(&p); err != nil {
				return err
			}
		}
	}

	return nil
}
