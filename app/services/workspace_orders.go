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

type WorkspaceOrderFilters struct {
	Query         string `json:"q"`
	Status        string `json:"status"`
	Sort          string `json:"sort"`
	Page          int    `json:"-"`
	PerPage       int    `json:"-"`
	SelectedOrder string `json:"-"`
}
type WorkspaceOrderPageMeta struct {
	Page          int                   `json:"page"`
	PerPage       int                   `json:"per_page"`
	Total         int64                 `json:"total"`
	FilteredTotal int64                 `json:"filtered_total"`
	LastPage      int                   `json:"last_page"`
	From          int                   `json:"from"`
	To            int                   `json:"to"`
	Filters       WorkspaceOrderFilters `json:"filters"`
}
type WorkspaceOrderPage struct {
	Orders              []models.Order                        `json:"orders"`
	OrderPage           WorkspaceOrderPageMeta                `json:"orderPage"`
	SelectedOrder       *models.Order                         `json:"selectedOrder"`
	SelectedOrderStatus string                                `json:"selectedOrderStatus"`
	OrderVisibility     string                                `json:"orderVisibility"`
	RecordCapabilities  map[string]map[string]map[string]bool `json:"recordCapabilities"`
}

func ParseWorkspaceOrderFilters(input map[string]string) (WorkspaceOrderFilters, error) {
	f := WorkspaceOrderFilters{Query: strings.TrimSpace(input["q"]), Status: input["status"], Sort: input["sort"], Page: 1, PerPage: 25, SelectedOrder: strings.TrimSpace(input["selected_order"])}
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
	return f, validateWorkspaceOrderFilters(f)
}
func validateWorkspaceOrderFilters(f WorkspaceOrderFilters) error {
	if f.Page < 1 || f.Page > 1000000 || f.PerPage < 1 || f.PerPage > 100 {
		return errors.New("invalid pagination")
	}
	if f.Sort != "newest" && f.Sort != "oldest" {
		return errors.New("invalid order sort")
	}
	switch f.Status {
	case "", "pending", "cancelled", "completed", "shipped":
	default:
		return errors.New("invalid order status")
	}
	if !utf8.ValidString(f.Query) || utf8.RuneCountInString(f.Query) > 120 || !utf8.ValidString(f.SelectedOrder) || utf8.RuneCountInString(f.SelectedOrder) > 64 || strings.ContainsAny(f.Query+f.SelectedOrder, "\x00\r\n") {
		return errors.New("invalid order search or selection")
	}
	return nil
}

// LoadWorkspaceOrders keeps the account loader unchanged. All counts and details
// share one read snapshot; supplier views never hydrate other vendors' items.
func LoadWorkspaceOrders(actor *models.User, f WorkspaceOrderFilters) (*WorkspaceOrderPage, error) {
	if actor == nil || !actor.CanAccessWorkspace() || !actor.HasPermission("orders.read") {
		return nil, errors.New("order workspace access denied")
	}
	if err := validateWorkspaceOrderFilters(f); err != nil {
		return nil, err
	}
	full := actor.Can("orders.read", "")
	out := &WorkspaceOrderPage{Orders: []models.Order{}, SelectedOrderStatus: "none", OrderVisibility: "supplier", OrderPage: WorkspaceOrderPageMeta{Page: f.Page, PerPage: f.PerPage, Filters: f}, RecordCapabilities: map[string]map[string]map[string]bool{"orders": {}}}
	if full {
		out.OrderVisibility = "full"
	}
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		scoped := func() orm.Query {
			q := tx.Model(&models.Order{})
			if !full {
				q = q.Where("id IN (SELECT order_id FROM order_items WHERE product_id IN (SELECT id FROM products WHERE vendor_id = ?))", actor.ID)
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
		out.OrderPage.Total, err = scoped().Count()
		if err != nil {
			return err
		}
		var ids []string
		if f.Query == "" {
			out.OrderPage.FilteredTotal, err = filtered().Count()
		} else {
			ids, out.OrderPage.FilteredTotal, err = scanWorkspaceOrderMatches(ordered(filtered()), f)
		}
		if err != nil {
			return err
		}
		out.OrderPage.LastPage = max(1, int((out.OrderPage.FilteredTotal+int64(f.PerPage)-1)/int64(f.PerPage)))
		out.OrderPage.Page = min(f.Page, out.OrderPage.LastPage)
		out.OrderPage.Filters.Page = out.OrderPage.Page
		offset := (out.OrderPage.Page - 1) * f.PerPage
		if f.Query == "" {
			err = ordered(filtered()).Offset(offset).Limit(f.PerPage).Get(&out.Orders)
		} else if len(ids) > 0 {
			err = ordered(filtered()).Where("id IN (?)", ids).Limit(f.PerPage).Get(&out.Orders)
		}
		if err != nil {
			return err
		}
		if len(out.Orders) > 0 {
			out.OrderPage.From = offset + 1
			out.OrderPage.To = offset + len(out.Orders)
		}
		if f.SelectedOrder != "" {
			out.SelectedOrderStatus = "unavailable"
			var selected models.Order
			if err := scoped().Where("id = ?", f.SelectedOrder).First(&selected); err != nil {
				return err
			}
			if selected.ID != "" {
				out.SelectedOrder = &selected
				out.SelectedOrderStatus = "available"
			}
		}
		records := make([]*models.Order, 0, len(out.Orders)+1)
		ids = []string{}
		for i := range out.Orders {
			records = append(records, &out.Orders[i])
			ids = append(ids, out.Orders[i].ID)
		}
		if out.SelectedOrder != nil {
			records = append(records, out.SelectedOrder)
			ids = append(ids, out.SelectedOrder.ID)
		}
		if len(ids) == 0 {
			return nil
		}
		itemsQuery := tx.Model(&models.OrderItem{}).Where("order_id IN (?)", ids)
		if !full {
			itemsQuery = itemsQuery.Where("product_id IN (SELECT id FROM products WHERE vendor_id = ?)", actor.ID)
		}
		var items []models.OrderItem
		if err := itemsQuery.OrderBy("created_at").OrderBy("id").Get(&items); err != nil {
			return err
		}
		byOrder := map[string][]models.OrderItem{}
		for _, item := range items {
			if !full {
				item.DigitalLicenseKey = ""
				item.DigitalDownloadUrl = ""
			}
			byOrder[item.OrderID] = append(byOrder[item.OrderID], item)
		}
		for _, order := range records {
			order.Items = byOrder[order.ID]
			if order.Items == nil {
				order.Items = []models.OrderItem{}
			}
			if !full {
				order.TotalAmount = 0
				order.HasDigital = false
				order.HasPhysical = false
				for _, item := range order.Items {
					order.TotalAmount += item.TotalPrice
					if item.ProductType != models.TypePhysical {
						order.HasDigital = true
					} else {
						order.HasPhysical = true
					}
				}
			}
			out.RecordCapabilities["orders"][order.ID] = map[string]bool{"cancel": full && order.Status == models.OrderStatusPending && actor.Can("orders.cancel", order.UserID)}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
func scanWorkspaceOrderMatches(q orm.Query, f WorkspaceOrderFilters) ([]string, int64, error) {
	lower := cases.Lower(language.Und)
	needle := lower.String(f.Query)
	requested := []string{}
	last := []string{}
	var matched int64
	var scanErr error
	cursor := q.Select("id", "order_number", "customer_name", "customer_email", "status").Cursor()
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
		var order models.Order
		if err := row.Scan(&order); err != nil {
			scanErr = err
			continue
		}
		found := false
		for _, value := range []string{order.OrderNumber, order.CustomerName, order.CustomerEmail, string(order.Status)} {
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
		last = append(last, order.ID)
		page := int(matched/int64(f.PerPage)) + 1
		if page == f.Page {
			requested = append(requested, order.ID)
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
