package services

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"elcomware/app/models"
)

type Store struct {
	mu         sync.RWMutex
	users      map[string]*models.User
	products   map[string]*models.Product
	orders     map[string]*models.Order
	rfqs       map[string]*models.Rfq
	offers     map[string]*models.Offer
	currencies map[string]models.Currency
	regions    map[string]models.MarketRegion
}

var (
	globalStore *Store
	once        sync.Once
)

func GetStore() *Store {
	once.Do(func() {
		globalStore = &Store{
			users:      make(map[string]*models.User),
			products:   make(map[string]*models.Product),
			orders:     make(map[string]*models.Order),
			rfqs:       make(map[string]*models.Rfq),
			offers:     make(map[string]*models.Offer),
			currencies: make(map[string]models.Currency),
			regions:    make(map[string]models.MarketRegion),
		}
		globalStore.seed()
	})
	return globalStore
}

func (s *Store) seed() {
	// 1. Regional Currencies
	currencies := []models.Currency{
		{Code: "USD", Symbol: "$", Name: "US Dollar", ExchangeRate: 1.0, Prefix: true},
		{Code: "EUR", Symbol: "€", Name: "Euro", ExchangeRate: 0.92, Prefix: true},
		{Code: "GBP", Symbol: "£", Name: "British Pound", ExchangeRate: 0.78, Prefix: true},
		{Code: "KES", Symbol: "KSh", Name: "Kenya Shilling", ExchangeRate: 130.0, Prefix: false},
		{Code: "NGN", Symbol: "₦", Name: "Nigerian Naira", ExchangeRate: 1550.0, Prefix: true},
		{Code: "ZAR", Symbol: "R", Name: "South African Rand", ExchangeRate: 18.5, Prefix: true},
		{Code: "AED", Symbol: "AED", Name: "UAE Dirham", ExchangeRate: 3.67, Prefix: false},
		{Code: "CNY", Symbol: "¥", Name: "Chinese Yuan", ExchangeRate: 7.23, Prefix: true},
		{Code: "JPY", Symbol: "¥", Name: "Japanese Yen", ExchangeRate: 152.0, Prefix: true},
		{Code: "INR", Symbol: "₹", Name: "Indian Rupee", ExchangeRate: 84.0, Prefix: true},
	}
	for _, c := range currencies {
		s.currencies[c.Code] = c
	}

	// 2. Regional Hubs (Tailorable to any market)
	regions := []models.MarketRegion{
		{
			ID:                 "global",
			Name:               "Global International",
			Flag:               "🌐",
			DefaultCurrency:    "USD",
			LocalLogistics:     "DHL Express, Maersk Multimodal & Air Cargo",
			PaymentRails:       []string{"Trade Escrow", "SWIFT Wire", "Credit Card", "USDC"},
			VatPercentage:      0.0,
			CustomsClearanceHr: 48,
		},
		{
			ID:                 "north_america",
			Name:               "North America (US & CA)",
			Flag:               "🇺🇸",
			DefaultCurrency:    "USD",
			LocalLogistics:     "FedEx Freight, UPS Domestic Next-Day",
			PaymentRails:       []string{"ACH", "Stripe Cards", "Trade Escrow"},
			VatPercentage:      7.5,
			CustomsClearanceHr: 24,
		},
		{
			ID:                 "europe",
			Name:               "Europe & UK",
			Flag:               "🇪🇺",
			DefaultCurrency:    "EUR",
			LocalLogistics:     "DB Schenker, DPD Inter-Europe, Royal Mail",
			PaymentRails:       []string{"SEPA Direct Debit", "Trade Escrow", "iDEAL"},
			VatPercentage:      20.0,
			CustomsClearanceHr: 12,
		},
		{
			ID:                 "east_africa",
			Name:               "East Africa Regional Hub",
			Flag:               "🇰🇪",
			DefaultCurrency:    "KES",
			LocalLogistics:     "SGR Freight, East Africa Cross-Border Clearing, Sendy",
			PaymentRails:       []string{"M-Pesa Express", "Pesalink", "Trade Escrow"},
			VatPercentage:      16.0,
			CustomsClearanceHr: 24,
		},
		{
			ID:                 "west_africa",
			Name:               "West Africa Hub",
			Flag:               "🇳🇬",
			DefaultCurrency:    "NGN",
			LocalLogistics:     "Lekki Deep Sea Logistics, GIG Logistics",
			PaymentRails:       []string{"Paystack / Flutterwave", "Bank Wire", "Escrow"},
			VatPercentage:      7.5,
			CustomsClearanceHr: 48,
		},
		{
			ID:                 "middle_east",
			Name:               "Middle East & GCC",
			Flag:               "🇦🇪",
			DefaultCurrency:    "AED",
			LocalLogistics:     "Aramex Middle East, Jebel Ali Port Transshipment",
			PaymentRails:       []string{"Trade Escrow", "Cards", "Direct Bank"},
			VatPercentage:      5.0,
			CustomsClearanceHr: 18,
		},
		{
			ID:                 "asia_pacific",
			Name:               "Asia-Pacific Manufacturing Hub",
			Flag:               "🇨🇳",
			DefaultCurrency:    "CNY",
			LocalLogistics:     "SF Express, COSCO Shipping Lines, Ningbo Cargo",
			PaymentRails:       []string{"Alipay Global", "WeChat Pay", "Telegraphic Transfer (T/T)"},
			VatPercentage:      13.0,
			CustomsClearanceHr: 24,
		},
	}
	for _, r := range regions {
		s.regions[r.ID] = r
	}

	// 3. Seed Users (Enterprise Admin, Global Suppliers, Regional Buyers)
	admin := &models.User{
		ID:          "usr-admin-1",
		Name:        "Soko Global Admin",
		Email:       "admin@soko.global",
		Password:    "admin123",
		Role:        models.RoleAdmin,
		CompanyName: "Soko Worldwide Commerce Corp",
		Country:     "United States",
		AccountType: "enterprise",
		Verified:    true,
		Avatar:      "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80",
		CreatedAt:   time.Now().Add(-180 * 24 * time.Hour),
	}
	s.users[admin.Email] = admin

	merchant := &models.User{
		ID:          "usr-merchant-1",
		Name:        "Apex Industrial & Cloud Tech Ltd",
		Email:       "merchant@apextech.com",
		Password:    "merchant123",
		Role:        models.RoleMerchant,
		CompanyName: "Apex Technologies Group",
		Country:     "Germany",
		AccountType: "supplier",
		Verified:    true,
		Avatar:      "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&auto=format&fit=crop&q=80",
		CreatedAt:   time.Now().Add(-90 * 24 * time.Hour),
	}
	s.users[merchant.Email] = merchant

	merchantAfrica := &models.User{
		ID:          "usr-merchant-2",
		Name:        "Kilima Agri-Tech & Machinery Ltd",
		Email:       "vendor@kilima.co.ke",
		Password:    "vendor123",
		Role:        models.RoleMerchant,
		CompanyName: "Kilima Engineering Works",
		Country:     "Kenya",
		AccountType: "supplier",
		Verified:    true,
		Avatar:      "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150&auto=format&fit=crop&q=80",
		CreatedAt:   time.Now().Add(-60 * 24 * time.Hour),
	}
	s.users[merchantAfrica.Email] = merchantAfrica

	buyer := &models.User{
		ID:          "usr-buyer-1",
		Name:        "Elena Rostova",
		Email:       "buyer@megacorp.com",
		Password:    "buyer123",
		Role:        models.RoleBuyer,
		CompanyName: "MegaCorp Procurement Int.",
		Country:     "United Kingdom",
		AccountType: "enterprise",
		Verified:    true,
		Avatar:      "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150&auto=format&fit=crop&q=80",
		CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
	}
	s.users[buyer.Email] = buyer

	// 4. Seed Diverse Products (Alibaba Wholesale + Amazon Retail + eBay Make Offer + Instant Digital)
	products := []*models.Product{
		// [Alibaba Style]: Heavy Machinery with Volume Tier Pricing
		{
			ID:          "prod-ali-1",
			Title:       "Industrial High-Speed CNC Fiber Laser Cutting Machine 3000W",
			Slug:        "industrial-cnc-fiber-laser-3000w",
			Description: "Heavy-duty aerospace-grade welded steel bed with Raycus laser source, CypCut CNC operating system, auto-exchange dual shuttle tables, and CE certification for precision sheet metal fabrication.",
			ShortDesc:   "3000W industrial fiber laser for 22mm carbon steel and 12mm stainless steel.",
			Category:    "Industrial & Machinery",
			SubCategory: "Fabrication & CNC",
			Type:        models.TypePhysical,
			Format:      models.FormatWholesale,
			Condition:   models.ConditionNew,
			Price:       18500.00,
			WholesalePrice: 15200.00,
			PriceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 2, Price: 18500.00},
				{MinQty: 3, MaxQty: 9, Price: 16800.00},
				{MinQty: 10, MaxQty: 999, Price: 15200.00},
			},
			MOQ:                 1,
			Stock:               18,
			Sku:                 "SKU-IND-LSR3000",
			VendorID:            merchant.ID,
			VendorName:          "Bavaria Precision Machine Works",
			VendorCountry:       "Germany",
			VendorRegion:        "europe",
			VendorVerified:      true,
			VendorBadge:         "Gold Verified Supplier",
			VendorPositiveRate:  "99.4%",
			TradeAssuranceCover: 2500000.00,
			AcceptsOffers:       true,
			Rating:              4.9,
			ReviewsCount:        84,
			ImageUrl:            "https://images.unsplash.com/photo-1581092160607-ee22621dd758?w=800&auto=format&fit=crop&q=80",
			Tags:                []string{"Industrial", "CNC", "Laser", "Machinery", "FOB", "CIF"},
			Featured:            true,
			LeadTimeDays:        14,
			Incoterm:            "FOB / CIF Hamburg",
			CreatedAt:           time.Now().Add(-60 * 24 * time.Hour),
		},

		// [Amazon + Digital Style]: Instant ERP Cloud Software License
		{
			ID:          "prod-dig-1",
			Title:       "Soko Cloud ERP Enterprise Suite v5 (Unlimited Multi-Tenant)",
			Slug:        "soko-cloud-erp-enterprise",
			Description: "Full-stack multi-warehouse supply chain and enterprise resource planning software with real-time financial tracking, multi-currency ledger, and REST/GraphQL APIs.",
			ShortDesc:   "Tier-1 enterprise ERP with automated billing, logistics, and multi-tenant ledger.",
			Category:    "Digital Software",
			SubCategory: "Enterprise SaaS",
			Type:        models.TypeDigitalSoftware,
			Format:      models.FormatDigital,
			Condition:   models.ConditionNew,
			Price:       2499.00,
			WholesalePrice: 1999.00,
			PriceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 4, Price: 2499.00},
				{MinQty: 5, MaxQty: 19, Price: 2199.00},
				{MinQty: 20, MaxQty: 999, Price: 1999.00},
			},
			MOQ:                  1,
			Stock:                -1, // Unlimited digital
			Sku:                  "SKU-SOK-ERP5",
			VendorID:             merchant.ID,
			VendorName:           merchant.Name,
			VendorCountry:        merchant.Country,
			VendorRegion:         "global",
			VendorVerified:       true,
			VendorBadge:          "Top Rated Software Partner",
			VendorPositiveRate:   "99.8%",
			TradeAssuranceCover:  1000000.00,
			AcceptsOffers:        false,
			Rating:               4.9,
			ReviewsCount:         136,
			ImageUrl:             "https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800&auto=format&fit=crop&q=80",
			Tags:                 []string{"SaaS", "ERP", "Cloud", "Instant License"},
			Featured:             true,
			LeadTimeDays:         0,
			Incoterm:             "Instant Cloud Delivery",
			DigitalDownloadUrl:   "https://cdn.soko.global/releases/soko-erp-v5.3.zip",
			DigitalLicenseFormat: "ERP5-XXXX-XXXX-ENTERPRISE",
			DigitalFileSize:      "1.42 GB",
			CreatedAt:            time.Now().Add(-45 * 24 * time.Hour),
		},

		// [eBay Style]: Certified Refurbished Enterprise Switch / Datacenter Router
		{
			ID:          "prod-ebay-1",
			Title:       "Cisco Catalyst 9300 48-Port PoE+ Switch (Factory Certified Refurbished)",
			Slug:        "cisco-catalyst-9300-48port-refurb",
			Description: "Inspected, cleaned and re-certified with 1-Year Soko Trade Warranty. 48 ports PoE+ with modular uplink support and dual power supply redundancy. Open to Make an Offer for bulk procurement.",
			ShortDesc:   "Certified refurbished 48-port enterprise network core switch with 1-year warranty.",
			Category:    "Electronics & Telecom",
			SubCategory: "Datacenter Core",
			Type:        models.TypePhysical,
			Format:      models.FormatMakeOffer,
			Condition:   models.ConditionRefurbished,
			Price:       1150.00,
			WholesalePrice: 920.00,
			PriceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 4, Price: 1150.00},
				{MinQty: 5, MaxQty: 19, Price: 1020.00},
				{MinQty: 20, MaxQty: 100, Price: 920.00},
			},
			MOQ:                 1,
			Stock:               45,
			Sku:                 "SKU-CISCO-9300-REF",
			VendorID:            merchant.ID,
			VendorName:          "Global Datacenter Surplus Corp",
			VendorCountry:       "United States",
			VendorRegion:        "north_america",
			VendorVerified:      true,
			VendorBadge:         "Top Rated Plus Seller",
			VendorPositiveRate:  "99.1%",
			TradeAssuranceCover: 500000.00,
			AcceptsOffers:       true,
			Rating:              4.8,
			ReviewsCount:        92,
			ImageUrl:            "https://images.unsplash.com/photo-1544197150-b99a580bb7a8?w=800&auto=format&fit=crop&q=80",
			Tags:                []string{"Cisco", "Refurbished", "Switch", "Make Offer", "Warranty"},
			Featured:            true,
			LeadTimeDays:        2,
			Incoterm:            "DDP Express 2-Day Air",
			LocalPickupAvailable: true,
			CreatedAt:           time.Now().Add(-35 * 24 * time.Hour),
		},

		// [Regional Local Hub - East Africa]: Optical Grain & Coffee Color Sorter Machine
		{
			ID:          "prod-reg-1",
			Title:       "Kilima OptiSort Multichannel Coffee & Grain Color Sorter 1.5 T/h",
			Slug:        "kilima-optisort-coffee-grain-sorter",
			Description: "Engineered in Nairobi for East African coffee cooperatives and grain millers. High-definition 5400-pixel CCD cameras with AI color and defect rejection. Eligible for local M-Pesa escrow and direct Mombasa port clearing.",
			ShortDesc:   "High-precision 1.5 ton/hour optical sorting machine for specialty Arabica coffee and grains.",
			Category:    "Agriculture & Processing",
			SubCategory: "Post-Harvest Milling",
			Type:        models.TypePhysical,
			Format:      models.FormatWholesale,
			Condition:   models.ConditionNew,
			Price:       8600.00,
			WholesalePrice: 7200.00,
			PriceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 1, Price: 8600.00},
				{MinQty: 2, MaxQty: 5, Price: 7800.00},
				{MinQty: 6, MaxQty: 20, Price: 7200.00},
			},
			MOQ:                 1,
			Stock:               8,
			Sku:                 "SKU-KLM-OPT15",
			VendorID:            merchantAfrica.ID,
			VendorName:          merchantAfrica.Name,
			VendorCountry:       merchantAfrica.Country,
			VendorRegion:        "east_africa",
			VendorVerified:      true,
			VendorBadge:         "Verified Regional Manufacturer",
			VendorPositiveRate:  "100%",
			TradeAssuranceCover: 300000.00,
			AcceptsOffers:       true,
			Rating:              5.0,
			ReviewsCount:        34,
			ImageUrl:            "https://images.unsplash.com/photo-1500651230702-0e2d8a49d4ad?w=800&auto=format&fit=crop&q=80",
			Tags:                []string{"AgriTech", "Coffee", "Sorter", "East Africa", "M-Pesa"},
			Featured:            true,
			LeadTimeDays:        4,
			Incoterm:            "Ex-Works Nairobi / Mombasa DDP",
			LocalPickupAvailable: true,
			CreatedAt:           time.Now().Add(-20 * 24 * time.Hour),
		},

		// [Digital AI Model]: Freight Route Neural Network Weights
		{
			ID:          "prod-dig-2",
			Title:       "Global Multimodal Freight Route & Fuel Optimization Neural Model",
			Slug:        "freight-route-opt-ai-model",
			Description: "Pretrained ONNX and PyTorch weights fine-tuned on 100M+ maritime and air-cargo shipment routes to minimize fuel consumption and customs clearance latency.",
			ShortDesc:   "Industrial AI neural network for container vessel routing and air freight scheduling.",
			Category:    "Digital Assets",
			SubCategory: "AI & Machine Learning",
			Type:        models.TypeDigitalAsset,
			Format:      models.FormatDigital,
			Condition:   models.ConditionNew,
			Price:       680.00,
			WholesalePrice: 540.00,
			PriceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 4, Price: 680.00},
				{MinQty: 5, MaxQty: 99, Price: 540.00},
			},
			MOQ:                  1,
			Stock:                -1,
			Sku:                  "SKU-AI-ROUTEX",
			VendorID:             merchant.ID,
			VendorName:           merchant.Name,
			VendorCountry:        merchant.Country,
			VendorRegion:         "global",
			VendorVerified:       true,
			VendorBadge:          "AI Research Lab Certified",
			VendorPositiveRate:   "99.5%",
			TradeAssuranceCover:  500000.00,
			AcceptsOffers:        true,
			Rating:               4.8,
			ReviewsCount:         64,
			ImageUrl:             "https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?w=800&auto=format&fit=crop&q=80",
			Tags:                 []string{"AI", "PyTorch", "Logistics", "Digital Download"},
			Featured:             false,
			LeadTimeDays:         0,
			Incoterm:             "Instant Download",
			DigitalDownloadUrl:   "https://cdn.soko.global/weights/freight-opt-v2.pt.tar.gz",
			DigitalLicenseFormat: "AI-MODEL-XXXX-XXXX-XXXX",
			DigitalFileSize:      "890 MB",
			CreatedAt:            time.Now().Add(-30 * 24 * time.Hour),
		},

		// [Alibaba Wholesale]: Bare Copper Cat6A Bulk Spools (1000ft)
		{
			ID:          "prod-ali-2",
			Title:       "Enterprise Shielded Cat6A 10Gbps Bare Copper Cable (1000ft Bulk Spool)",
			Slug:        "enterprise-cat6a-shielded-1000ft",
			Description: "Solid 23AWG pure oxygen-free copper, CMR riser rated, tested up to 650MHz bandwidth with foil shielding (F/UTP) for zero cross-talk in datacenter and commercial installations.",
			ShortDesc:   "Datacenter-certified 10Gbps Ethernet cable spools for commercial enterprise cabling.",
			Category:    "Electronics & Telecom",
			SubCategory: "Networking Infrastructure",
			Type:        models.TypePhysical,
			Format:      models.FormatWholesale,
			Condition:   models.ConditionNew,
			Price:       145.00,
			WholesalePrice: 98.00,
			PriceTiers: []models.PriceTier{
				{MinQty: 10, MaxQty: 49, Price: 125.00},
				{MinQty: 50, MaxQty: 199, Price: 108.00},
				{MinQty: 200, MaxQty: 9999, Price: 98.00},
			},
			MOQ:                 10,
			Stock:               650,
			Sku:                 "SKU-NET-C6A-1000",
			VendorID:            merchant.ID,
			VendorName:          "Shenzhen Global Opto-Cable Ltd",
			VendorCountry:       "China",
			VendorRegion:        "asia_pacific",
			VendorVerified:      true,
			VendorBadge:         "Gold Supplier 8 Years",
			VendorPositiveRate:  "98.9%",
			TradeAssuranceCover: 1200000.00,
			AcceptsOffers:       true,
			Rating:              4.8,
			ReviewsCount:        312,
			ImageUrl:            "https://images.unsplash.com/photo-1544197150-b99a580bb7a8?w=800&auto=format&fit=crop&q=80",
			Tags:                []string{"Datacenter", "Ethernet", "Wholesale", "FOB Shenzhen"},
			Featured:            true,
			LeadTimeDays:        5,
			Incoterm:            "FOB Shenzhen / CIF Global",
			CreatedAt:           time.Now().Add(-40 * 24 * time.Hour),
		},

		// [Amazon Retail Style]: Commercial LiFePO4 Energy Storage Cabinet
		{
			ID:          "prod-ali-3",
			Title:       "Commercial LiFePO4 Energy Storage Cabinet 100kWh 400V (Liquid Cooled)",
			Slug:        "commercial-lifepo4-energy-storage-100kwh",
			Description: "Scalable containerized lithium iron phosphate battery system with built-in aerosol fire suppression, active liquid cooling, 6000+ cycle life, and CAN/RS485 smart BMS integration.",
			ShortDesc:   "100kWh industrial power backup and peak-shaving energy storage unit.",
			Category:    "Energy & Utilities",
			SubCategory: "Commercial Storage",
			Type:        models.TypePhysical,
			Format:      models.FormatWholesale,
			Condition:   models.ConditionNew,
			Price:       22000.00,
			WholesalePrice: 18500.00,
			PriceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 2, Price: 22000.00},
				{MinQty: 3, MaxQty: 9, Price: 19800.00},
				{MinQty: 10, MaxQty: 50, Price: 18500.00},
			},
			MOQ:                 1,
			Stock:               12,
			Sku:                 "SKU-NRG-BATT100",
			VendorID:            merchant.ID,
			VendorName:          "Nordic Clean Energy Alliance",
			VendorCountry:       "Sweden",
			VendorRegion:        "europe",
			VendorVerified:      true,
			VendorBadge:         "Tier-1 BESS Manufacturer",
			VendorPositiveRate:  "100%",
			TradeAssuranceCover: 3000000.00,
			AcceptsOffers:       true,
			Rating:              5.0,
			ReviewsCount:        24,
			ImageUrl:            "https://images.unsplash.com/photo-1513836279014-a89f7a76ae86?w=800&auto=format&fit=crop&q=80",
			Tags:                []string{"Energy", "Storage", "LiFePO4", "DDP Europe"},
			Featured:            true,
			LeadTimeDays:        21,
			Incoterm:            "DDP / CIF Major Ports",
			CreatedAt:           time.Now().Add(-10 * 24 * time.Hour),
		},

		// [eBay Style - Make an Offer]: Autonomous Mobile Warehouse Robot AMR 500kg
		{
			ID:          "prod-ebay-2",
			Title:       "Autonomous Mobile Warehouse Robot (AMR) 500kg Payload with Docking Hub",
			Slug:        "amr-warehouse-robot-500kg",
			Description: "LiDAR SLAM navigation without magnetic floor tapes, 360-degree obstacle detection, smart swarm fleet coordination algorithm, and rapid 20-minute wireless contact charging.",
			ShortDesc:   "Next-gen automated guided vehicle for smart e-commerce logistics centers.",
			Category:    "Industrial & Machinery",
			SubCategory: "Warehouse Robotics",
			Type:        models.TypePhysical,
			Format:      models.FormatMakeOffer,
			Condition:   models.ConditionNew,
			Price:       9400.00,
			WholesalePrice: 7800.00,
			PriceTiers: []models.PriceTier{
				{MinQty: 1, MaxQty: 1, Price: 9400.00},
				{MinQty: 2, MaxQty: 5, Price: 8500.00},
				{MinQty: 6, MaxQty: 30, Price: 7800.00},
			},
			MOQ:                 1,
			Stock:               28,
			Sku:                 "SKU-ROB-AMR500",
			VendorID:            merchant.ID,
			VendorName:          "Tokyo Automation Dynamics",
			VendorCountry:       "Japan",
			VendorRegion:        "asia_pacific",
			VendorVerified:      true,
			VendorBadge:         "Certified Robotics Manufacturer",
			VendorPositiveRate:  "99.7%",
			TradeAssuranceCover: 1800000.00,
			AcceptsOffers:       true,
			Rating:              4.9,
			ReviewsCount:        52,
			ImageUrl:            "https://images.unsplash.com/photo-1485827404703-89b55fcc595e?w=800&auto=format&fit=crop&q=80",
			Tags:                []string{"Robotics", "AMR", "Automation", "Make Offer"},
			Featured:            false,
			LeadTimeDays:        10,
			Incoterm:            "DDP Air Cargo / Sea",
			CreatedAt:           time.Now().Add(-25 * 24 * time.Hour),
		},
	}

	for _, p := range products {
		s.products[p.ID] = p
	}

	// 5. Seed B2B RFQs (Request For Quotation - Alibaba Core Feature)
	rfqs := []*models.Rfq{
		{
			ID:                 "rfq-101",
			RfqNumber:          "RFQ-2026-101",
			BuyerID:            buyer.ID,
			BuyerName:          buyer.Name,
			BuyerCompany:       buyer.CompanyName,
			BuyerCountry:       buyer.Country,
			Title:              "Procurement of 500 Industrial Cat6A Shielded Ethernet Cable Spools",
			Category:           "Electronics & Telecom",
			Quantity:           500,
			TargetPrice:        90.00,
			DestinationCountry: "United Kingdom (Port of Felixstowe)",
			Specifications:     "23AWG Solid Bare Copper, CMR Rated, ETL Verified, Fluke Channel Tested.",
			Status:             models.RfqStatusOpen,
			Quotes: []models.RfqQuote{
				{
					ID:           "quot-1",
					RfqID:        "rfq-101",
					VendorID:     merchant.ID,
					VendorName:   "Shenzhen Global Opto-Cable Ltd",
					VendorCountry: "China",
					UnitPrice:    88.50,
					TotalAmount:  44250.00,
					LeadTimeDays: 7,
					Notes:        "CIF Felixstowe included with Lloyd's marine insurance.",
					Status:       "pending",
					CreatedAt:    time.Now().Add(-24 * time.Hour),
				},
			},
			CreatedAt: time.Now().Add(-48 * time.Hour),
		},
		{
			ID:                 "rfq-102",
			RfqNumber:          "RFQ-2026-102",
			BuyerID:            buyer.ID,
			BuyerName:          "East Africa Milling Corp",
			BuyerCompany:       "East Africa Grain Millers Ltd",
			BuyerCountry:       "Kenya",
			Title:              "3x Commercial 5-Ton Grain Dryers and Elevators",
			Category:           "Agriculture & Processing",
			Quantity:           3,
			TargetPrice:        14000.00,
			DestinationCountry: "Kenya (Nairobi Industrial Area)",
			Specifications:     "Diesel or biomass powered, automated humidity sensor, CE compliant.",
			Status:             models.RfqStatusOpen,
			Quotes:             []models.RfqQuote{},
			CreatedAt:          time.Now().Add(-12 * time.Hour),
		},
	}
	for _, r := range rfqs {
		s.rfqs[r.ID] = r
	}

	// 6. Seed Sample Orders
	order1 := &models.Order{
		ID:            "ord-9821-us",
		OrderNumber:   "SOKO-2026-9821",
		UserID:        buyer.ID,
		CustomerName:  buyer.Name,
		CustomerEmail: buyer.Email,
		TotalAmount:   2499.00,
		Currency:      "USD",
		Status:        models.OrderStatusCompleted,
		PaymentMethod: "Corporate Wire / Escrow",
		HasDigital:    true,
		HasPhysical:   false,
		Items: []models.OrderItem{
			{
				ProductID:          "prod-dig-1",
				ProductTitle:       "Soko Cloud ERP Enterprise Suite v5",
				ProductType:        models.TypeDigitalSoftware,
				Quantity:           1,
				UnitPrice:          2499.00,
				TotalPrice:         2499.00,
				DigitalLicenseKey:  "ERP5-9821-SOKO-ENTERPRISE-LIVE",
				DigitalDownloadUrl: "https://cdn.soko.global/releases/soko-erp-v5.3.zip",
			},
		},
		CreatedAt: time.Now().Add(-3 * 24 * time.Hour),
	}
	s.orders[order1.ID] = order1
}

// User methods
func (s *Store) Authenticate(email, password string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[strings.ToLower(strings.TrimSpace(email))]
	if !exists || user.Password != password {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}

func (s *Store) RegisterUser(user *models.User) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	email := strings.ToLower(strings.TrimSpace(user.Email))
	if _, exists := s.users[email]; exists {
		return nil, errors.New("an account with this email already exists")
	}

	user.ID = fmt.Sprintf("usr-%d", time.Now().UnixNano())
	user.CreatedAt = time.Now()
	user.Email = email
	if user.Avatar == "" {
		user.Avatar = fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=0D8ABC&color=fff", strings.ReplaceAll(user.Name, " ", "+"))
	}
	s.users[email] = user
	return user, nil
}

func (s *Store) GetUserByID(id string) *models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.ID == id {
			return u
		}
	}
	return nil
}

func (s *Store) GetUsers() []*models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.User, 0, len(s.users))
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

// Product Filtering & Operations
func (s *Store) GetProducts(filterType, format, region, category, condition, query string) []*models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*models.Product, 0)
	query = strings.ToLower(strings.TrimSpace(query))
	filterType = strings.ToLower(strings.TrimSpace(filterType))
	format = strings.ToLower(strings.TrimSpace(format))
	region = strings.ToLower(strings.TrimSpace(region))
	category = strings.ToLower(strings.TrimSpace(category))
	condition = strings.ToLower(strings.TrimSpace(condition))

	for _, p := range s.products {
		// Type filter: physical vs digital
		if filterType != "" && filterType != "all" {
			if filterType == "digital" && p.Type == models.TypePhysical {
				continue
			}
			if filterType == "physical" && p.Type != models.TypePhysical {
				continue
			}
		}

		// Format filter: wholesale, retail, digital, make_offer
		if format != "" && format != "all" {
			if string(p.Format) != format {
				continue
			}
		}

		// Region filter
		if region != "" && region != "global" && region != "all" {
			if p.VendorRegion != "global" && p.VendorRegion != region {
				continue
			}
		}

		// Category filter
		if category != "" && category != "all" {
			if !strings.EqualFold(p.Category, category) {
				continue
			}
		}

		// Condition filter
		if condition != "" && condition != "all" {
			if string(p.Condition) != condition {
				continue
			}
		}

		// Search query filter
		if query != "" {
			match := strings.Contains(strings.ToLower(p.Title), query) ||
				strings.Contains(strings.ToLower(p.Description), query) ||
				strings.Contains(strings.ToLower(p.Sku), query) ||
				strings.Contains(strings.ToLower(p.Category), query) ||
				strings.Contains(strings.ToLower(p.VendorName), query) ||
				strings.Contains(strings.ToLower(p.VendorCountry), query)
			if !match {
				continue
			}
		}

		list = append(list, p)
	}
	return list
}

func (s *Store) GetProductByID(id string) *models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.products[id]
}

func (s *Store) AddProduct(p *models.Product) *models.Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p.ID == "" {
		p.ID = fmt.Sprintf("prod-%d", time.Now().UnixNano())
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.Rating == 0 {
		p.Rating = 5.0
	}
	if p.VendorPositiveRate == "" {
		p.VendorPositiveRate = "100%"
	}
	if p.TradeAssuranceCover == 0 {
		p.TradeAssuranceCover = 500000.00
	}
	if p.Stock == 0 && p.Type != models.TypePhysical {
		p.Stock = -1
	}
	if p.ImageUrl == "" {
		p.ImageUrl = "https://images.unsplash.com/photo-1581092160607-ee22621dd758?w=800&auto=format&fit=crop&q=80"
	}
	s.products[p.ID] = p
	return p
}

func (s *Store) DeleteProduct(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.products[id]; exists {
		delete(s.products, id)
		return true
	}
	return false
}

// RFQ Operations
func (s *Store) GetRfqs() []*models.Rfq {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.Rfq, 0, len(s.rfqs))
	for _, r := range s.rfqs {
		list = append(list, r)
	}
	return list
}

func (s *Store) CreateRfq(rfq *models.Rfq) *models.Rfq {
	s.mu.Lock()
	defer s.mu.Unlock()

	rfq.ID = fmt.Sprintf("rfq-%d", time.Now().UnixNano())
	rfq.RfqNumber = fmt.Sprintf("RFQ-%d-%03d", time.Now().Year(), time.Now().Unix()%1000)
	rfq.Status = models.RfqStatusOpen
	rfq.CreatedAt = time.Now()
	rfq.Quotes = make([]models.RfqQuote, 0)

	s.rfqs[rfq.ID] = rfq
	return rfq
}

func (s *Store) AddRfqQuote(rfqID string, quote models.RfqQuote) (*models.Rfq, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rfq, exists := s.rfqs[rfqID]
	if !exists {
		return nil, errors.New("RFQ not found")
	}

	quote.ID = fmt.Sprintf("quot-%d", time.Now().UnixNano())
	quote.RfqID = rfqID
	quote.CreatedAt = time.Now()
	quote.Status = "pending"

	rfq.Quotes = append(rfq.Quotes, quote)
	return rfq, nil
}

// Offers Operations (eBay Make Offer)
func (s *Store) SubmitOffer(offer *models.Offer) *models.Offer {
	s.mu.Lock()
	defer s.mu.Unlock()

	offer.ID = fmt.Sprintf("off-%d", time.Now().UnixNano())
	offer.CreatedAt = time.Now()
	offer.Status = models.OfferStatusPending
	offer.TotalOffered = offer.OfferedPrice * float64(offer.Quantity)

	s.offers[offer.ID] = offer
	return offer
}

func (s *Store) GetOffers() []*models.Offer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*models.Offer, 0, len(s.offers))
	for _, o := range s.offers {
		list = append(list, o)
	}
	return list
}

// Orders Operations
func (s *Store) GetOrders() []*models.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]*models.Order, 0, len(s.orders))
	for _, o := range s.orders {
		orders = append(orders, o)
	}
	return orders
}

func (s *Store) CreateOrder(order *models.Order) *models.Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	order.ID = fmt.Sprintf("ord-%d", time.Now().UnixNano())
	order.OrderNumber = fmt.Sprintf("SOKO-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)
	order.CreatedAt = time.Now()
	order.Status = models.OrderStatusCompleted

	for i, item := range order.Items {
		if item.ProductType != models.TypePhysical {
			order.HasDigital = true
			if item.DigitalLicenseKey == "" {
				order.Items[i].DigitalLicenseKey = fmt.Sprintf("SOKO-LIC-%d-%s", time.Now().Unix()%100000, strings.ToUpper(item.ProductID))
			}
			if item.DigitalDownloadUrl == "" {
				order.Items[i].DigitalDownloadUrl = fmt.Sprintf("https://cdn.soko.global/download/%s?token=%d", item.ProductID, time.Now().UnixNano())
			}
		} else {
			order.HasPhysical = true
		}
	}

	s.orders[order.ID] = order
	return order
}

// Metadata getters
func (s *Store) GetCurrencies() []models.Currency {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]models.Currency, 0, len(s.currencies))
	for _, c := range s.currencies {
		list = append(list, c)
	}
	return list
}

func (s *Store) GetRegions() []models.MarketRegion {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]models.MarketRegion, 0, len(s.regions))
	for _, r := range s.regions {
		list = append(list, r)
	}
	return list
}

type DashboardStats struct {
	TotalGMV              float64 `json:"total_gmv"`
	TotalOrders           int     `json:"total_orders"`
	DigitalProductsCount  int     `json:"digital_products_count"`
	PhysicalProductsCount int     `json:"physical_products_count"`
	TotalMerchants        int     `json:"total_merchants"`
	TotalBuyers           int     `json:"total_buyers"`
	InstantDeliveries     int     `json:"instant_deliveries"`
	ActiveRfqsCount       int     `json:"active_rfqs_count"`
	ActiveOffersCount     int     `json:"active_offers_count"`
	TradeAssuranceReserves float64 `json:"trade_assurance_reserves"`
}

func (s *Store) GetDashboardStats() DashboardStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var gmv float64
	instantFulfilled := 0

	for _, o := range s.orders {
		gmv += o.TotalAmount
		if o.HasDigital {
			instantFulfilled++
		}
	}

	var digCount, phyCount int
	for _, p := range s.products {
		if p.Type == models.TypePhysical {
			phyCount++
		} else {
			digCount++
		}
	}

	var merchants, buyers int
	for _, u := range s.users {
		if u.Role == models.RoleMerchant {
			merchants++
		} else if u.Role == models.RoleBuyer {
			buyers++
		}
	}

	return DashboardStats{
		TotalGMV:              gmv,
		TotalOrders:           len(s.orders),
		DigitalProductsCount:  digCount,
		PhysicalProductsCount: phyCount,
		TotalMerchants:        merchants,
		TotalBuyers:           buyers,
		InstantDeliveries:     instantFulfilled,
		ActiveRfqsCount:       len(s.rfqs),
		ActiveOffersCount:     len(s.offers),
		TradeAssuranceReserves: 14250000.00,
	}
}
