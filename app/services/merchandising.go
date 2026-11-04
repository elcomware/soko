package services

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/queries"
	"encoding/json"
	"errors"
	"github.com/goravel/framework/contracts/database/orm"
)

func SaveMerchandising(actor *models.User, settings models.MerchandisingSettings) error {
	if !actor.Can("content.update", "") {
		return errors.New("storefront content permission required")
	}
	for _, placement := range []*models.MerchandisingPlacement{&settings.Hero, &settings.HeroSide, &settings.Discovery} {
		if !placement.Enabled || placement.Strategy != "manual" {
			placement.ProductIDs = []string{}
		}
	}
	if err := queries.ValidateMerchandising(settings); err != nil {
		return err
	}
	for index, placement := range []models.MerchandisingPlacement{settings.Hero, settings.HeroSide, settings.Discovery} {
		for _, id := range placement.ProductIDs {
			var product models.Product
			if err := facades.Orm().Query().Where("id = ?", id).First(&product); err != nil {
				return err
			}
			if product.ID == "" {
				return errors.New("a pinned product no longer exists")
			}
			if index == 2 && !product.IsDigital() {
				return errors.New("digital discovery pins must be digital products")
			}
			if placement.Category != "" && product.Category != placement.Category {
				return errors.New("a pinned product does not match its placement category")
			}
		}
	}
	encoded, _ := json.Marshal(settings)
	return facades.Orm().Transaction(func(tx orm.Query) error {
		var record models.MerchandisingRecord
		if err := tx.Where("id = ?", "store").First(&record); err != nil {
			return err
		}
		if record.Version != settings.Version {
			return errors.New("storefront settings changed; reload before saving")
		}
		if record.ID == "" {
			return tx.Create(&models.MerchandisingRecord{ID: "store", SettingsJSON: string(encoded), Version: 1})
		}
		result, err := tx.Model(&models.MerchandisingRecord{}).Where("id = ? AND version = ?", "store", settings.Version).Update(map[string]any{"settings_json": string(encoded), "version": settings.Version + 1})
		if err == nil && result.RowsAffected != 1 {
			return errors.New("storefront settings changed; reload before saving")
		}
		return err
	})
}
