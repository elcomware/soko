package commerce_test

import (
	"elcomware/app/actions"
	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/services"
	"fmt"
	frameworktesting "github.com/goravel/framework/testing"
	"github.com/stretchr/testify/require"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDashboardQuoteAwardRejectsChangedReviewedTerms(t *testing.T) {
	admin := models.User{ID: "version-award-admin", Email: "version-award-admin@test.invalid", Role: models.RoleAdmin}
	seller := models.User{ID: "version-award-seller", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&admin))
	r := models.Rfq{ID: "version-award-rfq", RfqNumber: "version-award-rfq", BuyerID: "version-award-buyer", Quantity: 10, Status: models.RfqStatusOpen}
	require.NoError(t, facades.Orm().Query().Create(&r))
	submit := actions.NewSubmitRfqQuoteAction()
	q, err := submit.ExecuteFromRequest(&requests.SubmitRfqQuoteRequest{RfqID: r.ID, UnitPrice: 10, LeadTimeDays: 7}, &seller)
	require.NoError(t, err)
	reviewed := q.UpdatedAt.Format(time.RFC3339Nano)
	updated, err := submit.ExecuteFromRequest(&requests.SubmitRfqQuoteRequest{RfqID: r.ID, UnitPrice: 100, LeadTimeDays: 30}, &seller)
	require.NoError(t, err)
	require.Equal(t, q.ID, updated.ID)
	require.NotEqual(t, reviewed, updated.UpdatedAt.Format(time.RFC3339Nano))
	tc := &frameworktesting.TestCase{}
	award := func(version string) {
		body := url.Values{"quote_version": {version}}
		response, err := tc.Http(t).WithHeaders(map[string]string{"Sec-Fetch-Site": "same-origin", "Content-Type": "application/x-www-form-urlencoded"}).WithSession(map[string]any{"user_id": admin.ID}).Post("/dashboard/rfqs/"+r.ID+"/quotes/"+q.ID+"/accept", strings.NewReader(body.Encode()))
		require.NoError(t, err)
		response.AssertStatus(303).AssertHeader("Location", "/dashboard?section=rfqs")
	}
	for _, version := range []string{reviewed, "", "invalid"} {
		award(version)
		var stored models.Rfq
		require.NoError(t, facades.Orm().Query().Where("id = ?", r.ID).First(&stored))
		require.Equal(t, models.RfqStatusOpen, stored.Status, "stale/missing/invalid version must not award")
		var quote models.RfqQuote
		require.NoError(t, facades.Orm().Query().Where("id = ?", q.ID).First(&quote))
		require.Equal(t, "pending", quote.Status)
		require.Equal(t, 100.0, quote.UnitPrice)
		require.Equal(t, 30, quote.LeadTimeDays)
	}
	award(updated.UpdatedAt.Format(time.RFC3339Nano))
	var stored models.Rfq
	require.NoError(t, facades.Orm().Query().Where("id = ?", r.ID).First(&stored))
	require.Equal(t, models.RfqStatusAwarded, stored.Status)
	_, err = submit.ExecuteFromRequest(&requests.SubmitRfqQuoteRequest{RfqID: r.ID, UnitPrice: 200, LeadTimeDays: 60}, &seller)
	require.Error(t, err)
	var accepted models.RfqQuote
	require.NoError(t, facades.Orm().Query().Where("id = ?", q.ID).First(&accepted))
	require.Equal(t, "accepted", accepted.Status)
	require.Equal(t, 100.0, accepted.UnitPrice)
	require.Equal(t, 30, accepted.LeadTimeDays)
}

// Whichever transaction wins, an accepted quote must retain the reviewed terms.
func TestConcurrentQuoteEditAndAwardPreserveReviewedTerms(t *testing.T) {
	for n := 0; n < 8; n++ {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			rfq := models.Rfq{ID: fmt.Sprintf("race-award-%d", n), RfqNumber: fmt.Sprintf("race-award-%d", n), BuyerID: "race-award-buyer", Quantity: 10, Status: models.RfqStatusOpen}
			require.NoError(t, facades.Orm().Query().Create(&rfq))
			seller := &models.User{ID: "race-award-seller", Role: models.RoleMerchant}
			admin := &models.User{ID: "race-award-admin", Role: models.RoleAdmin}
			submit := actions.NewSubmitRfqQuoteAction()
			quote, err := submit.ExecuteFromRequest(&requests.SubmitRfqQuoteRequest{RfqID: rfq.ID, UnitPrice: 10, LeadTimeDays: 7}, seller)
			require.NoError(t, err)
			start := make(chan struct{})
			var wg sync.WaitGroup
			var editErr, awardErr error
			wg.Add(2)
			go func() {
				defer wg.Done()
				<-start
				_, editErr = submit.ExecuteFromRequest(&requests.SubmitRfqQuoteRequest{RfqID: rfq.ID, UnitPrice: 100, LeadTimeDays: 30}, seller)
			}()
			go func() {
				defer wg.Done()
				<-start
				awardErr = services.ResolveRfqWithQuoteVersion(rfq.ID, quote.ID, quote.UpdatedAt.Format(time.RFC3339Nano), admin)
			}()
			close(start)
			wg.Wait()
			require.True(t, editErr == nil || awardErr == nil, "at least one transaction should complete")
			var stored models.RfqQuote
			require.NoError(t, facades.Orm().Query().Where("id = ?", quote.ID).First(&stored))
			var parent models.Rfq
			require.NoError(t, facades.Orm().Query().Where("id = ?", rfq.ID).First(&parent))
			if awardErr == nil {
				require.Error(t, editErr)
				require.Equal(t, models.RfqStatusAwarded, parent.Status)
				require.Equal(t, "accepted", stored.Status)
				require.Equal(t, 10.0, stored.UnitPrice)
				require.Equal(t, 7, stored.LeadTimeDays)
			} else {
				require.NoError(t, editErr)
				require.Equal(t, models.RfqStatusOpen, parent.Status)
				require.Equal(t, "pending", stored.Status)
				require.Equal(t, 100.0, stored.UnitPrice)
				require.Equal(t, 30, stored.LeadTimeDays)
			}
		})
	}
}
