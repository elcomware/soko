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

func TestWorkspaceRfqsScopeQuotesAndPaging(t *testing.T) {
	actor := &models.User{ID: "wr-vendor", Role: models.RoleMerchant}
	baseFilters, err := services.ParseWorkspaceRfqFilters(nil)
	require.NoError(t, err)
	baseline, err := services.LoadWorkspaceRfqs(actor, baseFilters)
	require.NoError(t, err)
	stamp := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, r := range []models.Rfq{{ID: "wr-a", RfqNumber: "WR-A", Title: "wr-scope ÉCOLE", Status: models.RfqStatusOpen, CreatedAt: stamp}, {ID: "wr-b", RfqNumber: "WR-B", Title: "wr-scope My history", Status: models.RfqStatusClosed, CreatedAt: stamp}, {ID: "wr-hidden", RfqNumber: "WR-HIDDEN", Title: "Hidden only", Status: models.RfqStatusClosed, CreatedAt: stamp}} {
		require.NoError(t, facades.Orm().Query().Create(&r))
	}
	for _, q := range []models.RfqQuote{{ID: "wr-qa", RfqID: "wr-a", VendorID: actor.ID, Status: "pending", UpdatedAt: stamp}, {ID: "wr-foreign", RfqID: "wr-a", VendorID: "other", Notes: "private"}, {ID: "wr-qb", RfqID: "wr-b", VendorID: actor.ID}, {ID: "wr-qhidden", RfqID: "wr-hidden", VendorID: "other"}} {
		require.NoError(t, facades.Orm().Query().Create(&q))
	}
	f, err := services.ParseWorkspaceRfqFilters(map[string]string{"q": "wr-scope", "per_page": "1", "selected_rfq": "wr-a"})
	require.NoError(t, err)
	page, err := services.LoadWorkspaceRfqs(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, baseline.RfqPage.Total+2, page.RfqPage.Total)
	require.Equal(t, "wr-b", page.Rfqs[0].ID)
	require.Equal(t, "supplier", page.RfqVisibility)
	require.Len(t, page.SelectedRfq.Quotes, 1)
	require.Equal(t, "wr-qa", page.SelectedRfq.Quotes[0].ID)
	require.Equal(t, stamp, page.SelectedRfq.Quotes[0].UpdatedAt)
	require.True(t, page.RecordCapabilities["rfqs"]["wr-a"]["quote"])
	require.False(t, page.RecordCapabilities["quotes"]["wr-qa"]["accept"])
	_, leak := page.RecordCapabilities["quotes"]["wr-foreign"]
	require.False(t, leak)
	f.Query = "école"
	f.Page = 999
	f.SelectedRfq = "wr-hidden"
	page, err = services.LoadWorkspaceRfqs(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, 1, page.RfqPage.FilteredTotal)
	require.Equal(t, 1, page.RfqPage.Page)
	require.Nil(t, page.SelectedRfq)
	require.Equal(t, "unavailable", page.SelectedRfqStatus)
	f.Query = "Hidden only"
	page, err = services.LoadWorkspaceRfqs(actor, f)
	require.NoError(t, err)
	require.Zero(t, page.RfqPage.FilteredTotal)
	f.Query = "wr-scope"
	f.Status = "closed"
	page, err = services.LoadWorkspaceRfqs(actor, f)
	require.NoError(t, err)
	require.EqualValues(t, 1, page.RfqPage.FilteredTotal)
	f.Status = ""
	f.Sort = "oldest"
	f.Page = 1
	page, err = services.LoadWorkspaceRfqs(actor, f)
	require.NoError(t, err)
	require.Equal(t, "wr-a", page.Rfqs[0].ID)
	f.SelectedRfq = "wr-a"
	page, err = services.LoadWorkspaceRfqs(&models.User{ID: "wr-admin", Role: models.RoleAdmin}, f)
	require.NoError(t, err)
	require.Equal(t, "full", page.RfqVisibility)
	require.Len(t, page.SelectedRfq.Quotes, 2)
	require.True(t, page.RecordCapabilities["quotes"]["wr-qa"]["accept"])
}
func TestWorkspaceRfqFiltersAndChildError(t *testing.T) {
	for _, input := range []map[string]string{{"status": "pending"}, {"sort": "price"}, {"page": "0"}, {"per_page": "101"}, {"q": "bad\nquery"}} {
		_, err := services.ParseWorkspaceRfqFilters(input)
		require.Error(t, err)
	}
	f, err := services.ParseWorkspaceRfqFilters(nil)
	require.NoError(t, err)
	_, err = services.LoadWorkspaceRfqs(nil, f)
	require.Error(t, err)
	require.NoError(t, facades.Orm().Query().Create(&models.Rfq{ID: "wr-error", RfqNumber: "WR-ERROR", Status: models.RfqStatusOpen}))
	_, err = facades.Orm().Query().Exec("ALTER TABLE rfq_quotes RENAME TO wr_unavailable_quotes")
	require.NoError(t, err)
	defer func() {
		_, e := facades.Orm().Query().Exec("ALTER TABLE wr_unavailable_quotes RENAME TO rfq_quotes")
		require.NoError(t, e)
	}()
	result, err := services.LoadWorkspaceRfqs(&models.User{ID: "wr-error-admin", Role: models.RoleAdmin}, f)
	require.Error(t, err)
	require.Nil(t, result)
}

func TestWorkspaceRfqsUnicodeQuoteOnlyAndStateRefresh(t *testing.T) {
	actor := &models.User{ID: "wr-quote-only", Role: models.RoleMerchant, PermissionsLoaded: true, PermissionGrants: []models.PermissionGrant{{Permission: "workspace.access", Scope: "own"}, {Permission: "rfqs.quote", Scope: "own"}}}
	require.NoError(t, facades.Orm().Query().Create(&models.AccessProfile{UserID: actor.ID, GrantsJSON: `[{"permission":"workspace.access","scope":"own"},{"permission":"rfqs.quote","scope":"own"}]`, Version: 1}))
	for i, title := range []string{"wr-u ΟΣ", "wr-u ПРИВЕТ", "wr-u İSTANBUL", "wr-u 100%_"} {
		id := fmt.Sprintf("wr-unicode-%d", i)
		require.NoError(t, facades.Orm().Query().Create(&models.Rfq{ID: id, RfqNumber: id, Title: title, Status: models.RfqStatusOpen}))
	}
	for _, query := range []string{"wr-u ος", "wr-u привет", "wr-u i\u0307stanbul", "wr-u 100%_"} {
		f, err := services.ParseWorkspaceRfqFilters(map[string]string{"q": query, "per_page": "1"})
		require.NoError(t, err)
		page, err := services.LoadWorkspaceRfqs(actor, f)
		require.NoError(t, err)
		require.EqualValues(t, 1, page.RfqPage.FilteredTotal, query)
	}
	f, err := services.ParseWorkspaceRfqFilters(map[string]string{"selected_rfq": "wr-unicode-0"})
	require.NoError(t, err)
	page, err := services.LoadWorkspaceRfqs(actor, f)
	require.NoError(t, err)
	require.Equal(t, "available", page.SelectedRfqStatus)
	require.True(t, page.RecordCapabilities["rfqs"]["wr-unicode-0"]["quote"])
	_, err = facades.Orm().Query().Model(&models.Rfq{}).Where("id = ?", "wr-unicode-0").Update("status", models.RfqStatusClosed)
	require.NoError(t, err)
	page, err = services.LoadWorkspaceRfqs(actor, f)
	require.NoError(t, err)
	require.Equal(t, "unavailable", page.SelectedRfqStatus)
	_, exists := page.RecordCapabilities["rfqs"]["wr-unicode-0"]
	require.False(t, exists)
	require.NoError(t, facades.Orm().Query().Create(&models.RfqQuote{ID: "wr-retained", RfqID: "wr-unicode-0", VendorID: actor.ID, Status: "rejected"}))
	page, err = services.LoadWorkspaceRfqs(actor, f)
	require.NoError(t, err)
	require.Equal(t, "available", page.SelectedRfqStatus)
	require.False(t, page.RecordCapabilities["rfqs"]["wr-unicode-0"]["quote"])
	require.False(t, page.RecordCapabilities["quotes"]["wr-retained"]["accept"])
}
