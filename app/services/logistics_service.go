package services

type ShippingEstimate struct {
	Method          string  `json:"method"`
	EstimatedDays   int     `json:"estimated_days"`
	EstimatedCost   float64 `json:"estimated_cost"`
	TrackingCarrier string  `json:"tracking_carrier"`
	Incoterm        string  `json:"incoterm"`
}

type LogisticsService struct{}

func NewLogisticsService() *LogisticsService {
	return &LogisticsService{}
}

func (s *LogisticsService) GetShippingOptions(fromRegion, toRegion string, isDigital bool) []ShippingEstimate {
	if isDigital {
		return []ShippingEstimate{
			{
				Method:          "Instant Cryptographic License & Download Vault",
				EstimatedDays:   0,
				EstimatedCost:   0.00,
				TrackingCarrier: "Automated Soko Vault CDN",
				Incoterm:        "Instant Digital Delivery",
			},
		}
	}

	if fromRegion == toRegion {
		return []ShippingEstimate{
			{
				Method:          "Regional Express Courier (Same Corridor)",
				EstimatedDays:   2,
				EstimatedCost:   45.00,
				TrackingCarrier: "Regional Express Direct",
				Incoterm:        "DDP (Delivered Duty Paid)",
			},
			{
				Method:          "Standard Regional Ground Freight",
				EstimatedDays:   5,
				EstimatedCost:   20.00,
				TrackingCarrier: "National Logistics Fleet",
				Incoterm:        "DDP",
			},
		}
	}

	// Inter-continental / cross-border options
	return []ShippingEstimate{
		{
			Method:          "International Priority Air Cargo",
			EstimatedDays:   4,
			EstimatedCost:   120.00,
			TrackingCarrier: "DHL Global Forwarding / FedEx Cargo",
			Incoterm:        "DDP (Delivered Duty Paid)",
		},
		{
			Method:          "Inter-Continental Ocean Container (LCL/FCL)",
			EstimatedDays:   25,
			EstimatedCost:   35.00,
			TrackingCarrier: "Maersk / MSC Line",
			Incoterm:        "CIF / FOB Available",
		},
	}
}

