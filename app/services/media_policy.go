package services

import (
	"crypto/sha256"
	"elcomware/app/facades"
	"elcomware/app/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
			if candidate.Version == "" {
				candidate.Version = legacyMediaVersion(row.ContentJSON)
			}
			return candidate
		}
	}
	return p
}

var ErrMediaPolicyConflict = errors.New("Image policy changed since you opened it. Reload the current policy and review your changes before saving.")

func legacyMediaVersion(raw string) string {
	return fmt.Sprintf("legacy-%x", sha256.Sum256([]byte(raw)))
}

// SaveMediaPolicy atomically compares the acknowledged record before replacing it.
// Existing unversioned records remain readable; clients must reload before editing.
func SaveMediaPolicy(p MediaPolicy) (MediaPolicy, error) {
	if err := ValidateMediaPolicy(p); err != nil {
		return p, err
	}
	var current models.StorefrontContent
	if err := facades.Orm().Query().Where("id = ?", "media-policy").First(&current); err != nil {
		return p, err
	}
	expected := p.Version
	if current.ID != "" {
		var previous MediaPolicy
		if err := json.Unmarshal([]byte(current.ContentJSON), &previous); err != nil {
			return p, err
		}
		actual := previous.Version
		if actual == "" {
			actual = legacyMediaVersion(current.ContentJSON)
		}
		if expected != actual {
			return p, ErrMediaPolicyConflict
		}
	} else if expected != "" {
		return p, ErrMediaPolicyConflict
	}
	p.Version = uuid.NewString()
	raw, err := json.Marshal(p)
	if err != nil {
		return p, err
	}
	if current.ID == "" {
		err = facades.Orm().Query().Create(&models.StorefrontContent{ID: "media-policy", ContentJSON: string(raw)})
		if err != nil {
			var existing models.StorefrontContent
			if facades.Orm().Query().Where("id = ?", "media-policy").First(&existing) == nil && existing.ID != "" {
				return p, ErrMediaPolicyConflict
			}
			return p, err
		}
	} else {
		result, updateErr := facades.Orm().Query().Model(&models.StorefrontContent{}).Where("id = ? AND content_json = ?", "media-policy", current.ContentJSON).Update(map[string]any{"content_json": string(raw)})
		if updateErr != nil {
			return p, updateErr
		}
		if result.RowsAffected != 1 {
			return p, ErrMediaPolicyConflict
		}
	}
	return p, nil
}
