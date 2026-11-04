package queries

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"github.com/goravel/framework/contracts/database/orm"
	"time"
)

// ActivePromotions uses one captured instant. Expiration is enforced at every
// price calculation and checkout, without relying on an asynchronous scheduler.
func ActivePromotions(query orm.Query, now time.Time) ([]models.Promotion, error) {
	campaigns := []models.Promotion{}
	err := query.Where("enabled = ? AND starts_at <= ? AND ends_at > ?", true, now, now).Get(&campaigns)
	for i := range campaigns {
		campaigns[i].Hydrate()
	}
	return campaigns, err
}
func ApplyProductPricing(products []models.Product) error {
	now := time.Now()
	campaigns, err := ActivePromotions(facades.Orm().Query(), now)
	if err != nil {
		return err
	}
	for i := range products {
		qty := max(1, products[i].MOQ)
		p, e := models.PriceProduct(&products[i], qty, campaigns, now)
		if e == nil {
			products[i].Pricing = &p
		}
	}
	return nil
}
