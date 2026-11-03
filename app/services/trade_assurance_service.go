package services

type TradeAssuranceService struct{}

func NewTradeAssuranceService() *TradeAssuranceService {
	return &TradeAssuranceService{}
}

// CalculateProtectionCoverage computes escrow coverage cap based on product & order value
func (s *TradeAssuranceService) CalculateProtectionCoverage(orderValue float64, supplierCover float64) float64 {
	if supplierCover <= 0 {
		return 100000.00 // Standard baseline protection: $100,000
	}
	if orderValue > supplierCover {
		return supplierCover
	}
	return orderValue
}

// EscrowFeeRate returns standard escrow guarantee percentage (typically 0% for qualified buyers, 1.5% for extended warranty)
func (s *TradeAssuranceService) EscrowFeeRate(withExtendedInspection bool) float64 {
	if withExtendedInspection {
		return 0.015 // 1.5% for third-party lab inspection certificate
	}
	return 0.0 // 0% standard Trade Assurance
}

// DisputeWindowDays returns buyer inspection and dispute window
func (s *TradeAssuranceService) DisputeWindowDays(isDigital bool) int {
	if isDigital {
		return 7 // 7 days for software license activation / token revocation
	}
	return 30 // 30 days post-customs delivery for physical wholesale inspection
}

