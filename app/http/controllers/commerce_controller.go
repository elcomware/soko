package controllers

import (
	"elcomware/app/actions"
	appmiddleware "elcomware/app/http/middleware"
	"elcomware/app/queries"
	"elcomware/app/services"
	"encoding/json"
	"errors"
	"github.com/goravel/framework/contracts/http"
	inertiafacades "github.com/goravel/inertia/facades"
	"math"
	"sort"
	"strconv"
)

type CommerceController struct{}

func NewCommerceController() *CommerceController { return &CommerceController{} }
func cartFrom(ctx http.Context) map[string]int {
	cart := map[string]int{}
	if ctx.Request().HasSession() {
		raw, _ := ctx.Request().Session().Get("cart").(string)
		_ = json.Unmarshal([]byte(raw), &cart)
	}
	if len(cart) > 100 {
		return map[string]int{}
	}
	return cart
}
func saveCart(ctx http.Context, cart map[string]int) {
	b, _ := json.Marshal(cart)
	ctx.Request().Session().Put("cart", string(b))
}
func CartSummary(ctx http.Context) map[string]any {
	cart := cartFrom(ctx)
	items := []map[string]any{}
	issues := []string{}
	var total float64
	ids := []string{}
	for id := range cart {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		qty := cart[id]
		p, err := queries.NewProductQueries().FindByID(id)
		if err != nil || p == nil || p.ID == "" {
			issues = append(issues, "A listing is no longer available. It has been removed from your cart.")
			delete(cart, id)
			continue
		}
		price, err := p.PriceFor(qty)
		if err != nil {
			issues = append(issues, p.Title+": "+err.Error())
			price = p.Price
		}
		line := math.Round(price*float64(qty)*100) / 100
		items = append(items, map[string]any{"product_id": id, "quantity": qty, "product": p.Public(), "unit_price": price, "total_price": line})
		total += line
	}
	saveCart(ctx, cart)
	total = math.Round(total*100) / 100
	return map[string]any{"items": items, "subtotal": total, "total_amount": total, "currency": "USD", "issues": issues}
}
func (c *CommerceController) Cart(ctx http.Context) http.Response {
	currency := services.NewCurrencyService()
	return inertiafacades.Inertia().Render(ctx, "Store/Cart", map[string]any{"cart": CartSummary(ctx), "currencies": currency.GetCurrencies(), "regions": currency.GetRegions()})
}
func (c *CommerceController) AddItem(ctx http.Context) http.Response    { return c.change(ctx, true) }
func (c *CommerceController) UpdateItem(ctx http.Context) http.Response { return c.change(ctx, false) }
func (c *CommerceController) change(ctx http.Context, add bool) http.Response {
	id := ctx.Request().Route("id")
	if add {
		id = ctx.Request().Input("product_id")
	}
	qty, err := strconv.Atoi(ctx.Request().Input("quantity"))
	cart := cartFrom(ctx)
	if err != nil || qty < 1 || qty > 1000000 {
		err = errors.New("quantity must be a positive whole number within its limit")
	} else {
		if add {
			qty += cart[id]
		}
		if len(cart) >= 100 && cart[id] == 0 {
			err = errors.New("cart has reached its listing limit")
		} else {
			p, e := queries.NewProductQueries().FindByID(id)
			if e != nil || p == nil || p.ID == "" {
				err = errors.New("listing not found")
			} else {
				_, err = p.PriceFor(qty)
			}
		}
	}
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
	} else {
		cart[id] = qty
		saveCart(ctx, cart)
		ctx.Request().Session().Flash("success", "Cart updated.")
	}
	return ctx.Response().Redirect(http.StatusSeeOther, "/cart")
}
func (c *CommerceController) RemoveItem(ctx http.Context) http.Response {
	cart := cartFrom(ctx)
	delete(cart, ctx.Request().Route("id"))
	saveCart(ctx, cart)
	return ctx.Response().Redirect(http.StatusSeeOther, "/cart")
}
func (c *CommerceController) Checkout(ctx http.Context) http.Response {
	cart := cartFrom(ctx)
	if id := ctx.Request().Input("product_id"); id != "" {
		qty, err := strconv.Atoi(ctx.Request().Input("quantity"))
		if err != nil {
			ctx.Request().Session().Flash("error", "Quantity must be a whole number.")
			return ctx.Response().Redirect(http.StatusSeeOther, "/cart")
		}
		cart = map[string]int{id: qty}
	}
	order, err := actions.NewExecuteCheckoutAction().ExecuteCart(cart, appmiddleware.CurrentUser(ctx))
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
		return ctx.Response().Redirect(http.StatusSeeOther, "/cart")
	}
	saveCart(ctx, map[string]int{})
	ctx.Request().Session().Flash("success", "Order "+order.OrderNumber+" submitted. Payment has not been collected; fulfillment is pending.")
	return ctx.Response().Redirect(http.StatusSeeOther, "/account?section=orders")
}
func (c *CommerceController) Account(ctx http.Context) http.Response {
	u := appmiddleware.CurrentUser(ctx)
	data, err := services.LoadWorkspace(u, true)
	if err != nil {
		return ctx.Response().String(http.StatusInternalServerError, "Your account could not be loaded. Please try again.")
	}
	currency := services.NewCurrencyService()
	return inertiafacades.Inertia().Render(ctx, "Account/Index", map[string]any{"currentUser": u, "orders": data.Orders, "offers": data.Offers, "rfqs": data.Rfqs, "currencies": currency.GetCurrencies(), "regions": currency.GetRegions()})
}
func (c *CommerceController) AcceptQuote(ctx http.Context) http.Response {
	return c.resolve(ctx, false)
}
func (c *CommerceController) CloseRfq(ctx http.Context) http.Response { return c.resolve(ctx, true) }
func (c *CommerceController) resolve(ctx http.Context, closeOnly bool) http.Response {
	err := services.ResolveRfq(ctx.Request().Route("id"), ctx.Request().Route("quoteId"), closeOnly, appmiddleware.CurrentUser(ctx))
	if err != nil {
		ctx.Request().Session().Flash("error", err.Error())
	} else {
		ctx.Request().Session().Flash("success", "Sourcing request updated.")
	}
	return ctx.Response().Redirect(http.StatusSeeOther, "/account?section=requests")
}

func (c *CommerceController) CancelOrder(ctx http.Context) http.Response {
	if err := services.CancelOrder(ctx.Request().Route("id"), appmiddleware.CurrentUser(ctx)); err != nil {
		ctx.Request().Session().Flash("error", err.Error())
	} else {
		ctx.Request().Session().Flash("success", "Order cancelled. Reserved inventory has been released; no payment was collected.")
	}
	return ctx.Response().Redirect(http.StatusSeeOther, "/account?section=orders")
}
