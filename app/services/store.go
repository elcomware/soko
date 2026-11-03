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
	mu       sync.RWMutex
	users    map[string]*models.User
	products map[string]*models.Product
	orders   map[string]*models.Order
}

var (
	globalStore *Store
	once        sync.Once
)

func GetStore() *Store {
	once.Do(func() {
		globalStore = &Store{
			users:    make(map[string]*models.User),
			products: make(map[string]*models.Product),
			orders:   make(map[string]*models.Order),
		}
		globalStore.seed()
	})
	return globalStore
}

func (s *Store) seed() {
	// Seed Users: Admin, Merchant/Supplier, Buyer
	admin := &models.User{
		ID:          "usr-admin-1",
		Name:        "Soko Global Admin",
		Email:       "admin@soko.global",
		Password:    "admin123",
		Role:        models.RoleAdmin,
		CompanyName: "Soko Marketplace Holdings Inc.",
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

	// Seed Products (Physical B2B/Retail & Digital Downloads/SaaS/Licenses)
	products := []*models.Product{
		// 1. Digital Enterprise Software
		{
			ID:                   "prod-dig-1",
			Title:                "Soko Cloud ERP Enterprise Suite v5",
			Slug:                 "soko-cloud-erp-enterprise",
			Description:          "Full-stack multi-warehouse supply chain and enterprise resource planning software with real-time financial tracking, multi-currency ledger, and REST/GraphQL APIs.",
			ShortDesc:            "Tier-1 enterprise ERP with automated billing, logistics, and multi-tenant ledger.",
			Category:             "Digital Software",
			SubCategory:          "Enterprise SaaS",
			Type:                 models.TypeDigitalSoftware,
			Price:                2499.00,
			WholesalePrice:       1999.00,
			MOQ:                  1,
			Stock:                -1, // Unlimited
			Sku:                  "SKU-SOK-ERP5",
			VendorID:             merchant.ID,
			VendorName:           merchant.Name,
			VendorCountry:        merchant.Country,
			VendorVerified:       true,
			Rating:               4.9,
			ReviewsCount:         128,
			ImageUrl:             "https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800&auto=format&fit=crop&q=80",
			Tags:                 []string{"SaaS", "ERP", "Cloud", "Finance", "Instant Delivery"},
			Featured:             true,
			LeadTimeDays:         0,
			DigitalDownloadUrl:   "https://cdn.soko.global/releases/soko-erp-v5.3.zip",
			DigitalLicenseFormat: "ERP5-XXXX-XXXX-ENTERPRISE",
			DigitalFileSize:      "1.42 GB",
			CreatedAt:            time.Now().Add(-45 * 24 * time.Hour),
		},
		// 2. Digital Machine Learning Algorithm / Weights
		{
			ID:                   "prod-dig-2",
			Title:                "Global Route & Freight Optimization Neural Model",
			Slug:                 "freight-route-opt-ai-model",
			Description:          "Pretrained ONNX and PyTorch weights fine-tuned on 100M+ maritime and air-cargo shipment routes to minimize fuel consumption and customs clearance latency.",
			ShortDesc:            "Industrial AI neural network for container vessel routing and air freight scheduling.",
			Category:             "Digital Assets",
			SubCategory:          "AI & Machine Learning",
			Type:                 models.TypeDigitalAsset,
			Price:                680.00,
			WholesalePrice:       540.00,
			MOQ:                  1,
			Stock:                -1,
			Sku:                  "SKU-AI-ROUTEX",
			VendorID:             merchant.ID,
			VendorName:           merchant.Name,
			VendorCountry:        merchant.Country,
			VendorVerified:       true,
			Rating:               4.8,
			ReviewsCount:         64,
			ImageUrl:             "https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?w=800&auto=format&fit=crop&q=80",
			Tags:                 []string{"AI", "PyTorch", "Logistics", "Digital Download"},
			Featured:             true,
			LeadTimeDays:         0,
			DigitalDownloadUrl:   "https://cdn.soko.global/weights/freight-opt-v2.pt.tar.gz",
			DigitalLicenseFormat: "AI-MODEL-XXXX-XXXX-XXXX",
			DigitalFileSize:      "890 MB",
			CreatedAt:            time.Now().Add(-30 * 24 * time.Hour),
		},
		// 3. Digital Financial Gateway License
		{
			ID:                   "prod-dig-3",
			Title:                "Cross-Border High-Throughput FX Settlement API Gateway",
			Slug:                 "fx-settlement-api-gateway",
			Description:          "Ultra-low latency microservice package in Go + Rust for multi-currency automated forex clearing, SWIFT ISO20022 message parsing, and SEPA/ACH integrations.",
			ShortDesc:            "High-throughput banking gateway for cross-border treasury settlements.",
			Category:             "Digital Software",
			SubCategory:          "Fintech APIs",
			Type:                 models.TypeDigitalLicense,
			Price:                1250.00,
			WholesalePrice:       990.00,
			MOQ:                  1,
			Stock:                -1,
			Sku:                  "SKU-FIN-GATE88",
			VendorID:             merchant.ID,
			VendorName:           merchant.Name,
			VendorCountry:        merchant.Country,
			VendorVerified:       true,
			Rating:               5.0,
			ReviewsCount:         39,
			ImageUrl:             "https://images.unsplash.com/photo-1559526324-4b87b5e36e44?w=800&auto=format&fit=crop&q=80",
			Tags:                 []string{"Fintech", "API", "Golang", "Microservice"},
			Featured:             false,
			LeadTimeDays:         0,
			DigitalDownloadUrl:   "https://cdn.soko.global/binaries/fx-gateway-darwin-linux-amd64.tar.gz",
			DigitalLicenseFormat: "FX-GATE-PRO-XXXX-XXXX",
			DigitalFileSize:      "185 MB",
			CreatedAt:            time.Now().Add(-20 * 24 * time.Hour),
		},
		// 4. Digital Industrial CAD / 3D Asset
		{
			ID:                   "prod-dig-4",
			Title:                "Precision Modular Solar Inverter 50kW CAD Blueprints",
			Slug:                 "modular-solar-inverter-50kw-cad",
			Description:          "Complete production-ready SolidWorks, STEP, and Gerber PCB manufacturing schematic package for high-efficiency grid-tied 50kW photovoltaic inverters.",
			ShortDesc:            "Factory-ready CAD, thermal analysis and PCB design files for 50kW solar inverters.",
			Category:             "Digital Assets",
			SubCategory:          "Engineering & 3D Blueprints",
			Type:                 models.TypeDigitalAsset,
			Price:                349.00,
			WholesalePrice:       280.00,
			MOQ:                  1,
			Stock:                -1,
			Sku:                  "SKU-ENG-SOLAR50",
			VendorID:             merchant.ID,
			VendorName:           merchant.Name,
			VendorCountry:        merchant.Country,
			VendorVerified:       true,
			Rating:               4.7,
			ReviewsCount:         22,
			ImageUrl:             "https://images.unsplash.com/photo-1509391365360-2e959784a276?w=800&auto=format&fit=crop&q=80",
			Tags:                 []string{"CAD", "Solar", "Engineering", "Schematics"},
			Featured:             false,
			LeadTimeDays:         0,
			DigitalDownloadUrl:   "https://cdn.soko.global/blueprints/solar-50kw-complete-cad.zip",
			DigitalLicenseFormat: "CAD-IP-ROYALTY-FREE-XXXX",
			DigitalFileSize:      "620 MB",
			CreatedAt:            time.Now().Add(-15 * 24 * time.Hour),
		},
		// 5. Physical B2B Heavy Machinery
		{
			ID:                   "prod-phy-1",
			Title:                "Industrial High-Speed CNC Fiber Laser Cutting Machine 3000W",
			Slug:                 "industrial-cnc-fiber-laser-3000w",
			Description:          "Heavy-duty aerospace-grade welded steel bed with Raycus laser source, CypCut CNC operating system, auto-exchange dual shuttle tables, and CE certification for precision sheet metal fabrication.",
			ShortDesc:            "3000W industrial fiber laser for 22mm carbon steel and 12mm stainless steel.",
			Category:             "Industrial & Machinery",
			SubCategory:          "Fabrication & CNC",
			Type:                 models.TypePhysical,
			Price:                18500.00,
			WholesalePrice:       15200.00,
			MOQ:                  1,
			Stock:                15,
			Sku:                  "SKU-IND-LSR3000",
			VendorID:             merchant.ID,
			VendorName:           "Bavaria Precision Machine Works",
			VendorCountry:        "Germany",
			VendorVerified:       true,
			Rating:               4.9,
			ReviewsCount:         78,
			ImageUrl:             "https://images.unsplash.com/photo-1581092160607-ee22621dd758?w=800&auto=format&fit=crop&q=80",
			Tags:                 []string{"Industrial", "CNC", "Laser", "Machinery", "B2B Freight"},
			Featured:             true,
			LeadTimeDays:         14,
			CreatedAt:            time.Now().Add(-60 * 24 * time.Hour),
		},
		// 6. Physical Wholesale Electronics / Networking
		{
			ID:                   "prod-phy-2",
			Title:                "Enterprise Shielded Cat6A 10Gbps Bare Copper Cable (1000ft Bulk Spool)",
			Slug:                 "enterprise-cat6a-shielded-1000ft",
			Description:          "Solid 23AWG pure oxygen-free copper, CMR riser rated, tested up to 650MHz bandwidth with foil shielding (F/UTP) for zero cross-talk in datacenter and commercial installations.",
			ShortDesc:            "Datacenter-certified 10Gbps Ethernet cable spools for commercial enterprise cabling.",
			Category:             "Electronics & Telecom",
			SubCategory:          "Networking Infrastructure",
			Type:                 models.TypePhysical,
			Price:                145.00,
			WholesalePrice:       98.00,
			MOQ:                  10, // 10 Spools minimum
			Stock:                420,
			Sku:                  "SKU-NET-C6A-1000",
			VendorID:             merchant.ID,
			VendorName:           "Shenzhen Global Opto-Cable Ltd",
			VendorCountry:        "China",
			VendorVerified:       true,
			Rating:               4.8,
			ReviewsCount:         312,
			ImageUrl:             "https://images.unsplash.com/photo-1544197150-b99a580bb7a8?w=800&auto=format&fit=crop&q=80",
			Tags:                 []string{"Datacenter", "Ethernet", "Wholesale", "Hardware"},
			Featured:             true,
			LeadTimeDays:         5,
			CreatedAt:            time.Now().Add(-40 * 24 * time.Hour),
		},
		// 7. Physical Commercial Robotics & Automation
		{
			ID:                   "prod-phy-3",
			Title:                "Autonomous Mobile Warehouse Robot (AMR) 500kg Payload",
			Slug:                 "amr-warehouse-robot-500kg",
			Description:          "LiDAR SLAM navigation without magnetic floor tapes, 360-degree obstacle detection, smart swarm fleet coordination algorithm, and rapid 20-minute wireless contact charging.",
			ShortDesc:            "Next-gen automated guided vehicle for smart e-commerce logistics centers.",
			Category:             "Industrial & Machinery",
			SubCategory:          "Warehouse Robotics",
			Type:                 models.TypePhysical,
			Price:                9400.00,
			WholesalePrice:       7800.00,
			MOQ:                  2,
			Stock:                28,
			Sku:                  "SKU-ROB-AMR500",
			VendorID:             merchant.ID,
			VendorName:           "Tokyo Automation Dynamics",
			VendorCountry:        "Japan",
			VendorVerified:       true,
			Rating:               4.9,
			ReviewsCount:         45,
			ImageUrl:             "https://images.unsplash.com/photo-1485827404703-89b55fcc595e?w=800&auto=format&fit=crop&q=80",
			Tags:                 []string{"Robotics", "AMR", "Automation", "Warehouse"},
			Featured:             false,
			LeadTimeDays:         10,
			CreatedAt:            time.Now().Add(-25 * 24 * time.Hour),
		},
		// 8. Physical Commercial Smart Energy / Battery
		{
			ID:                   "prod-phy-4",
			Title:                "Commercial LiFePO4 Energy Storage Cabinet 100kWh 400V",
			Slug:                 "commercial-lifepo4-energy-storage-100kwh",
			Description:          "Scalable containerized lithium iron phosphate battery system with built-in aerosol fire suppression, active liquid cooling, 6000+ cycle life, and CAN/RS485 smart BMS integration.",
			ShortDesc:            "100kWh industrial power backup and peak-shaving energy storage unit.",
			Category:             "Energy & Utilities",
			SubCategory:          "Commercial Storage",
			Type:                 models.TypePhysical,
			Price:                22000.00,
			WholesalePrice:       18500.00,
			MOQ:                  1,
			Stock:                8,
			Sku:                  "SKU-NRG-BATT100",
			VendorID:             merchant.ID,
			VendorName:           "Nordic Clean Energy Alliance",
			VendorCountry:        "Sweden",
			VendorVerified:       true,
			Rating:               5.0,
			ReviewsCount:         18,
			ImageUrl:             "https://images.unsplash.com/photo-1513836279014-a89f7a76ae86?w=800&auto=format&fit=crop&q=80",
			Tags:                 []string{"Energy", "Storage", "LiFePO4", "Commercial"},
			Featured:             true,
			LeadTimeDays:         21,
			CreatedAt:            time.Now().Add(-10 * 24 * time.Hour),
		},
	}

	for _, p := range products {
		s.products[p.ID] = p
	}

	// Seed Sample Orders (Showing completed orders with digital keys and physical shipments)
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

	order2 := &models.Order{
		ID:            "ord-9822-eu",
		OrderNumber:   "SOKO-2026-9822",
		UserID:        buyer.ID,
		CustomerName:  buyer.Name,
		CustomerEmail: buyer.Email,
		TotalAmount:   18500.00,
		Currency:      "USD",
		Status:        models.OrderStatusShipped,
		PaymentMethod: "Letter of Credit (L/C)",
		HasDigital:    false,
		HasPhysical:   true,
		Items: []models.OrderItem{
			{
				ProductID:    "prod-phy-1",
				ProductTitle: "Industrial High-Speed CNC Fiber Laser Cutting Machine 3000W",
				ProductType:  models.TypePhysical,
				Quantity:     1,
				UnitPrice:    18500.00,
				TotalPrice:   18500.00,
			},
		},
		CreatedAt: time.Now().Add(-5 * 24 * time.Hour),
	}
	s.orders[order2.ID] = order2
}

func (s *Store) Authenticate(email, password string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[strings.ToLower(strings.TrimSpace(email))]
	if !exists {
		return nil, errors.New("invalid email or password")
	}
	if user.Password != password {
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

func (s *Store) GetProducts(filterType, category, query string) []*models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*models.Product, 0)
	query = strings.ToLower(strings.TrimSpace(query))
	filterType = strings.ToLower(strings.TrimSpace(filterType))
	category = strings.ToLower(strings.TrimSpace(category))

	for _, p := range s.products {
		// Filter by Type: "physical", "digital", or specific
		if filterType != "" && filterType != "all" {
			if filterType == "digital" {
				if p.Type == models.TypePhysical {
					continue
				}
			} else if filterType == "physical" {
				if p.Type != models.TypePhysical {
					continue
				}
			} else if string(p.Type) != filterType {
				continue
			}
		}

		// Filter by category
		if category != "" && category != "all" {
			if !strings.EqualFold(p.Category, category) {
				continue
			}
		}

		// Filter by search query
		if query != "" {
			match := strings.Contains(strings.ToLower(p.Title), query) ||
				strings.Contains(strings.ToLower(p.Description), query) ||
				strings.Contains(strings.ToLower(p.Sku), query) ||
				strings.Contains(strings.ToLower(p.Category), query) ||
				strings.Contains(strings.ToLower(p.VendorName), query)
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
	if p.Stock == 0 && p.Type != models.TypePhysical {
		p.Stock = -1 // unlimited for digital
	}
	if p.ImageUrl == "" {
		if p.Type == models.TypePhysical {
			p.ImageUrl = "https://images.unsplash.com/photo-1581092160607-ee22621dd758?w=800&auto=format&fit=crop&q=80"
		} else {
			p.ImageUrl = "https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800&auto=format&fit=crop&q=80"
		}
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

	// Process digital licenses
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

type DashboardStats struct {
	TotalGMV             float64 `json:"total_gmv"`
	TotalOrders          int     `json:"total_orders"`
	DigitalProductsCount int     `json:"digital_products_count"`
	PhysicalProductsCount int    `json:"physical_products_count"`
	TotalMerchants       int     `json:"total_merchants"`
	TotalBuyers          int     `json:"total_buyers"`
	InstantDeliveries    int     `json:"instant_deliveries"`
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
		TotalGMV:             gmv,
		TotalOrders:          len(s.orders),
		DigitalProductsCount: digCount,
		PhysicalProductsCount: phyCount,
		TotalMerchants:       merchants,
		TotalBuyers:          buyers,
		InstantDeliveries:    instantFulfilled,
	}
}

