package services

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"errors"
	"github.com/goravel/framework/contracts/database/orm"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type WorkspacePromotionFilters struct {
	Query            string `json:"q"`
	Status           string `json:"status"`
	Sort             string `json:"sort"`
	Page             int    `json:"-"`
	PerPage          int    `json:"-"`
	SelectedCampaign string `json:"-"`
}
type WorkspacePromotionPageMeta struct {
	Page          int                       `json:"page"`
	PerPage       int                       `json:"per_page"`
	Total         int64                     `json:"total"`
	FilteredTotal int64                     `json:"filtered_total"`
	LastPage      int                       `json:"last_page"`
	From          int                       `json:"from"`
	To            int                       `json:"to"`
	Filters       WorkspacePromotionFilters `json:"filters"`
}
type WorkspacePromotionPage struct {
	Products               []map[string]string        `json:"products"`
	CampaignCounts         map[string]int64           `json:"campaignCounts"`
	ServerTime             string                     `json:"serverTime"`
	NextStatusChange       string                     `json:"nextStatusChange"`
	Campaigns              []models.Promotion         `json:"campaigns"`
	CampaignPage           WorkspacePromotionPageMeta `json:"campaignPage"`
	SelectedCampaign       *models.Promotion          `json:"selectedCampaign"`
	SelectedCampaignStatus string                     `json:"selectedCampaignStatus"`
}

func ParseWorkspacePromotionFilters(input map[string]string) (WorkspacePromotionFilters, error) {
	f := WorkspacePromotionFilters{Query: strings.TrimSpace(input["q"]), Status: input["status"], Sort: input["sort"], Page: 1, PerPage: 25, SelectedCampaign: strings.TrimSpace(input["selected_campaign"])}
	if f.Sort == "" {
		f.Sort = "newest"
	}
	for key, target := range map[string]*int{"page": &f.Page, "per_page": &f.PerPage} {
		if raw := input[key]; raw != "" {
			v, err := strconv.Atoi(raw)
			if err != nil {
				return f, errors.New("invalid " + key)
			}
			*target = v
		}
	}
	return f, validateWorkspacePromotionFilters(f)
}
func validateWorkspacePromotionFilters(f WorkspacePromotionFilters) error {
	if f.Page < 1 || f.Page > 1000000 || f.PerPage < 1 || f.PerPage > 100 {
		return errors.New("invalid pagination")
	}
	if f.Sort != "newest" && f.Sort != "oldest" {
		return errors.New("invalid promotion sort")
	}
	switch f.Status {
	case "", "active", "scheduled", "ended", "disabled":
	default:
		return errors.New("invalid promotion status")
	}
	if !utf8.ValidString(f.Query) || utf8.RuneCountInString(f.Query) > 120 || !utf8.ValidString(f.SelectedCampaign) || utf8.RuneCountInString(f.SelectedCampaign) > 64 || strings.ContainsAny(f.Query+f.SelectedCampaign, "\x00\r\n") {
		return errors.New("invalid promotion search or selection")
	}
	return nil
}

// LoadWorkspacePromotions keeps the account loader unchanged. All counts and details
// share one read snapshot; current product ownership scopes supplier access.
func LoadWorkspacePromotions(actor *models.User, f WorkspacePromotionFilters, now time.Time) (*WorkspacePromotionPage, error) {
	if actor == nil || !actor.CanAccessWorkspace() || !actor.HasPermission("promotions.manage") {
		return nil, errors.New("promotion workspace access denied")
	}
	if err := validateWorkspacePromotionFilters(f); err != nil {
		return nil, err
	}
	full := actor.Can("promotions.manage", "")
	out := &WorkspacePromotionPage{Campaigns: []models.Promotion{}, SelectedCampaignStatus: "none", CampaignPage: WorkspacePromotionPageMeta{Page: f.Page, PerPage: f.PerPage, Filters: f}}
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		scoped := func() orm.Query {
			q := tx.Model(&models.Promotion{})
			if !full {
				q = q.Where("vendor_id = ?", actor.ID)
			}
			return q
		}
		// Target choices remain complete until a separate bounded picker contract exists.
		productsQuery := tx.Model(&models.Product{}).Select("id", "title", "category", "sub_category", "vendor_id", "vendor_name")
		if !full {
			productsQuery = productsQuery.Where("vendor_id = ?", actor.ID)
		}
		var products []models.Product
		if err := productsQuery.OrderBy("id").Get(&products); err != nil {
			return err
		}
		merchantNames := map[string]string{}
		out.Products = []map[string]string{}
		for _, p := range products {
			name := p.VendorName
			if name == "" {
				name = p.VendorID
			}
			merchantNames[p.VendorID] = name
			out.Products = append(out.Products, map[string]string{"id": p.ID, "title": p.Title, "category": p.Category, "sub_category": p.SubCategory, "vendor_id": p.VendorID, "vendor_name": p.VendorName})
		}
		out.ServerTime = now.UTC().Format(time.RFC3339Nano)
		out.CampaignCounts = map[string]int64{}
		for _, status := range []string{"active", "scheduled", "ended", "disabled"} {
			count, err := promotionStatusQuery(scoped(), status, now).Count()
			if err != nil {
				return err
			}
			out.CampaignCounts[status] = count
		}
		for _, field := range []string{"starts_at", "ends_at"} {
			var next models.Promotion
			if err := scoped().Where("enabled = ?", true).Where(field+" > ?", now).OrderBy(field).First(&next); err != nil {
				return err
			}
			if next.ID != "" {
				boundary := next.StartsAt
				if field == "ends_at" {
					boundary = next.EndsAt
				}
				if out.NextStatusChange == "" {
					out.NextStatusChange = boundary.UTC().Format(time.RFC3339Nano)
				} else {
					previous, _ := time.Parse(time.RFC3339Nano, out.NextStatusChange)
					if boundary.Before(previous) {
						out.NextStatusChange = boundary.UTC().Format(time.RFC3339Nano)
					}
				}
			}
		}
		filtered := func() orm.Query {
			q := scoped()
			if f.Status != "" {
				q = promotionStatusQuery(q, f.Status, now)
			}
			return q
		}
		ordered := func(q orm.Query) orm.Query {
			if f.Sort == "oldest" {
				return q.OrderBy("created_at").OrderBy("id")
			}
			return q.OrderByDesc("created_at").OrderByDesc("id")
		}
		var err error
		out.CampaignPage.Total, err = scoped().Count()
		out.CampaignCounts["total"] = out.CampaignPage.Total
		if err != nil {
			return err
		}
		var ids []string
		if f.Query == "" {
			out.CampaignPage.FilteredTotal, err = filtered().Count()
		} else {
			ids, out.CampaignPage.FilteredTotal, err = scanWorkspacePromotionMatches(ordered(filtered()), f, merchantNames)
		}
		if err != nil {
			return err
		}
		out.CampaignPage.LastPage = max(1, int((out.CampaignPage.FilteredTotal+int64(f.PerPage)-1)/int64(f.PerPage)))
		out.CampaignPage.Page = min(f.Page, out.CampaignPage.LastPage)
		out.CampaignPage.Filters.Page = out.CampaignPage.Page
		offset := (out.CampaignPage.Page - 1) * f.PerPage
		if f.Query == "" {
			err = ordered(filtered()).Offset(offset).Limit(f.PerPage).Get(&out.Campaigns)
		} else if len(ids) > 0 {
			err = ordered(filtered()).Where("id IN (?)", ids).Limit(f.PerPage).Get(&out.Campaigns)
		}
		if err != nil {
			return err
		}
		if len(out.Campaigns) > 0 {
			out.CampaignPage.From = offset + 1
			out.CampaignPage.To = offset + len(out.Campaigns)
		}
		if f.SelectedCampaign != "" {
			out.SelectedCampaignStatus = "unavailable"
			var selected models.Promotion
			if err := scoped().Where("id = ?", f.SelectedCampaign).First(&selected); err != nil {
				return err
			}
			if selected.ID != "" {
				out.SelectedCampaign = &selected
				out.SelectedCampaignStatus = "available"
			}
		}
		for i := range out.Campaigns {
			out.Campaigns[i].Hydrate()
		}
		if out.SelectedCampaign != nil {
			out.SelectedCampaign.Hydrate()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
func scanWorkspacePromotionMatches(q orm.Query, f WorkspacePromotionFilters, merchantNames map[string]string) ([]string, int64, error) {
	lower := cases.Lower(language.Und)
	needle := lower.String(f.Query)
	requested := []string{}
	last := []string{}
	var matched int64
	var scanErr error
	cursor := q.Select("id", "name", "vendor_id", "kind").Cursor()
	for row := range cursor {
		// Always drain: the ORM cursor's producer cannot be left blocked on send.
		if row.Err() != nil {
			if scanErr == nil {
				scanErr = row.Err()
			}
			continue
		}
		if scanErr != nil {
			continue
		}
		var promotion models.Promotion
		if err := row.Scan(&promotion); err != nil {
			scanErr = err
			continue
		}
		found := false
		for _, value := range []string{strings.Join([]string{promotion.Name, promotionMerchantName(promotion.VendorID, merchantNames), promotion.Kind}, " ")} {
			if strings.Contains(lower.String(value), needle) {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		if matched%int64(f.PerPage) == 0 {
			last = last[:0]
		}
		last = append(last, promotion.ID)
		page := int(matched/int64(f.PerPage)) + 1
		if page == f.Page {
			requested = append(requested, promotion.ID)
		}
		matched++
	}
	if scanErr != nil {
		return nil, 0, scanErr
	}
	if len(requested) == 0 && matched > 0 {
		return last, matched, nil
	}
	return requested, matched, nil
}

func promotionMerchantName(id string, names map[string]string) string {
	if name, ok := names[id]; ok {
		return name
	}
	return id
}
func promotionStatusQuery(q orm.Query, status string, now time.Time) orm.Query {
	switch status {
	case "disabled":
		return q.Where("enabled = ?", false)
	case "scheduled":
		return q.Where("enabled = ?", true).Where("starts_at > ?", now)
	case "ended":
		return q.Where("enabled = ?", true).Where("starts_at <= ?", now).Where("ends_at <= ?", now)
	case "active":
		return q.Where("enabled = ?", true).Where("starts_at <= ?", now).Where("ends_at > ?", now)
	}
	return q
}
