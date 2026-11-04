package services

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"encoding/json"
)

func LoadMediaPolicy() MediaPolicy {
	p := DefaultMediaPolicy()
	if !facades.Schema().HasTable("storefront_contents") {
		return p
	}
	var row models.StorefrontContent
	if facades.Orm().Query().Where("id = ?", "media-policy").First(&row) == nil && row.ID != "" {
		var candidate MediaPolicy
		if json.Unmarshal([]byte(row.ContentJSON), &candidate) == nil && ValidateMediaPolicy(candidate) == nil {
			return candidate
		}
	}
	return p
}
