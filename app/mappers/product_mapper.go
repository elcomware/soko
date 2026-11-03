package mappers

import (
	"encoding/json"
	"time"

	"elcomware/app/dto"
	"elcomware/app/models"
)

type ProductMapper struct{}

func NewProductMapper() *ProductMapper {
	return &ProductMapper{}
}

func (m *ProductMapper) ToModel(input dto.CreateProductInput) models.Product {
	tiersBytes, _ := json.Marshal(input.PriceTiers)
	galleryBytes, _ := json.Marshal(input.GalleryUrls)
	tagsBytes, _ := json.Marshal(input.Tags)

	return models.Product{
		Title:                input.Title,
		Slug:                 input.Slug,
		Description:          input.Description,
		ShortDesc:            input.ShortDesc,
		Category:             input.Category,
		SubCategory:          input.SubCategory,
		Type:                 input.Type,
		Format:               input.Format,
		Condition:            input.Condition,
		Price:                input.Price,
		WholesalePrice:       input.WholesalePrice,
		PriceTiersJSON:       string(tiersBytes),
		MOQ:                  input.MOQ,
		Stock:                input.Stock,
		Sku:                  input.Sku,
		VendorID:             input.VendorID,
		VendorName:           input.VendorName,
		VendorCountry:        input.VendorCountry,
		VendorRegion:         input.VendorRegion,
		VendorVerified:       input.VendorVerified,
		VendorBadge:          input.VendorBadge,
		VendorPositiveRate:   input.VendorPositiveRate,
		TradeAssuranceCover:  input.TradeAssuranceCover,
		AcceptsOffers:        input.AcceptsOffers,
		Rating:               5.0,
		ReviewsCount:         0,
		ImageUrl:             input.ImageUrl,
		GalleryUrlsJSON:      string(galleryBytes),
		TagsJSON:             string(tagsBytes),
		Featured:             input.Featured,
		LeadTimeDays:         input.LeadTimeDays,
		Incoterm:             input.Incoterm,
		LocalPickupAvailable: input.LocalPickupAvailable,
		DigitalDownloadUrl:   input.DigitalDownloadUrl,
		DigitalLicenseFormat: input.DigitalLicenseFormat,
		DigitalFileSize:      input.DigitalFileSize,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

func (m *ProductMapper) ToDTO(model models.Product) dto.ProductOutput {
	var priceTiers []models.PriceTier
	if model.PriceTiersJSON != "" {
		_ = json.Unmarshal([]byte(model.PriceTiersJSON), &priceTiers)
	}

	var galleryUrls []string
	if model.GalleryUrlsJSON != "" {
		_ = json.Unmarshal([]byte(model.GalleryUrlsJSON), &galleryUrls)
	}

	var tags []string
	if model.TagsJSON != "" {
		_ = json.Unmarshal([]byte(model.TagsJSON), &tags)
	}

	return dto.ProductOutput{
		ID:                   model.ID,
		Title:                model.Title,
		Slug:                 model.Slug,
		Description:          model.Description,
		ShortDesc:            model.ShortDesc,
		Category:             model.Category,
		SubCategory:          model.SubCategory,
		Type:                 model.Type,
		Format:               model.Format,
		Condition:            model.Condition,
		Price:                model.Price,
		WholesalePrice:       model.WholesalePrice,
		PriceTiers:           priceTiers,
		MOQ:                  model.MOQ,
		Stock:                model.Stock,
		Sku:                  model.Sku,
		VendorID:             model.VendorID,
		VendorName:           model.VendorName,
		VendorCountry:        model.VendorCountry,
		VendorRegion:         model.VendorRegion,
		VendorVerified:       model.VendorVerified,
		VendorBadge:          model.VendorBadge,
		VendorPositiveRate:   model.VendorPositiveRate,
		TradeAssuranceCover:  model.TradeAssuranceCover,
		AcceptsOffers:        model.AcceptsOffers,
		Rating:               model.Rating,
		ReviewsCount:         model.ReviewsCount,
		ImageUrl:             model.ImageUrl,
		GalleryUrls:          galleryUrls,
		Tags:                 tags,
		Featured:             model.Featured,
		LeadTimeDays:         model.LeadTimeDays,
		Incoterm:             model.Incoterm,
		LocalPickupAvailable: model.LocalPickupAvailable,
		DigitalDownloadUrl:   model.DigitalDownloadUrl,
		DigitalLicenseFormat: model.DigitalLicenseFormat,
		DigitalFileSize:      model.DigitalFileSize,
		CreatedAt:            model.CreatedAt,
	}
}

func (m *ProductMapper) ToDTOList(list []models.Product) []dto.ProductOutput {
	dtos := make([]dto.ProductOutput, 0, len(list))
	for _, item := range list {
		dtos = append(dtos, m.ToDTO(item))
	}
	return dtos
}
