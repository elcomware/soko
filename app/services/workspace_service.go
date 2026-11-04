package services

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/queries"
	"errors"
)

type WorkspaceData struct {
	Products []models.Product
	Orders   []models.Order
	Offers   []models.Offer
	Rfqs     []models.Rfq
	Users    []models.User
	Stats    map[string]any
}

func LoadWorkspace(u *models.User, account bool) (*WorkspaceData, error) {
	if u == nil || u.ID == "" {
		return nil, errors.New("sign in required")
	}
	d := &WorkspaceData{Products: []models.Product{}, Orders: []models.Order{}, Offers: []models.Offer{}, Rfqs: []models.Rfq{}, Users: []models.User{}}
	var err error
	if account || u.Role == models.RoleBuyer {
		d.Orders, err = queries.NewOrderQueries().FindByUserID(u.ID)
		if err != nil {
			return nil, err
		}
		if err = facades.Orm().Query().Where("buyer_id = ?", u.ID).OrderByDesc("created_at").Get(&d.Offers); err != nil {
			return nil, err
		}
		if err = facades.Orm().Query().Where("buyer_id = ?", u.ID).OrderByDesc("created_at").Get(&d.Rfqs); err != nil {
			return nil, err
		}
	} else if u.Role == models.RoleAdmin {
		if err = facades.Orm().Query().Get(&d.Products); err != nil {
			return nil, err
		}
		d.Orders, err = queries.NewOrderQueries().FindAll()
		if err != nil {
			return nil, err
		}
		if err = facades.Orm().Query().OrderByDesc("created_at").Get(&d.Offers); err != nil {
			return nil, err
		}
		if err = facades.Orm().Query().OrderByDesc("created_at").Get(&d.Rfqs); err != nil {
			return nil, err
		}
		if err = facades.Orm().Query().Get(&d.Users); err != nil {
			return nil, err
		}
	} else if u.Role == models.RoleMerchant {
		if err = facades.Orm().Query().Where("vendor_id = ?", u.ID).Get(&d.Products); err != nil {
			return nil, err
		}
		if err = facades.Orm().Query().Where("product_id IN (SELECT id FROM products WHERE vendor_id = ?)", u.ID).OrderByDesc("created_at").Get(&d.Offers); err != nil {
			return nil, err
		}
		if err = facades.Orm().Query().Where("status = ? OR id IN (SELECT rfq_id FROM rfq_quotes WHERE vendor_id = ?)", models.RfqStatusOpen, u.ID).OrderByDesc("created_at").Get(&d.Rfqs); err != nil {
			return nil, err
		}
		if err = facades.Orm().Query().Where("id IN (SELECT order_id FROM order_items WHERE product_id IN (SELECT id FROM products WHERE vendor_id = ?))", u.ID).OrderByDesc("created_at").Get(&d.Orders); err != nil {
			return nil, err
		}
		for i := range d.Orders {
			var items []models.OrderItem
			if err = facades.Orm().Query().Where("order_id = ? AND product_id IN (SELECT id FROM products WHERE vendor_id = ?)", d.Orders[i].ID, u.ID).Get(&items); err != nil {
				return nil, err
			}
			d.Orders[i].Items = items
			d.Orders[i].TotalAmount = 0
			for j := range items {
				items[j].DigitalLicenseKey = ""
				items[j].DigitalDownloadUrl = ""
				d.Orders[i].TotalAmount += items[j].TotalPrice
			}
		}
	} else {
		return nil, errors.New("account role is not supported")
	}
	for i := range d.Products {
		d.Products[i].Hydrate()
		d.Products[i].VendorRegion = models.NormalizeRegion(d.Products[i].VendorRegion)
	}
	for i := range d.Rfqs {
		q := facades.Orm().Query().Where("rfq_id = ?", d.Rfqs[i].ID)
		if !account && u.Role == models.RoleMerchant {
			q = q.Where("vendor_id = ?", u.ID)
		}
		d.Rfqs[i].Quotes = []models.RfqQuote{}
		if err = q.OrderByDesc("created_at").Get(&d.Rfqs[i].Quotes); err != nil {
			return nil, err
		}
	}
	var revenue float64
	var pending, activeOffers, activeOrders int
	for _, o := range d.Orders {
		if o.Status == models.OrderStatusCompleted || o.Status == models.OrderStatusShipped {
			revenue += o.TotalAmount
		}
		if o.Status == models.OrderStatusPending {
			activeOrders++
		}
	}
	for _, r := range d.Rfqs {
		if r.Status == models.RfqStatusOpen {
			pending++
		}
	}
	for _, o := range d.Offers {
		if o.Status == models.OfferStatusPending {
			activeOffers++
		}
	}
	d.Stats = map[string]any{"total_revenue": revenue, "active_orders": activeOrders, "total_products": len(d.Products), "pending_rfqs": pending, "active_negotiations": activeOffers, "trade_volume_covered": 0}
	return d, nil
}
