package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/services"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWorkspaceOrdersScopePageAndSelection(t *testing.T) {
	actor := &models.User{ID: "paged-order-vendor", Role: models.RoleMerchant}
	for _, p := range []models.Product{{ID: "wo-own", VendorID: actor.ID}, {ID: "wo-foreign", VendorID: "other"}} {
		require.NoError(t, facades.Orm().Query().Create(&p))
	}
	stamp := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, o := range []models.Order{{ID: "wo-a", OrderNumber: "WO-A", CustomerName: "ÉCOLE", Status: models.OrderStatusPending, TotalAmount: 999, CreatedAt: stamp}, {ID: "wo-b", OrderNumber: "WO-B", Status: models.OrderStatusShipped, CreatedAt: stamp}, {ID: "wo-secret", OrderNumber: "WO-SECRET", CreatedAt: stamp}} {
		require.NoError(t, facades.Orm().Query().Create(&o))
	}
	for _, i := range []models.OrderItem{{ID: "wo-item-a", OrderID: "wo-a", ProductID: "wo-own", TotalPrice: 10, ProductType: models.TypePhysical, DigitalLicenseKey: "secret", DigitalDownloadUrl: "secret"}, {ID: "wo-item-other", OrderID: "wo-a", ProductID: "wo-foreign", TotalPrice: 989, ProductType: models.TypeDigitalLicense}, {ID: "wo-item-b", OrderID: "wo-b", ProductID: "wo-own", TotalPrice: 20}, {ID: "wo-item-secret", OrderID: "wo-secret", ProductID: "wo-foreign"}} {
		require.NoError(t, facades.Orm().Query().Create(&i))
	}
	f, err := services.ParseWorkspaceOrderFilters(map[string]string{"per_page": "1", "selected_order": "wo-a"})
	require.NoError(t, err)
	page, err := services.LoadWorkspaceOrders(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, 2, page.OrderPage.Total)
	require.EqualValues(t, 2, page.OrderPage.FilteredTotal)
	require.Equal(t, "wo-b", page.Orders[0].ID)
	require.Equal(t, "supplier", page.OrderVisibility)
	require.Equal(t, "available", page.SelectedOrderStatus)
	require.Len(t, page.SelectedOrder.Items, 1)
	require.Equal(t, float64(10), page.SelectedOrder.TotalAmount)
	require.False(t, page.SelectedOrder.HasDigital)
	require.True(t, page.SelectedOrder.HasPhysical)
	require.Empty(t, page.SelectedOrder.Items[0].DigitalLicenseKey)
	require.Empty(t, page.SelectedOrder.Items[0].DigitalDownloadUrl)
	require.False(t, page.RecordCapabilities["orders"]["wo-a"]["cancel"])
	f.Query = "école"
	f.Page = 999
	f.SelectedOrder = "wo-secret"
	page, err = services.LoadWorkspaceOrders(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, 1, page.OrderPage.FilteredTotal)
	require.Equal(t, 1, page.OrderPage.Page)
	require.Equal(t, "wo-a", page.Orders[0].ID)
	require.Nil(t, page.SelectedOrder)
	require.Equal(t, "unavailable", page.SelectedOrderStatus)
	f.Query = "WO-SECRET"
	page, err = services.LoadWorkspaceOrders(actor, f)
	require.NoError(t, err)
	require.Zero(t, page.OrderPage.FilteredTotal)
	require.Empty(t, page.Orders)
	f.Query = ""
	f.Status = "shipped"
	page, err = services.LoadWorkspaceOrders(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, 1, page.OrderPage.FilteredTotal)
	f.Status = ""
	f.Sort = "oldest"
	f.Page = 1
	page, err = services.LoadWorkspaceOrders(actor, f)
	require.NoError(t, err)
	require.Equal(t, "wo-a", page.Orders[0].ID)
	admin := &models.User{ID: "wo-admin", Role: models.RoleAdmin}
	f.SelectedOrder = "wo-a"
	page, err = services.LoadWorkspaceOrders(admin, f)
	require.NoError(t, err)
	require.Equal(t, "full", page.OrderVisibility)
	require.Len(t, page.SelectedOrder.Items, 2)
	require.Equal(t, "secret", page.SelectedOrder.Items[0].DigitalLicenseKey)
	require.Equal(t, float64(999), page.SelectedOrder.TotalAmount)
	require.True(t, page.RecordCapabilities["orders"]["wo-a"]["cancel"])
}

func TestWorkspaceOrderFiltersRejectInvalid(t *testing.T) {
	for _, input := range []map[string]string{{"page": "0"}, {"per_page": "101"}, {"sort": "amount"}, {"status": "paid"}, {"q": "bad\nquery"}, {"page": "NaN"}} {
		_, err := services.ParseWorkspaceOrderFilters(input)
		require.Error(t, err)
	}
	f, err := services.ParseWorkspaceOrderFilters(nil)
	require.NoError(t, err)
	_, err = services.LoadWorkspaceOrders(nil, f)
	require.Error(t, err)
}

func TestWorkspaceOrdersUnicodeLiteralSearchAndChildFailure(t *testing.T) {
	actor := &models.User{ID: "wo-unicode", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&models.Product{ID: "wo-unicode-product", VendorID: actor.ID}))
	for i, name := range []string{"ΟΣ", "ПРИВЕТ", "İSTANBUL", "100%_"} {
		id := fmt.Sprintf("wo-unicode-%d", i)
		require.NoError(t, facades.Orm().Query().Create(&models.Order{ID: id, OrderNumber: id, CustomerName: name, Status: models.OrderStatusPending}))
		require.NoError(t, facades.Orm().Query().Create(&models.OrderItem{ID: id, OrderID: id, ProductID: "wo-unicode-product"}))
	}
	for _, query := range []string{"ος", "привет", "i\u0307stanbul", "100%_"} {
		f, err := services.ParseWorkspaceOrderFilters(map[string]string{"q": query, "per_page": "1"})
		require.NoError(t, err)
		page, err := services.LoadWorkspaceOrders(actor, f)
		require.NoError(t, err)
		require.EqualValues(t, 1, page.OrderPage.FilteredTotal, query)
		require.Len(t, page.Orders, 1)
	}
	// Full readers' count/page queries remain valid when item hydration fails.
	_, err := facades.Orm().Query().Exec("ALTER TABLE order_items RENAME TO wo_unavailable_items")
	require.NoError(t, err)
	defer func() {
		_, e := facades.Orm().Query().Exec("ALTER TABLE wo_unavailable_items RENAME TO order_items")
		require.NoError(t, e)
	}()
	f, err := services.ParseWorkspaceOrderFilters(map[string]string{"selected_order": "wo-unicode-0"})
	require.NoError(t, err)
	result, err := services.LoadWorkspaceOrders(&models.User{ID: "wo-full", Role: models.RoleAdmin}, f)
	require.Error(t, err)
	require.Nil(t, result)
}
