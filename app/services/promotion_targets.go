package services

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"errors"
	"github.com/goravel/framework/contracts/database/orm"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

type PromotionTarget struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	SubCategory string `json:"sub_category"`
	VendorID    string `json:"vendor_id"`
	VendorName  string `json:"vendor_name"`
}
type PromotionTargetFilters struct {
	VendorID, Query string
	Page, PerPage   int
	Scope           models.PromotionScope
}
type PromotionTargetPage struct {
	Page          int   `json:"page"`
	PerPage       int   `json:"per_page"`
	Total         int64 `json:"total"`
	FilteredTotal int64 `json:"filtered_total"`
	LastPage      int   `json:"last_page"`
	From          int   `json:"from"`
	To            int   `json:"to"`
}
type PromotionCategoryFacet struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}
type PromotionSubcategoryFacet struct {
	Category    string `json:"category"`
	SubCategory string `json:"sub_category"`
	Count       int64  `json:"count"`
}
type PromotionSelectedTarget struct {
	ID      string           `json:"id"`
	Status  string           `json:"status"`
	Product *PromotionTarget `json:"product"`
}
type PromotionTargetPreview struct {
	Total    int64             `json:"total"`
	Products []PromotionTarget `json:"products"`
}
type PromotionTargets struct {
	VendorID         string                      `json:"vendor_id"`
	Products         []PromotionTarget           `json:"products"`
	ProductPage      PromotionTargetPage         `json:"productPage"`
	Categories       []PromotionCategoryFacet    `json:"categories"`
	Subcategories    []PromotionSubcategoryFacet `json:"subcategories"`
	SelectedProducts []PromotionSelectedTarget   `json:"selectedProducts"`
	Preview          PromotionTargetPreview      `json:"preview"`
}

func ValidatePromotionTargetFilters(f PromotionTargetFilters) error {
	if f.Page < 1 || f.Page > 1000000 || f.PerPage < 1 || f.PerPage > 100 {
		return errors.New("invalid target pagination")
	}
	if !utf8.ValidString(f.Query) || utf8.RuneCountInString(f.Query) > 120 || len(f.VendorID) > 64 || strings.ContainsAny(f.Query+f.VendorID, "\x00\r\n") {
		return errors.New("invalid target search or merchant")
	}
	// Validate bounded arrays without requiring an unfinished draft to select an
	// inclusion. The original scope, not this validation copy, drives coverage.
	copyScope := f.Scope
	copyScope.All = true
	p := models.Promotion{Name: "Target preview", Kind: "fixed", Amount: 1, StartsAt: time.Unix(1, 0), EndsAt: time.Unix(2, 0), Scope: copyScope}
	return p.Validate()
}
func LoadPromotionTargets(actor *models.User, f PromotionTargetFilters) (*PromotionTargets, error) {
	if actor == nil || !actor.Can("promotions.manage", f.VendorID) {
		return nil, errors.New("campaign target access denied")
	}
	if err := ValidatePromotionTargetFilters(f); err != nil {
		return nil, err
	}
	out := &PromotionTargets{VendorID: f.VendorID, Products: []PromotionTarget{}, Categories: []PromotionCategoryFacet{}, Subcategories: []PromotionSubcategoryFacet{}, SelectedProducts: []PromotionSelectedTarget{}, Preview: PromotionTargetPreview{Products: []PromotionTarget{}}, ProductPage: PromotionTargetPage{Page: f.Page, PerPage: f.PerPage}}
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		q := tx.Model(&models.Product{}).Select("id", "title", "category", "sub_category", "vendor_id", "vendor_name")
		if f.VendorID != "" {
			q = q.Where("vendor_id = ?", f.VendorID)
		}
		selected := map[string]*PromotionTarget{}
		selectedIDs := []string{}
		for _, ids := range [][]string{f.Scope.ProductIDs, f.Scope.ExcludeProductIDs} {
			for _, id := range ids {
				if _, exists := selected[id]; !exists {
					selected[id] = nil
					selectedIDs = append(selectedIDs, id)
				}
			}
		}
		categories := map[string]int64{}
		subcategories := map[models.CategorySelection]int64{}
		includeIDs := targetSet(f.Scope.ProductIDs)
		excludeIDs := targetSet(f.Scope.ExcludeProductIDs)
		includeCategories := targetSet(f.Scope.Categories)
		excludeCategories := targetSet(f.Scope.ExcludeCategories)
		includeSubs := map[models.CategorySelection]bool{}
		excludeSubs := map[models.CategorySelection]bool{}
		for _, v := range f.Scope.SubCategories {
			includeSubs[v] = true
		}
		for _, v := range f.Scope.ExcludeSubCategories {
			excludeSubs[v] = true
		}
		lower := cases.Lower(language.Und)
		needle := lower.String(strings.TrimSpace(f.Query))
		last := []PromotionTarget{}
		var scanErr error
		for row := range q.OrderBy("title").OrderBy("id").Cursor() {
			if row.Err() != nil {
				if scanErr == nil {
					scanErr = row.Err()
				}
				continue
			}
			if scanErr != nil {
				continue
			}
			var p PromotionTarget
			if err := row.Scan(&p); err != nil {
				scanErr = err
				continue
			}
			out.ProductPage.Total++
			if strings.TrimSpace(p.Category) != "" {
				categories[p.Category]++
				if strings.TrimSpace(p.SubCategory) != "" {
					subcategories[models.CategorySelection{Category: p.Category, SubCategory: p.SubCategory}]++
				}
			}
			if _, ok := selected[p.ID]; ok {
				copy := p
				selected[p.ID] = &copy
			}
			sub := models.CategorySelection{Category: p.Category, SubCategory: p.SubCategory}
			if !excludeIDs[p.ID] && !excludeCategories[p.Category] && !excludeSubs[sub] && (f.Scope.All || includeIDs[p.ID] || includeCategories[p.Category] || includeSubs[sub]) {
				out.Preview.Total++
				if len(out.Preview.Products) < 6 {
					out.Preview.Products = append(out.Preview.Products, p)
				}
			}
			if needle != "" && !strings.Contains(lower.String(strings.Join([]string{p.Title, p.Category, p.SubCategory, p.VendorName}, " ")), needle) {
				continue
			}
			matched := out.ProductPage.FilteredTotal
			if matched%int64(f.PerPage) == 0 {
				last = last[:0]
			}
			last = append(last, p)
			if int(matched/int64(f.PerPage))+1 == f.Page {
				out.Products = append(out.Products, p)
			}
			out.ProductPage.FilteredTotal++
		}
		if scanErr != nil {
			return scanErr
		}
		out.ProductPage.LastPage = max(1, int((out.ProductPage.FilteredTotal+int64(f.PerPage)-1)/int64(f.PerPage)))
		out.ProductPage.Page = min(f.Page, out.ProductPage.LastPage)
		if len(out.Products) == 0 && out.ProductPage.FilteredTotal > 0 {
			out.Products = last
		}
		if len(out.Products) > 0 {
			out.ProductPage.From = (out.ProductPage.Page-1)*f.PerPage + 1
			out.ProductPage.To = out.ProductPage.From + len(out.Products) - 1
		}
		for category, count := range categories {
			out.Categories = append(out.Categories, PromotionCategoryFacet{Category: category, Count: count})
		}
		sort.Slice(out.Categories, func(i, j int) bool { return out.Categories[i].Category < out.Categories[j].Category })
		for sub, count := range subcategories {
			out.Subcategories = append(out.Subcategories, PromotionSubcategoryFacet{Category: sub.Category, SubCategory: sub.SubCategory, Count: count})
		}
		sort.Slice(out.Subcategories, func(i, j int) bool {
			a, b := out.Subcategories[i], out.Subcategories[j]
			if a.Category == b.Category {
				return a.SubCategory < b.SubCategory
			}
			return a.Category < b.Category
		})
		for _, id := range selectedIDs {
			status := "unavailable"
			if selected[id] != nil {
				status = "available"
			}
			out.SelectedProducts = append(out.SelectedProducts, PromotionSelectedTarget{ID: id, Status: status, Product: selected[id]})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
func targetSet(values []string) map[string]bool {
	out := map[string]bool{}
	for _, v := range values {
		out[v] = true
	}
	return out
}
