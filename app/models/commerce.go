package models

import (
	"errors"
	"math"
)

func (u *User) CanManageCatalog() bool {
	return u.HasPermission("products.create") || u.HasPermission("products.update")
}
func (u *User) CanManageProduct(p *Product) bool {
	return p != nil && u.Can("products.update", p.VendorID)
}
func (p *Product) IsDigital() bool {
	return (p.Type != TypePhysical && p.Type != "") || p.Format == FormatDigital
}
func (p Product) Public() Product {
	p.DigitalDownloadUrl = ""
	p.DigitalLicenseFormat = ""
	p.VendorRegion = NormalizeRegion(p.VendorRegion)
	return p
}
func PublicProducts(products []Product) []Product {
	result := make([]Product, 0, len(products))
	for _, p := range products {
		result = append(result, p.Public())
	}
	return result
}
func (p *Product) PriceFor(quantity int) (float64, error) {
	if quantity < 1 || quantity > 1000000 {
		return 0, errors.New("quantity must be between 1 and 1,000,000")
	}
	minimum := p.MOQ
	if minimum < 1 {
		minimum = 1
	}
	if quantity < minimum {
		return 0, errors.New("quantity is below the listing minimum")
	}
	if p.Stock < -1 || (!p.IsDigital() && p.Stock < 0) || (p.Stock >= 0 && quantity > p.Stock) {
		return 0, errors.New("requested quantity exceeds available stock")
	}
	p.Hydrate()
	price := p.Price
	if len(p.PriceTiers) > 0 {
		for _, tier := range p.PriceTiers {
			if quantity >= tier.MinQty && (tier.MaxQty == 0 || quantity <= tier.MaxQty) {
				price = tier.Price
				break
			}
		}
	} else if p.Format == FormatWholesale && quantity >= 10 && p.WholesalePrice > 0 {
		price = p.WholesalePrice
	}
	if math.IsNaN(price) || math.IsInf(price, 0) || price <= 0 || price > 1e12 {
		return 0, errors.New("listing has no valid price")
	}
	rounded := math.Round(price*100) / 100
	if rounded <= 0 {
		return 0, errors.New("listing price must be at least one cent")
	}
	return rounded, nil
}
