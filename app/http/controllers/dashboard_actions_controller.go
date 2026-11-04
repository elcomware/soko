package controllers

import (
	"elcomware/app/services"
	"github.com/goravel/framework/contracts/http"
)

// These endpoints preserve the operational workspace after an existing domain
// action. They accept no caller-selected return URL and introduce no new state
// transitions; the services remain responsible for ownership and status checks.
func (d *DashboardController) CancelOrder(ctx http.Context) http.Response {
	actor := d.getAuthenticatedUser(ctx)
	if !actor.CanAccessWorkspace() {
		return ctx.Response().String(http.StatusForbidden, "Workspace access is not permitted.")
	}
	// Supplier workspaces contain only owned lines, not the complete order.
	// Whole-order cancellation requires the full-order read scope here; buyers
	// retain their own full-order cancellation flow under /account.
	if !actor.Can("orders.read", "") {
		return ctx.Response().String(http.StatusForbidden, "Full order access is required to cancel an order from this workspace.")
	}
	if err := services.CancelOrder(ctx.Request().Route("id"), actor); err != nil {
		ctx.Request().Session().Flash("error", err.Error())
	} else {
		ctx.Request().Session().Flash("success", "Order cancelled. Reserved inventory has been released; no payment was collected.")
	}
	return ctx.Response().Redirect(http.StatusSeeOther, "/dashboard?section=orders")
}

func (d *DashboardController) AcceptQuote(ctx http.Context) http.Response {
	return d.resolveWorkspaceRfq(ctx, false)
}

func (d *DashboardController) CloseRfq(ctx http.Context) http.Response {
	return d.resolveWorkspaceRfq(ctx, true)
}

func (d *DashboardController) resolveWorkspaceRfq(ctx http.Context, closeOnly bool) http.Response {
	actor := d.getAuthenticatedUser(ctx)
	if !actor.CanAccessWorkspace() {
		return ctx.Response().String(http.StatusForbidden, "Workspace access is not permitted.")
	}
	var err error
	if closeOnly {
		err = services.ResolveRfq(ctx.Request().Route("id"), "", true, actor)
	} else {
		err = services.ResolveRfqWithQuoteVersion(ctx.Request().Route("id"), ctx.Request().Route("quoteId"), ctx.Request().Input("quote_version"), actor)
	}
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
	} else {
		ctx.Request().Session().Flash("success", "Sourcing request updated.")
	}
	return ctx.Response().Redirect(http.StatusSeeOther, "/dashboard?section=rfqs")
}
