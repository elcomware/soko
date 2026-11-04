package services

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

func ListPromotions(actor *models.User) ([]models.Promotion, error) {
	if actor == nil || !actor.HasPermission("promotions.manage") {
		return nil, errors.New("campaign access denied")
	}
	q := facades.Orm().Query()
	if !actor.Can("promotions.manage", "") {
		q = q.Where("vendor_id = ?", actor.ID)
	}
	result := []models.Promotion{}
	err := q.OrderByDesc("created_at").Get(&result)
	for i := range result {
		result[i].Hydrate()
	}
	return result, err
}
func SavePromotion(input models.Promotion, actor *models.User, version string) (*models.Promotion, error) {
	if actor == nil || !actor.Can("promotions.manage", input.VendorID) {
		return nil, errors.New("campaign access denied")
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	input.Name = strings.TrimSpace(input.Name)
	input.PrepareSave()
	if input.ID == "" {
		input.ID = "promo_" + uuid.NewString()
		input.CreatedBy = actor.ID
		input.CreatedAt = now
		input.UpdatedAt = now
		return &input, facades.Orm().Query().Create(&input)
	}
	var current models.Promotion
	if err := facades.Orm().Query().Where("id = ?", input.ID).First(&current); err != nil {
		return nil, err
	}
	if current.ID == "" || !actor.Can("promotions.manage", current.VendorID) {
		return nil, errors.New("campaign not found or access denied")
	}
	previous, err := time.Parse(time.RFC3339Nano, version)
	if err != nil || !previous.Equal(current.UpdatedAt) {
		return nil, errors.New("campaign changed; refresh before saving")
	}
	input.CreatedBy = current.CreatedBy
	input.CreatedAt = current.CreatedAt
	input.UpdatedAt = now
	values := map[string]any{"name": input.Name, "vendor_id": input.VendorID, "kind": input.Kind, "amount": input.Amount, "enabled": input.Enabled, "starts_at": input.StartsAt, "ends_at": input.EndsAt, "scope_json": input.ScopeJSON, "updated_at": now}
	result, err := facades.Orm().Query().Model(&models.Promotion{}).Where("id = ? AND updated_at = ?", input.ID, current.UpdatedAt).Update(values)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected != 1 {
		return nil, errors.New("campaign changed; refresh before saving")
	}
	return &input, nil
}
