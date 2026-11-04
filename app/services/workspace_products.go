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

// WorkspaceProductFilters is separate from public catalog filters. Operational
// queries must retain private resource scope before applying search or paging.
type WorkspaceProductFilters struct {
	SelectedProduct string `json:"-"`
	Query           string `json:"q"`
	Category        string `json:"category"`
	Sort            string `json:"sort"`
	Page            int    `json:"-"`
	PerPage         int    `json:"-"`
}
type WorkspaceProductPageMeta struct {
	Page          int                     `json:"page"`
	PerPage       int                     `json:"per_page"`
	Total         int64                   `json:"total"`
	FilteredTotal int64                   `json:"filtered_total"`
	LastPage      int                     `json:"last_page"`
	From          int                     `json:"from"`
	To            int                     `json:"to"`
	Filters       WorkspaceProductFilters `json:"filters"`
}
type WorkspaceProductPage struct {
	ProductCategories     []string                              `json:"productCategories"`
	SelectedProduct       *models.Product                       `json:"selectedProduct"`
	SelectedProductStatus string                                `json:"selectedProductStatus"`
	Products              []models.Product                      `json:"products"`
	ProductPage           WorkspaceProductPageMeta              `json:"productPage"`
	RecordCapabilities    map[string]map[string]map[string]bool `json:"recordCapabilities"`
}

func ParseWorkspaceProductFilters(input map[string]string) (WorkspaceProductFilters, error) {
	f := WorkspaceProductFilters{Query: strings.TrimSpace(input["q"]), Category: strings.TrimSpace(input["category"]), Sort: input["sort"], SelectedProduct: strings.TrimSpace(input["selected_product"]), Page: 1, PerPage: 25}
	if f.Sort == "" {
		f.Sort = "newest"
	}
	for key, target := range map[string]*int{"page": &f.Page, "per_page": &f.PerPage} {
		if raw := input[key]; raw != "" {
			value, err := strconv.Atoi(raw)
			if err != nil {
				return f, errors.New("invalid " + key)
			}
			*target = value
		}
	}
	return f, validateWorkspaceProductFilters(f)
}
func validateWorkspaceProductFilters(f WorkspaceProductFilters) error {
	if f.Page < 1 || f.Page > 1000000 || f.PerPage < 1 || f.PerPage > 100 {
		return errors.New("page must be 1–1000000 and per_page must be 1–100")
	}
	if !utf8.ValidString(f.Query) || !utf8.ValidString(f.Category) || utf8.RuneCountInString(f.Query) > 120 || utf8.RuneCountInString(f.Category) > 128 || strings.ContainsAny(f.Query+f.Category, "\x00\r\n") {
		return errors.New("invalid product search or category")
	}
	if !utf8.ValidString(f.SelectedProduct) || utf8.RuneCountInString(f.SelectedProduct) > 64 || strings.ContainsAny(f.SelectedProduct, "\x00\r\n") {
		return errors.New("invalid selected product")
	}
	switch f.Sort {
	case "newest", "oldest", "title_asc", "title_desc":
	default:
		return errors.New("invalid product sort")
	}
	return nil
}

// LoadWorkspaceProducts is an opt-in dashboard query path. The existing account
// loader remains unchanged; only the dashboard products section consumes it.
// All query failures return nil instead of apparently empty successful results.
func LoadWorkspaceProducts(actor *models.User, f WorkspaceProductFilters) (*WorkspaceProductPage, error) {
	if actor == nil || !actor.CanAccessWorkspace() || !actor.HasPermission("products.read") {
		return nil, errors.New("product workspace access denied")
	}
	if err := validateWorkspaceProductFilters(f); err != nil {
		return nil, err
	}
	out := &WorkspaceProductPage{ProductCategories: []string{}, SelectedProductStatus: "none", Products: []models.Product{}, ProductPage: WorkspaceProductPageMeta{Page: f.Page, PerPage: f.PerPage, Filters: f}, RecordCapabilities: map[string]map[string]map[string]bool{"products": {}}}
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		scoped := func() orm.Query {
			q := tx.Model(&models.Product{})
			if !actor.Can("products.read", "") {
				q = q.Where("vendor_id = ?", actor.ID)
			}
			return q
		}
		filtered := func() orm.Query {
			q := scoped()
			if f.Category != "" {
				q = q.Where("category = ?", f.Category)
			}
			return q
		}
		var err error
		out.ProductPage.Total, err = scoped().Count()
		if err != nil {
			return err
		}
		var matchedIDs []string
		if f.Query == "" {
			out.ProductPage.FilteredTotal, err = filtered().Count()
		} else {
			matchedIDs, out.ProductPage.FilteredTotal, err = scanWorkspaceProductMatches(orderedWorkspaceProducts(filtered(), f.Sort), f)
		}
		if err != nil {
			return err
		}
		out.ProductPage.LastPage = max(1, int((out.ProductPage.FilteredTotal+int64(f.PerPage)-1)/int64(f.PerPage)))
		out.ProductPage.Page = min(f.Page, out.ProductPage.LastPage)
		out.ProductPage.Filters.Page = out.ProductPage.Page
		var categories []struct{ Category string }
		if err := scoped().Select("category").Distinct("category").Where("category <> ?", "").OrderBy("category").Get(&categories); err != nil {
			return err
		}
		for _, category := range categories {
			out.ProductCategories = append(out.ProductCategories, category.Category)
		}
		q := orderedWorkspaceProducts(filtered(), f.Sort)
		offset := (out.ProductPage.Page - 1) * f.PerPage
		if f.Query == "" {
			err = q.Offset(offset).Limit(f.PerPage).Get(&out.Products)
		} else if len(matchedIDs) > 0 {
			err = q.Where("id IN (?)", matchedIDs).Limit(f.PerPage).Get(&out.Products)
		}
		if err != nil {
			return err
		}
		if len(out.Products) > 0 {
			out.ProductPage.From = offset + 1
			out.ProductPage.To = offset + len(out.Products)
		}
		records := append([]models.Product{}, out.Products...)
		if f.SelectedProduct != "" {
			out.SelectedProductStatus = "unavailable"
			var selected models.Product
			if err := scoped().Where("id = ?", f.SelectedProduct).First(&selected); err != nil {
				return err
			}
			if selected.ID != "" {
				selected.Hydrate()
				selected.VendorRegion = models.NormalizeRegion(selected.VendorRegion)
				out.SelectedProduct = &selected
				out.SelectedProductStatus = "available"
				records = append(records, selected)
			}
		}
		ids := make([]string, 0, len(records))
		for i := range out.Products {
			out.Products[i].Hydrate()
			out.Products[i].VendorRegion = models.NormalizeRegion(out.Products[i].VendorRegion)

		}
		for _, p := range records {
			ids = append(ids, p.ID)
		}
		if len(ids) == 0 {
			return nil
		}
		// One bounded history query, never one order/offer lookup per product.
		var referenced []models.Product
		if err := tx.Model(&models.Product{}).Select("id").Where("id IN (?)", ids).Where("(EXISTS (SELECT 1 FROM order_items WHERE product_id = products.id) OR EXISTS (SELECT 1 FROM offers WHERE product_id = products.id))").Get(&referenced); err != nil {
			return err
		}
		history := map[string]bool{}
		for _, p := range referenced {
			history[p.ID] = true
		}
		for _, p := range records {
			update := actor.CanManageProduct(&p)
			out.RecordCapabilities["products"][p.ID] = map[string]bool{"update": update, "delete": actor.Can("products.delete", p.VendorID) && !history[p.ID], "media": update && actor.Can("media.create", actor.ID)}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func orderedWorkspaceProducts(q orm.Query, sort string) orm.Query {
	switch sort {
	case "oldest":
		return q.OrderBy("created_at").OrderBy("id")
	case "title_asc":
		return q.OrderBy("title").OrderBy("id")
	case "title_desc":
		return q.OrderByDesc("title").OrderByDesc("id")
	default:
		return q.OrderByDesc("created_at").OrderByDesc("id")
	}
}

// SQL LOWER is not Unicode-aware on the supported SQLite database. A narrow
// cursor preserves multilingual substring semantics without loading full product
// records. Search is O(N) over authorized category-filtered rows; retained memory
// is O(per_page), and every cursor error aborts rather than returning false totals.
// This runs inside the same read transaction as counts and page hydration.
func scanWorkspaceProductMatches(q orm.Query, f WorkspaceProductFilters) ([]string, int64, error) {
	lower := cases.Lower(language.Und)
	needle := lower.String(f.Query)
	requested := []string{}
	last := []string{}
	var matched int64
	var scanErr error
	cursor := q.Select("id", "title", "sku", "category", "vendor_name").Cursor()
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
		var product models.Product
		if err := row.Scan(&product); err != nil {
			scanErr = err
			continue
		}
		found := false
		for _, value := range []string{product.Title, product.Sku, product.Category, product.VendorName} {
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
		last = append(last, product.ID)
		page := int(matched/int64(f.PerPage)) + 1
		if page == f.Page {
			requested = append(requested, product.ID)
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
