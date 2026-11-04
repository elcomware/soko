package middleware

import "elcomware/app/models"

type WorkspaceLink struct {
	Label string `json:"label"`
	Href  string `json:"href"`
	Icon  string `json:"icon"`
	Group string `json:"group"`
}

func WorkspaceNavigation(user *models.User) []WorkspaceLink {
	links := []WorkspaceLink{}
	if !user.CanAccessWorkspace() {
		return links
	}
	add := func(allowed bool, label, href, icon, group string) {
		if allowed {
			links = append(links, WorkspaceLink{label, href, icon, group})
		}
	}
	add(true, "Overview", "/dashboard", "grid", "Workspace")
	add(user.HasPermission("products.read"), "Products", "/dashboard?section=products", "box", "Commerce")
	add(user.HasPermission("orders.read"), "Orders", "/dashboard?section=orders", "receipt", "Commerce")
	add(user.HasPermission("offers.read"), "Offers", "/dashboard?section=offers", "tag", "Commerce")
	add(user.HasPermission("rfqs.read") || user.HasPermission("rfqs.quote"), "Sourcing requests", "/dashboard?section=rfqs", "document", "Commerce")
	add(user.HasPermission("promotions.manage"), "Promotions", "/dashboard/promotions", "tag", "Commerce")
	add(user.HasPermission("media.create") || user.Can("media.settings.manage", ""), "Media library", "/dashboard/media-library", "image", "Storefront")
	add(user.Can("content.update", ""), "Storefront editor", "/dashboard/merchandising", "layout", "Storefront")
	add(user.Can("brand.update", ""), "Brand & presentation", "/dashboard/settings?panel=brand", "palette", "Storefront")
	add(user.Can("users.read", ""), "Accounts", "/dashboard?section=users", "users", "Administration")
	add(user.Can("permissions.manage", ""), "Access & permissions", "/dashboard/settings?panel=access", "shield", "Administration")
	add(user.HasPermission("payments.configure"), "Payment configuration", "/dashboard/settings?panel=payments", "card", "Administration")
	return links
}
