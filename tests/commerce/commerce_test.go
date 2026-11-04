package commerce_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"elcomware/app/actions"
	"elcomware/app/facades"
	"elcomware/app/http/middleware"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/queries"
	"elcomware/app/services"
	"elcomware/bootstrap"
	"encoding/json"
	"fmt"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"io"
	"math"
	nethttp "net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "soko-commerce-")
	if err != nil {
		panic(err)
	}
	_ = os.Setenv("DB_CONNECTION", "sqlite")
	_ = os.Setenv("DB_DATABASE", filepath.Join(dir, "test.sqlite"))
	facades.Config().Add("database.default", "sqlite")
	facades.Config().Add("database.connections.sqlite.database", filepath.Join(dir, "test.sqlite"))
	facades.Config().Add("session.files", filepath.Join(dir, "sessions"))
	bootstrap.Boot()
	for _, migration := range bootstrap.Migrations() {
		if err := migration.Up(); err != nil {
			panic(err)
		}
	}
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

func TestCheckoutReservesStockAndRemainsPending(t *testing.T) {
	p := models.Product{ID: "checkout-stock", Title: "Machine", Type: models.TypePhysical, Format: models.FormatWholesale, Price: 100, MOQ: 2, Stock: 5, PriceTiers: []models.PriceTier{{MinQty: 2, MaxQty: 5, Price: 80}}}
	p.PrepareSave()
	require.NoError(t, facades.Orm().Query().Create(&p))
	buyer := &models.User{ID: "buyer-1", Name: "Buyer", Email: "buyer@example.com", Role: models.RoleBuyer}
	req := &requests.CheckoutRequest{ProductID: p.ID, Quantity: 3, Currency: "KES", UnitPrice: 1}
	order, err := actions.NewExecuteCheckoutAction().ExecuteFromRequest(req, buyer)
	require.NoError(t, err)
	require.Equal(t, models.OrderStatusPending, order.Status)
	require.Equal(t, "USD", order.Currency)
	require.Equal(t, 240.0, order.TotalAmount)
	require.Empty(t, order.Items[0].DigitalLicenseKey)
	require.Empty(t, order.Items[0].DigitalDownloadUrl)
	var stored models.Product
	require.NoError(t, facades.Orm().Query().Where("id = ?", p.ID).First(&stored))
	require.Equal(t, 2, stored.Stock)
	_, err = actions.NewExecuteCheckoutAction().ExecuteFromRequest(req, buyer)
	require.Error(t, err)
	count, err := facades.Orm().Query().Model(&models.Order{}).Where("user_id = ?", buyer.ID).Count()
	require.NoError(t, err)
	require.EqualValues(t, 1, count)
}
func TestCheckoutRejectsMOQAndGuests(t *testing.T) {
	p := models.Product{ID: "checkout-moq", Title: "Bulk", Price: 10, MOQ: 5, Stock: 10, Type: models.TypePhysical}
	require.NoError(t, facades.Orm().Query().Create(&p))
	_, err := actions.NewExecuteCheckoutAction().ExecuteFromRequest(&requests.CheckoutRequest{ProductID: p.ID, Quantity: 1}, &models.User{ID: "buyer"})
	require.Error(t, err)
	_, err = actions.NewExecuteCheckoutAction().ExecuteFromRequest(&requests.CheckoutRequest{ProductID: p.ID, Quantity: 5}, nil)
	require.Error(t, err)
}
func TestDigitalCheckoutDoesNotGrantEntitlement(t *testing.T) {
	p := models.Product{ID: "checkout-digital", Title: "Software", Price: 10, MOQ: 1, Stock: -1, Type: models.TypeDigitalSoftware, Format: models.FormatDigital, DigitalDownloadUrl: "https://private.example/download"}
	require.NoError(t, facades.Orm().Query().Create(&p))
	order, err := actions.NewExecuteCheckoutAction().ExecuteFromRequest(&requests.CheckoutRequest{ProductID: p.ID, Quantity: 1}, &models.User{ID: "digital-buyer"})
	require.NoError(t, err)
	require.Equal(t, models.OrderStatusPending, order.Status)
	require.Empty(t, order.Items[0].DigitalLicenseKey)
	require.Empty(t, order.Items[0].DigitalDownloadUrl)
}
func TestBuyerCannotCreateProductAndZeroStockIsPreserved(t *testing.T) {
	req := &requests.CreateProductRequest{Title: "Valid product", Category: "Tools", Type: "physical", Format: "retail", Condition: "new", Price: 5, MOQ: 1, Stock: 0}
	_, err := actions.NewCreateProductAction().ExecuteFromRequest(req, &models.User{ID: "buyer-create", Role: models.RoleBuyer})
	require.Error(t, err)
	p, err := actions.NewCreateProductAction().ExecuteFromRequest(req, &models.User{ID: "seller-create", Role: models.RoleMerchant})
	require.NoError(t, err)
	require.Zero(t, p.Stock)
	require.False(t, p.VendorVerified)
	require.Zero(t, p.TradeAssuranceCover)
}
func TestPlaintextPasswordIsRejected(t *testing.T) {
	u := models.User{ID: "plaintext-user", Email: "plaintext@example.com", Password: "plain:secret123"}
	require.NoError(t, facades.Orm().Query().Create(&u))
	_, err := actions.NewAuthenticateUserAction().Execute(u.Email, "secret123")
	require.Error(t, err)
}

func TestCartCheckoutRollsBackEveryLine(t *testing.T) {
	a := models.Product{ID: "rollback-a", Price: 10, MOQ: 1, Stock: 3, Type: models.TypePhysical}
	b := models.Product{ID: "rollback-b", Price: 10, MOQ: 1, Stock: 0, Type: models.TypePhysical}
	require.NoError(t, facades.Orm().Query().Create(&a))
	require.NoError(t, facades.Orm().Query().Create(&b))
	_, err := actions.NewExecuteCheckoutAction().ExecuteCart(map[string]int{a.ID: 2, b.ID: 1}, &models.User{ID: "rollback-buyer"})
	require.Error(t, err)
	var stored models.Product
	require.NoError(t, facades.Orm().Query().Where("id = ?", a.ID).First(&stored))
	require.Equal(t, 3, stored.Stock)
	count, err := facades.Orm().Query().Model(&models.Order{}).Where("user_id = ?", "rollback-buyer").Count()
	require.NoError(t, err)
	require.Zero(t, count)
}
func TestCatalogOwnerMayEditButOtherActorsCannot(t *testing.T) {
	owner := &models.User{ID: "edit-owner", Role: models.RoleMerchant}
	other := &models.User{ID: "edit-other", Role: models.RoleMerchant}
	buyer := &models.User{ID: "edit-buyer", Role: models.RoleBuyer}
	r := &requests.CreateProductRequest{Title: "Owned item", Category: "Tools", Type: "physical", Format: "wholesale", Condition: "new", Price: 100, MOQ: 1, Stock: 4, AcceptsOffers: true}
	service := services.NewCatalogService()
	p, err := service.CreateProduct(r, owner)
	require.NoError(t, err)
	_, err = service.UpdateProduct(p.ID, r, other)
	require.Error(t, err)
	require.Error(t, service.DeleteProduct(p.ID, buyer))
	require.Error(t, service.DeleteProduct(p.ID, other))
	r.Version = p.UpdatedAt.Format(time.RFC3339Nano)
	r.Stock = 0
	r.Format = "retail"
	r.AcceptsOffers = false
	r.ImageUrl = ""
	updated, err := service.UpdateProduct(p.ID, r, owner)
	require.NoError(t, err)
	require.Zero(t, updated.Stock)
	require.False(t, updated.AcceptsOffers)
	require.Empty(t, updated.PriceTiers)
	require.Equal(t, owner.ID, updated.VendorID)
	var stored models.Product
	require.NoError(t, facades.Orm().Query().Where("id = ?", p.ID).First(&stored))
	stored.Hydrate()
	require.Empty(t, stored.PriceTiers)
	require.Zero(t, stored.Stock)
}
func TestWorkspaceScopesOrdersOffersUsersAndQuotes(t *testing.T) {
	buyerA := &models.User{ID: "scope-buyer-a", Role: models.RoleBuyer}
	buyerB := &models.User{ID: "scope-buyer-b", Role: models.RoleBuyer}
	seller := &models.User{ID: "scope-seller", Role: models.RoleMerchant}
	p := models.Product{ID: "scope-product", VendorID: seller.ID, Price: 20, MOQ: 1, Stock: 10, Type: models.TypePhysical}
	require.NoError(t, facades.Orm().Query().Create(&p))
	for _, u := range []*models.User{buyerA, buyerB} {
		_, err := actions.NewExecuteCheckoutAction().ExecuteCart(map[string]int{p.ID: 1}, u)
		require.NoError(t, err)
	}
	offer := models.Offer{ID: "scope-offer", BuyerID: buyerB.ID, ProductID: p.ID, Status: models.OfferStatusPending}
	require.NoError(t, facades.Orm().Query().Create(&offer))
	rfq := models.Rfq{ID: "scope-rfq", RfqNumber: "SCOPE-RFQ", BuyerID: buyerB.ID, Status: models.RfqStatusOpen}
	require.NoError(t, facades.Orm().Query().Create(&rfq))
	for _, v := range []string{seller.ID, "scope-other-seller"} {
		quote := models.RfqQuote{ID: "scope-quote-" + v, RfqID: rfq.ID, VendorID: v, UnitPrice: 1, Status: "pending"}
		require.NoError(t, facades.Orm().Query().Create(&quote))
	}
	data, err := services.LoadWorkspace(buyerA, true)
	require.NoError(t, err)
	require.Len(t, data.Orders, 1)
	require.Equal(t, buyerA.ID, data.Orders[0].UserID)
	require.Empty(t, data.Users)
	require.Empty(t, data.Offers)
	require.Empty(t, data.Rfqs)
	data, err = services.LoadWorkspace(seller, false)
	require.NoError(t, err)
	require.Empty(t, data.Users)
	require.Len(t, data.Products, 1)
	require.Len(t, data.Offers, 1)
	found := false
	for _, r := range data.Rfqs {
		if r.ID == rfq.ID {
			found = true
			require.Len(t, r.Quotes, 1)
			require.Equal(t, seller.ID, r.Quotes[0].VendorID)
		}
	}
	require.True(t, found)
}
func TestOfferDecisionsAndRfqAwardRequireCorrectActors(t *testing.T) {
	seller := &models.User{ID: "decision-seller", Role: models.RoleMerchant}
	other := &models.User{ID: "decision-other", Role: models.RoleMerchant}
	buyer := &models.User{ID: "decision-buyer", Role: models.RoleBuyer}
	p := models.Product{ID: "decision-product", VendorID: seller.ID, Price: 10, MOQ: 1, Stock: 10, Type: models.TypePhysical, AcceptsOffers: true}
	require.NoError(t, facades.Orm().Query().Create(&p))
	offer, err := actions.NewSubmitOfferAction().ExecuteFromRequest(&requests.SubmitOfferRequest{ProductID: p.ID, Quantity: 1, OfferedPrice: 9}, buyer)
	require.NoError(t, err)
	require.Error(t, services.DecideOffer(offer.ID, "accepted", buyer))
	require.Error(t, services.DecideOffer(offer.ID, "accepted", other))
	require.NoError(t, services.DecideOffer(offer.ID, "accepted", seller))
	require.Error(t, services.DecideOffer(offer.ID, "declined", seller))
	rfq := models.Rfq{ID: "decision-rfq", RfqNumber: "DECISION-RFQ", BuyerID: buyer.ID, Status: models.RfqStatusOpen, Quantity: 5}
	require.NoError(t, facades.Orm().Query().Create(&rfq))
	request := &requests.SubmitRfqQuoteRequest{RfqID: rfq.ID, UnitPrice: 4, LeadTimeDays: 2}
	_, err = actions.NewSubmitRfqQuoteAction().ExecuteFromRequest(request, buyer)
	require.Error(t, err)
	quote, err := actions.NewSubmitRfqQuoteAction().ExecuteFromRequest(request, seller)
	require.NoError(t, err)
	request.UnitPrice = 3
	replacement, err := actions.NewSubmitRfqQuoteAction().ExecuteFromRequest(request, seller)
	require.NoError(t, err)
	require.Equal(t, quote.ID, replacement.ID)
	require.Equal(t, 15.0, replacement.TotalAmount)
	require.Error(t, services.ResolveRfq(rfq.ID, quote.ID, false, other))
	require.NoError(t, services.ResolveRfq(rfq.ID, quote.ID, false, buyer))
	require.Error(t, services.ResolveRfq(rfq.ID, "", true, buyer))
	_, err = actions.NewSubmitRfqQuoteAction().ExecuteFromRequest(request, seller)
	require.Error(t, err)
	var stored models.RfqQuote
	require.NoError(t, facades.Orm().Query().Where("id = ?", quote.ID).First(&stored))
	require.Equal(t, "accepted", stored.Status)
}
func TestPublicCatalogAndSourcingBoardDoNotExposeSecrets(t *testing.T) {
	p := models.Product{ID: "private-download", DigitalDownloadUrl: "https://secret.example/file", DigitalLicenseFormat: "secret-format"}
	raw, err := json.Marshal(p.Public())
	require.NoError(t, err)
	require.NotContains(t, string(raw), "secret")
	rfqs, err := services.NewRfqService().GetOpenRfqs()
	require.NoError(t, err)
	for _, r := range rfqs {
		require.Equal(t, models.RfqStatusOpen, r.Status)
		require.Empty(t, r.Quotes)
		require.Empty(t, r.BuyerID)
		require.Empty(t, r.BuyerName)
	}
}
func TestContentPersistenceUsesSingletonAndRetainsDefaults(t *testing.T) {
	var row models.StorefrontContent
	require.NoError(t, facades.Orm().Query().UpdateOrCreate(&row, models.StorefrontContent{ID: "store"}, map[string]any{"content_json": "{\"hero_heading\":\"Edited hero\"}"}))
	require.Equal(t, "Edited hero", queries.LoadStorefrontContent()["hero_heading"])
	require.NotEmpty(t, queries.LoadStorefrontContent()["hero_accent"])
}
func TestSameOriginMutationsRejectCrossOriginBrowserRequests(t *testing.T) {
	cases := []struct {
		name, origin, referer, site, cookie string
		allowed                             bool
	}{
		{"same origin", "http://localhost:8000", "", "same-origin", "session=x", true}, {"foreign", "https://evil.example", "", "cross-site", "session=x", false}, {"same site subdomain", "http://other.localhost:8000", "", "same-site", "session=x", false}, {"foreign port", "http://localhost:9000", "", "same-origin", "session=x", false}, {"null origin", "null", "", "", "session=x", false}, {"referer fallback", "", "http://localhost:8000/cart", "", "session=x", true}, {"cookie without metadata", "", "", "", "session=x", false}, {"nonbrowser without cookie", "", "", "", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.allowed, middleware.ValidMutationOrigin("POST", "localhost:8000", "http", tc.origin, tc.referer, tc.site, tc.cookie))
		})
	}
}

func TestHTTPAccountIdentityAndContentRoleChecks(t *testing.T) {
	users := []models.User{{ID: "http-buyer", Name: "Buyer", Email: "httpbuyer@example.com", Role: models.RoleBuyer}, {ID: "http-seller", Name: "Seller", Email: "httpseller@example.com", Role: models.RoleMerchant}, {ID: "http-admin", Name: "Admin", Email: "httpadmin@example.com", Role: models.RoleAdmin}}
	for i := range users {
		require.NoError(t, facades.Orm().Query().Create(&users[i]))
	}
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithSession(map[string]any{"user_id": "does-not-exist", "user_role": "admin"}).Get("/dashboard")
	require.NoError(t, err)
	response.AssertFound().AssertHeader("Location", "/login")
	response, err = tc.Http(t).WithSession(map[string]any{"user_id": users[0].ID}).Get("/dashboard")
	require.NoError(t, err)
	response.AssertFound().AssertHeader("Location", "/account")
	response, err = tc.Http(t).WithHeader("X-Inertia", "true").WithSession(map[string]any{"user_id": users[0].ID}).Get("/account")
	require.NoError(t, err)
	response.AssertOk()
	body, err := response.Json()
	require.NoError(t, err)
	props := body["props"].(map[string]any)
	require.Equal(t, users[0].ID, props["currentUser"].(map[string]any)["id"])
	for _, u := range users[:2] {
		response, err = tc.Http(t).WithHeader("Sec-Fetch-Site", "same-origin").WithSession(map[string]any{"user_id": u.ID, "user_role": "admin"}).Post("/dashboard/content", strings.NewReader(""))
		require.NoError(t, err)
		response.AssertForbidden()
	}
	values := queries.DefaultStorefrontContent()
	values["hero_heading"] = "HTTP edited hero"
	encoded, err := json.Marshal(values)
	require.NoError(t, err)
	response, err = tc.Http(t).WithHeaders(map[string]string{"Sec-Fetch-Site": "same-origin", "Content-Type": "application/json"}).WithSession(map[string]any{"user_id": users[2].ID}).Post("/dashboard/content", strings.NewReader(string(encoded)))
	require.NoError(t, err)
	response.AssertFound()
	require.Equal(t, "HTTP edited hero", queries.LoadStorefrontContent()["hero_heading"])
	response, err = tc.Http(t).WithHeaders(map[string]string{"Origin": "https://evil.example", "Sec-Fetch-Site": "cross-site"}).WithSession(map[string]any{"user_id": users[2].ID}).Post("/dashboard/content", strings.NewReader(""))
	require.NoError(t, err)
	response.AssertForbidden()
}
func TestHTTPCartRejectsFractionsAndCheckoutUsesSessionIdentity(t *testing.T) {
	p := models.Product{ID: "http-cart", Title: "HTTP cart listing", VendorID: "http-seller", Type: models.TypePhysical, Format: models.FormatRetail, Price: 20, MOQ: 1, Stock: 5}
	require.NoError(t, facades.Orm().Query().Create(&p))
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).Post("/cart/items", strings.NewReader(`{"product_id":"http-cart","quantity":1.5}`))
	require.NoError(t, err)
	response.AssertStatus(303)
	follow, err := tc.Http(t).WithCookies(response.Cookies()).WithHeader("X-Inertia", "true").Get("/cart")
	require.NoError(t, err)
	follow.AssertOk()
	body, err := follow.Json()
	require.NoError(t, err)
	props := body["props"].(map[string]any)
	cart := props["cart"].(map[string]any)
	require.Empty(t, cart["items"])
	require.NotEmpty(t, props["flash"].(map[string]any)["error"])
	user := models.User{ID: "http-checkout-buyer", Name: "Actual buyer", Email: "actual@example.com", Role: models.RoleBuyer}
	require.NoError(t, facades.Orm().Query().Create(&user))
	response, err = tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": user.ID, "cart": `{"http-cart":2}`}).Post("/checkout", strings.NewReader(`{"unit_price":1,"currency":"KES","customer_email":"spoof@example.com"}`))
	require.NoError(t, err)
	response.AssertStatus(303).AssertHeader("Location", "/account?section=orders")
	var order models.Order
	require.NoError(t, facades.Orm().Query().Where("user_id = ?", user.ID).First(&order))
	require.Equal(t, 40.0, order.TotalAmount)
	require.Equal(t, user.Email, order.CustomerEmail)
	require.Equal(t, "USD", order.Currency)
	require.Equal(t, models.OrderStatusPending, order.Status)
}

func TestHTTPPendingOrderCancellationRestoresInventoryExactlyOnce(t *testing.T) {
	owner := models.User{ID: "cancel-owner", Name: "Owner", Email: "cancelowner@example.com", Role: models.RoleBuyer}
	other := models.User{ID: "cancel-other", Name: "Other", Email: "cancelother@example.com", Role: models.RoleBuyer}
	admin := models.User{ID: "cancel-admin", Name: "Admin", Email: "canceladmin@example.com", Role: models.RoleAdmin}
	for _, u := range []*models.User{&owner, &other, &admin} {
		require.NoError(t, facades.Orm().Query().Create(u))
	}
	physical := models.Product{ID: "cancel-physical", Title: "Machine", Type: models.TypePhysical, MOQ: 1, Stock: 7, Price: 10}
	digital := models.Product{ID: "cancel-digital", Title: "Software", Type: models.TypeDigitalSoftware, Format: models.FormatDigital, MOQ: 1, Stock: -1, Price: 5}
	for _, p := range []*models.Product{&physical, &digital} {
		require.NoError(t, facades.Orm().Query().Create(p))
	}
	order, err := actions.NewExecuteCheckoutAction().ExecuteCart(map[string]int{physical.ID: 3, digital.ID: 2}, &owner)
	require.NoError(t, err)
	tc := &frameworktesting.TestCase{}
	cancel := func(actor string) {
		response, err := tc.Http(t).WithHeader("Sec-Fetch-Site", "same-origin").WithSession(map[string]any{"user_id": actor}).Post("/account/orders/"+order.ID+"/cancel", strings.NewReader(""))
		require.NoError(t, err)
		response.AssertStatus(303).AssertHeader("Location", "/account?section=orders")
	}
	cancel(other.ID)
	var stored models.Order
	require.NoError(t, facades.Orm().Query().Where("id = ?", order.ID).First(&stored))
	require.Equal(t, models.OrderStatusPending, stored.Status)
	cancel(owner.ID)
	require.NoError(t, facades.Orm().Query().Where("id = ?", order.ID).First(&stored))
	require.Equal(t, models.OrderStatus("cancelled"), stored.Status)
	cancel(owner.ID)
	var stock models.Product
	require.NoError(t, facades.Orm().Query().Where("id = ?", physical.ID).First(&stock))
	require.Equal(t, 7, stock.Stock)
	stock = models.Product{}
	require.NoError(t, facades.Orm().Query().Where("id = ?", digital.ID).First(&stock))
	require.Equal(t, -1, stock.Stock)
	count, err := facades.Orm().Query().Model(&models.OrderItem{}).Where("order_id = ?", order.ID).Count()
	require.NoError(t, err)
	require.EqualValues(t, 2, count)
	second, err := actions.NewExecuteCheckoutAction().ExecuteCart(map[string]int{physical.ID: 1}, &owner)
	require.NoError(t, err)
	response, err := tc.Http(t).WithHeader("Sec-Fetch-Site", "same-origin").WithSession(map[string]any{"user_id": admin.ID}).Post("/account/orders/"+second.ID+"/cancel", strings.NewReader(""))
	require.NoError(t, err)
	response.AssertStatus(303)
	stored = models.Order{}
	require.NoError(t, facades.Orm().Query().Where("id = ?", second.ID).First(&stored))
	require.Equal(t, models.OrderStatus("cancelled"), stored.Status)
}

func TestRealCookieJarLoginRotatesSessionAndPreservesCart(t *testing.T) {
	hash, err := facades.Hash().Make("secret123")
	require.NoError(t, err)
	user := models.User{ID: "jar-buyer", Name: "Cookie Buyer", Email: "jarbuyer@example.com", Password: hash, Role: models.RoleBuyer}
	require.NoError(t, facades.Orm().Query().Create(&user))
	p := models.Product{ID: "jar-product", Title: "Cookie cart listing", Type: models.TypePhysical, Format: models.FormatRetail, Price: 20, MOQ: 1, Stock: 4}
	require.NoError(t, facades.Orm().Query().Create(&p))
	server := httptest.NewServer(facades.Route())
	defer server.Close()
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := &nethttp.Client{Jar: jar, CheckRedirect: func(*nethttp.Request, []*nethttp.Request) error { return nethttp.ErrUseLastResponse }}
	endpoint, err := url.Parse(server.URL)
	require.NoError(t, err)
	call := func(method, path, body string) *nethttp.Response {
		req, err := nethttp.NewRequest(method, server.URL+path, strings.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("X-Inertia", "true")
		req.Header.Set("Content-Type", "application/json")
		if method != "GET" {
			req.Header.Set("Origin", server.URL)
			req.Header.Set("Sec-Fetch-Site", "same-origin")
		}
		response, err := client.Do(req)
		require.NoError(t, err)
		return response
	}
	closeResponse := func(r *nethttp.Response) { _, _ = io.Copy(io.Discard, r.Body); require.NoError(t, r.Body.Close()) }
	response := call("GET", "/login", "")
	require.Equal(t, 200, response.StatusCode)
	closeResponse(response)
	response = call("POST", "/cart/items", `{"product_id":"jar-product","quantity":2}`)
	require.Equal(t, 303, response.StatusCode)
	closeResponse(response)
	cookies := jar.Cookies(endpoint)
	require.NotEmpty(t, cookies)
	oldID := cookies[0].Value
	cookieName := cookies[0].Name
	response = call("POST", "/login", `{"email":"jarbuyer@example.com","password":"secret123"}`)
	require.Equal(t, 302, response.StatusCode)
	require.Equal(t, "/account", response.Header.Get("Location"))
	closeResponse(response)
	require.NotEqual(t, oldID, jar.Cookies(endpoint)[0].Value, "login must send its rotated session ID to the browser")
	response = call("GET", "/account", "")
	require.Equal(t, 200, response.StatusCode)
	var page map[string]any
	require.NoError(t, json.NewDecoder(response.Body).Decode(&page))
	require.NoError(t, response.Body.Close())
	require.Equal(t, user.ID, page["props"].(map[string]any)["currentUser"].(map[string]any)["id"])
	response = call("GET", "/cart", "")
	require.Equal(t, 200, response.StatusCode)
	page = nil
	require.NoError(t, json.NewDecoder(response.Body).Decode(&page))
	require.NoError(t, response.Body.Close())
	require.Len(t, page["props"].(map[string]any)["cart"].(map[string]any)["items"], 1)
	staleRequest, err := nethttp.NewRequest("GET", server.URL+"/account", nil)
	require.NoError(t, err)
	staleRequest.AddCookie(&nethttp.Cookie{Name: cookieName, Value: oldID})
	staleClient := &nethttp.Client{CheckRedirect: client.CheckRedirect}
	staleResponse, err := staleClient.Do(staleRequest)
	require.NoError(t, err)
	require.Equal(t, 302, staleResponse.StatusCode)
	require.Equal(t, "/login", staleResponse.Header.Get("Location"))
	closeResponse(staleResponse)
	response = call("POST", "/logout", `{}`)
	require.Equal(t, 302, response.StatusCode)
	closeResponse(response)
	response = call("GET", "/account", "")
	require.Equal(t, 302, response.StatusCode)
	require.Equal(t, "/login", response.Header.Get("Location"))
	closeResponse(response)
}

func TestRealRegistrationRotatesCookieAndCannotChooseAdministrator(t *testing.T) {
	server := httptest.NewServer(facades.Route())
	defer server.Close()
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := &nethttp.Client{Jar: jar, CheckRedirect: func(*nethttp.Request, []*nethttp.Request) error { return nethttp.ErrUseLastResponse }}
	initial, err := client.Get(server.URL + "/login")
	require.NoError(t, err)
	_, _ = io.Copy(io.Discard, initial.Body)
	require.NoError(t, initial.Body.Close())
	endpoint, err := url.Parse(server.URL)
	require.NoError(t, err)
	oldID := jar.Cookies(endpoint)[0].Value
	register := func(password string) *nethttp.Response {
		body, _ := json.Marshal(map[string]string{"name": "New Buyer", "email": "registered@example.com", "password": password, "account_type": "buyer", "role": "admin"})
		req, err := nethttp.NewRequest("POST", server.URL+"/register", strings.NewReader(string(body)))
		require.NoError(t, err)
		req.Header.Set("Origin", server.URL)
		req.Header.Set("Content-Type", "application/json")
		r, err := client.Do(req)
		require.NoError(t, err)
		_, _ = io.Copy(io.Discard, r.Body)
		require.NoError(t, r.Body.Close())
		return r
	}
	response := register("short77")
	require.Equal(t, 302, response.StatusCode)
	require.Equal(t, "/register", response.Header.Get("Location"))
	response = register("secret123")
	require.Equal(t, 302, response.StatusCode)
	require.Equal(t, "/account", response.Header.Get("Location"))
	require.NotEqual(t, oldID, jar.Cookies(endpoint)[0].Value)
	var user models.User
	require.NoError(t, facades.Orm().Query().Where("email = ?", "registered@example.com").First(&user))
	require.NotEmpty(t, user.ID)
	require.Equal(t, models.RoleBuyer, user.Role)
	require.True(t, facades.Hash().Check("secret123", user.Password))
	require.NotEqual(t, "secret123", user.Password)
	req, err := nethttp.NewRequest("GET", server.URL+"/account", nil)
	require.NoError(t, err)
	req.Header.Set("X-Inertia", "true")
	response, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.NoError(t, response.Body.Close())
}

func TestCatalogRejectsStaleEditsAndRetainsOrderAndOfferHistory(t *testing.T) {
	owner := &models.User{ID: "history-owner", Role: models.RoleMerchant}
	buyer := &models.User{ID: "history-buyer", Role: models.RoleBuyer}
	r := &requests.CreateProductRequest{Title: "History listing", Category: "Tools", Type: "physical", Format: "retail", Condition: "new", Price: 10, MOQ: 1, Stock: 5}
	catalog := services.NewCatalogService()
	p, err := catalog.CreateProduct(r, owner)
	require.NoError(t, err)
	r.Version = p.UpdatedAt.Format(time.RFC3339Nano)
	_, err = actions.NewExecuteCheckoutAction().ExecuteCart(map[string]int{p.ID: 2}, buyer)
	require.NoError(t, err)
	_, err = catalog.UpdateProduct(p.ID, r, owner)
	require.ErrorContains(t, err, "changed")
	stored, err := catalog.GetProductByID(p.ID)
	require.NoError(t, err)
	require.Equal(t, 3, stored.Stock)
	require.ErrorContains(t, catalog.DeleteProduct(p.ID, owner), "history")
	other := models.Product{ID: "offer-history-product", VendorID: owner.ID, Price: 10, MOQ: 1, Stock: 3, Type: models.TypePhysical}
	require.NoError(t, facades.Orm().Query().Create(&other))
	offer := models.Offer{ID: "offer-history", ProductID: other.ID, BuyerID: buyer.ID}
	require.NoError(t, facades.Orm().Query().Create(&offer))
	require.ErrorContains(t, catalog.DeleteProduct(other.ID, owner), "history")
}
func TestFiniteDigitalStockReservesAndCancellationRestoresOnlyReservations(t *testing.T) {
	buyer := &models.User{ID: "finite-digital-buyer", Role: models.RoleBuyer}
	finite := models.Product{ID: "finite-digital", Type: models.TypeDigitalSoftware, Format: models.FormatDigital, Price: 5, MOQ: 1, Stock: 3}
	unlimited := models.Product{ID: "unlimited-digital", Type: models.TypeDigitalSoftware, Format: models.FormatDigital, Price: 5, MOQ: 1, Stock: -1}
	for _, p := range []*models.Product{&finite, &unlimited} {
		require.NoError(t, facades.Orm().Query().Create(p))
	}
	order, err := actions.NewExecuteCheckoutAction().ExecuteCart(map[string]int{finite.ID: 2, unlimited.ID: 2}, buyer)
	require.NoError(t, err)
	var stock models.Product
	require.NoError(t, facades.Orm().Query().Where("id = ?", finite.ID).First(&stock))
	require.Equal(t, 1, stock.Stock)
	_, err = actions.NewExecuteCheckoutAction().ExecuteCart(map[string]int{finite.ID: 2}, buyer)
	require.Error(t, err)
	_, err = facades.Orm().Query().Model(&models.Product{}).Where("id = ?", unlimited.ID).Update("stock", 10)
	require.NoError(t, err)
	require.NoError(t, services.CancelOrder(order.ID, buyer))
	stock = models.Product{}
	require.NoError(t, facades.Orm().Query().Where("id = ?", finite.ID).First(&stock))
	require.Equal(t, 3, stock.Stock)
	stock = models.Product{}
	require.NoError(t, facades.Orm().Query().Where("id = ?", unlimited.ID).First(&stock))
	require.Equal(t, 10, stock.Stock)
}
func TestRoundedZeroPriceAndExhaustedDigitalStockAreRejected(t *testing.T) {
	p := models.Product{Type: models.TypeDigitalAsset, Format: models.FormatDigital, Price: .001, MOQ: 1, Stock: -1}
	_, err := p.PriceFor(1)
	require.Error(t, err)
	p.Price = 5
	p.Stock = 0
	_, err = p.PriceFor(1)
	require.Error(t, err)
}

func TestFriendlyRegionDataMatchesCanonicalAndUnderscoreFilters(t *testing.T) {
	owner := &models.User{ID: "region-owner", Role: models.RoleMerchant}
	spellings := []string{"East Africa", "east-africa", "east_africa"}
	for i, region := range spellings {
		p := models.Product{ID: fmt.Sprintf("region-fixture-%d", i), Title: "Region fixture listing", VendorID: owner.ID, VendorRegion: region, Type: models.TypePhysical, Format: models.FormatRetail, Price: 5, MOQ: 1, Stock: 3}
		require.NoError(t, facades.Orm().Query().Create(&p))
	}
	catalog := services.NewCatalogService()
	for _, region := range []string{"east-africa", "east_africa", "East Africa"} {
		products, err := catalog.GetProducts(queries.ProductFilter{Region: region, Query: "Region fixture"})
		require.NoError(t, err)
		require.Len(t, products, 3)
	}
	original, err := catalog.GetProductByID("region-fixture-0")
	require.NoError(t, err)
	require.Equal(t, "East Africa", original.VendorRegion)
	require.Equal(t, "east-africa", original.Public().VendorRegion)
	workspace, err := services.LoadWorkspace(owner, false)
	require.NoError(t, err)
	require.Len(t, workspace.Products, 3)
	for _, p := range workspace.Products {
		require.Equal(t, "east-africa", p.VendorRegion)
	}
	input := &requests.CreateProductRequest{Title: "Region fixture edited", Category: "Tools", Region: "East Africa", Type: "physical", Format: "retail", Condition: "new", Price: 5, MOQ: 1, Stock: 3, Version: original.UpdatedAt.Format(time.RFC3339Nano)}
	edited, err := catalog.UpdateProduct(original.ID, input, owner)
	require.NoError(t, err)
	require.Equal(t, "east-africa", edited.VendorRegion)
	input.Region = "Asia-Pacific"
	created, err := catalog.CreateProduct(input, owner)
	require.NoError(t, err)
	require.Equal(t, "asia-pacific", created.VendorRegion)
	products, err := catalog.GetProducts(queries.ProductFilter{Region: "asia_pacific", Query: "Region fixture"})
	require.NoError(t, err)
	require.Len(t, products, 1)
}

func TestTitleOnlyEditPreservesActualVolumePricesAndExplicitEmptyClears(t *testing.T) {
	owner := &models.User{ID: "tier-owner", Role: models.RoleMerchant}
	p := models.Product{ID: "tier-cnc", Title: "CNC machine", VendorID: owner.ID, Category: "Tools", Type: models.TypePhysical, Format: models.FormatWholesale, Condition: models.ConditionNew, MOQ: 1, Stock: 50, Price: 48500, WholesalePrice: 39500, PriceTiers: []models.PriceTier{{MinQty: 1, MaxQty: 2, Price: 48500}, {MinQty: 3, MaxQty: 5, Price: 44000}, {MinQty: 6, MaxQty: 0, Price: 39500}}}
	p.PrepareSave()
	require.NoError(t, facades.Orm().Query().Create(&p))
	original := append([]models.PriceTier{}, p.PriceTiers...)
	r := &requests.CreateProductRequest{Title: "CNC machine renamed", Category: "Tools", Type: "physical", Format: "wholesale", Condition: "new", MOQ: 1, Stock: 50, Price: 48500, WholesalePrice: 39500, Version: p.UpdatedAt.Format(time.RFC3339Nano)}
	catalog := services.NewCatalogService()
	updated, err := catalog.UpdateProduct(p.ID, r, owner)
	require.NoError(t, err)
	require.Equal(t, original, updated.PriceTiers)
	price, err := updated.PriceFor(3)
	require.NoError(t, err)
	require.Equal(t, 44000.0, price)
	r.PriceTiers = []models.PriceTier{}
	r.Version = updated.UpdatedAt.Format(time.RFC3339Nano)
	updated, err = catalog.UpdateProduct(p.ID, r, owner)
	require.NoError(t, err)
	require.Empty(t, updated.PriceTiers)
	price, err = updated.PriceFor(3)
	require.NoError(t, err)
	require.Equal(t, 48500.0, price)
	r.PriceTiers = nil
	created, err := catalog.CreateProduct(r, owner)
	require.NoError(t, err)
	require.Empty(t, created.PriceTiers)
}
func TestInvalidVolumePriceTiersAreRejected(t *testing.T) {
	cases := [][]models.PriceTier{{{MinQty: 0, MaxQty: 4, Price: 10}}, {{MinQty: 1, MaxQty: 0, Price: 10}}, {{MinQty: 2, MaxQty: 1, Price: 10}}, {{MinQty: 2, MaxQty: 4, Price: 10}, {MinQty: 4, MaxQty: 6, Price: 9}}, {{MinQty: 2, MaxQty: 0, Price: 10}, {MinQty: 5, MaxQty: 0, Price: 9}}, {{MinQty: 2, MaxQty: 1000001, Price: 10}}, {{MinQty: 2, MaxQty: 4, Price: math.Inf(1)}}, {{MinQty: 2, MaxQty: 4, Price: .001}}}
	for i, tiers := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			r := &requests.CreateProductRequest{Title: "Valid listing", Category: "Tools", Type: "physical", Format: "wholesale", Condition: "new", MOQ: 2, Stock: 50, Price: 10, PriceTiers: tiers}
			require.Error(t, actions.ValidateProduct(r))
		})
	}
	r := &requests.CreateProductRequest{Title: "Valid listing", Category: "Tools", Type: "physical", Format: "wholesale", Condition: "new", MOQ: 1, Stock: 50, Price: 10}
	for i := 1; i <= 21; i++ {
		r.PriceTiers = append(r.PriceTiers, models.PriceTier{MinQty: i, MaxQty: i, Price: 1})
	}
	require.Error(t, actions.ValidateProduct(r))
}

func TestHTTPProductCreationBindsSuppliedVolumeTiers(t *testing.T) {
	seller := models.User{ID: "http-tier-seller", Name: "Tier Seller", Email: "httptiers@example.com", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&seller))
	payload := map[string]any{"title": "HTTP volume listing", "category": "Tools", "type": "physical", "format": "wholesale", "condition": "new", "region": "global", "price": 10, "wholesale_price": 8, "moq": 2, "stock": 20, "price_tiers": []map[string]any{{"min_qty": 2, "max_qty": 5, "price": 10}, {"min_qty": 6, "max_qty": 0, "price": 8}}}
	body, err := json.Marshal(payload)
	require.NoError(t, err)
	tc := &frameworktesting.TestCase{}
	response, err := tc.Http(t).WithHeaders(map[string]string{"Content-Type": "application/json", "Sec-Fetch-Site": "same-origin"}).WithSession(map[string]any{"user_id": seller.ID}).Post("/dashboard/products", strings.NewReader(string(body)))
	require.NoError(t, err)
	response.AssertFound().AssertHeader("Location", "/dashboard?section=products")
	var stored models.Product
	require.NoError(t, facades.Orm().Query().Where("vendor_id = ?", seller.ID).First(&stored))
	require.NotEmpty(t, stored.ID, "HTTP create must bind nested quantity fields and persist the listing")
	stored.Hydrate()
	require.Equal(t, []models.PriceTier{{MinQty: 2, MaxQty: 5, Price: 10}, {MinQty: 6, MaxQty: 0, Price: 8}}, stored.PriceTiers)
	unitPrice, err := stored.PriceFor(6)
	require.NoError(t, err)
	require.Equal(t, 8.0, unitPrice)
}
