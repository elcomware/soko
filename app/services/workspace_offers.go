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
	"unicode/utf8"
)

type WorkspaceOfferFilters struct {
	Query         string `json:"q"`
	Status        string `json:"status"`
	Sort          string `json:"sort"`
	Page          int    `json:"-"`
	PerPage       int    `json:"-"`
	SelectedOffer string `json:"-"`
}
type WorkspaceOfferPageMeta struct {
	Page          int                   `json:"page"`
	PerPage       int                   `json:"per_page"`
	Total         int64                 `json:"total"`
	FilteredTotal int64                 `json:"filtered_total"`
	LastPage      int                   `json:"last_page"`
	From          int                   `json:"from"`
	To            int                   `json:"to"`
	Filters       WorkspaceOfferFilters `json:"filters"`
}
type WorkspaceOfferPage struct {
	Offers              []models.Offer                        `json:"offers"`
	OfferPage           WorkspaceOfferPageMeta                `json:"offerPage"`
	SelectedOffer       *models.Offer                         `json:"selectedOffer"`
	SelectedOfferStatus string                                `json:"selectedOfferStatus"`
	OfferVisibility     string                                `json:"offerVisibility"`
	RecordCapabilities  map[string]map[string]map[string]bool `json:"recordCapabilities"`
}

func ParseWorkspaceOfferFilters(input map[string]string) (WorkspaceOfferFilters, error) {
	f := WorkspaceOfferFilters{Query: strings.TrimSpace(input["q"]), Status: input["status"], Sort: input["sort"], Page: 1, PerPage: 25, SelectedOffer: strings.TrimSpace(input["selected_offer"])}
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
	return f, validateWorkspaceOfferFilters(f)
}
func validateWorkspaceOfferFilters(f WorkspaceOfferFilters) error {
	if f.Page < 1 || f.Page > 1000000 || f.PerPage < 1 || f.PerPage > 100 {
		return errors.New("invalid pagination")
	}
	if f.Sort != "newest" && f.Sort != "oldest" {
		return errors.New("invalid offer sort")
	}
	switch f.Status {
	case "", "pending", "accepted", "declined", "countered":
	default:
		return errors.New("invalid offer status")
	}
	if !utf8.ValidString(f.Query) || utf8.RuneCountInString(f.Query) > 120 || !utf8.ValidString(f.SelectedOffer) || utf8.RuneCountInString(f.SelectedOffer) > 64 || strings.ContainsAny(f.Query+f.SelectedOffer, "\x00\r\n") {
		return errors.New("invalid offer search or selection")
	}
	return nil
}

// LoadWorkspaceOffers keeps the account loader unchanged. All counts and details
// share one read snapshot; current product ownership scopes supplier access.
func LoadWorkspaceOffers(actor *models.User, f WorkspaceOfferFilters) (*WorkspaceOfferPage, error) {
	if actor == nil || !actor.CanAccessWorkspace() || !actor.HasPermission("offers.read") {
		return nil, errors.New("offer workspace access denied")
	}
	if err := validateWorkspaceOfferFilters(f); err != nil {
		return nil, err
	}
	full := actor.Can("offers.read", "")
	out := &WorkspaceOfferPage{Offers: []models.Offer{}, SelectedOfferStatus: "none", OfferVisibility: "supplier", OfferPage: WorkspaceOfferPageMeta{Page: f.Page, PerPage: f.PerPage, Filters: f}, RecordCapabilities: map[string]map[string]map[string]bool{"offers": {}}}
	if full {
		out.OfferVisibility = "full"
	}
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		scoped := func() orm.Query {
			q := tx.Model(&models.Offer{})
			if !full {
				q = q.Where("product_id IN (SELECT id FROM products WHERE vendor_id = ?)", actor.ID)
			}
			return q
		}
		filtered := func() orm.Query {
			q := scoped()
			if f.Status != "" {
				q = q.Where("status = ?", f.Status)
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
		out.OfferPage.Total, err = scoped().Count()
		if err != nil {
			return err
		}
		var ids []string
		if f.Query == "" {
			out.OfferPage.FilteredTotal, err = filtered().Count()
		} else {
			ids, out.OfferPage.FilteredTotal, err = scanWorkspaceOfferMatches(ordered(filtered()), f)
		}
		if err != nil {
			return err
		}
		out.OfferPage.LastPage = max(1, int((out.OfferPage.FilteredTotal+int64(f.PerPage)-1)/int64(f.PerPage)))
		out.OfferPage.Page = min(f.Page, out.OfferPage.LastPage)
		out.OfferPage.Filters.Page = out.OfferPage.Page
		offset := (out.OfferPage.Page - 1) * f.PerPage
		if f.Query == "" {
			err = ordered(filtered()).Offset(offset).Limit(f.PerPage).Get(&out.Offers)
		} else if len(ids) > 0 {
			err = ordered(filtered()).Where("id IN (?)", ids).Limit(f.PerPage).Get(&out.Offers)
		}
		if err != nil {
			return err
		}
		if len(out.Offers) > 0 {
			out.OfferPage.From = offset + 1
			out.OfferPage.To = offset + len(out.Offers)
		}
		if f.SelectedOffer != "" {
			out.SelectedOfferStatus = "unavailable"
			var selected models.Offer
			if err := scoped().Where("id = ?", f.SelectedOffer).First(&selected); err != nil {
				return err
			}
			if selected.ID != "" {
				out.SelectedOffer = &selected
				out.SelectedOfferStatus = "available"
			}
		}
		records := make([]*models.Offer, 0, len(out.Offers)+1)
		ids = []string{}
		for i := range out.Offers {
			records = append(records, &out.Offers[i])
			ids = append(ids, out.Offers[i].ID)
		}
		if out.SelectedOffer != nil {
			records = append(records, out.SelectedOffer)
			ids = append(ids, out.SelectedOffer.ID)
		}
		if len(ids) == 0 {
			return nil
		}
		productIDs := make([]string, 0, len(records))
		for _, offer := range records {
			productIDs = append(productIDs, offer.ProductID)
		}
		var products []models.Product
		if err := tx.Model(&models.Product{}).Select("id", "vendor_id").Where("id IN (?)", productIDs).Get(&products); err != nil {
			return err
		}
		owners := map[string]string{}
		for _, product := range products {
			owners[product.ID] = product.VendorID
		}
		for _, offer := range records {
			owner, exists := owners[offer.ProductID]
			out.RecordCapabilities["offers"][offer.ID] = map[string]bool{"decide": exists && offer.Status == models.OfferStatusPending && actor.Can("offers.decide", owner)}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
func scanWorkspaceOfferMatches(q orm.Query, f WorkspaceOfferFilters) ([]string, int64, error) {
	lower := cases.Lower(language.Und)
	needle := lower.String(f.Query)
	requested := []string{}
	last := []string{}
	var matched int64
	var scanErr error
	cursor := q.Select("id", "product_title", "buyer_company", "buyer_name", "buyer_email", "status").Cursor()
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
		var offer models.Offer
		if err := row.Scan(&offer); err != nil {
			scanErr = err
			continue
		}
		found := false
		for _, value := range []string{offer.ProductTitle, offer.BuyerCompany, offer.BuyerName, offer.BuyerEmail, string(offer.Status)} {
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
		last = append(last, offer.ID)
		page := int(matched/int64(f.PerPage)) + 1
		if page == f.Page {
			requested = append(requested, offer.ID)
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
