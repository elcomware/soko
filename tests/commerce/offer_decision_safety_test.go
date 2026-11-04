package commerce_test

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"elcomware/app/services"
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestOfferDecisionRejectsOrphanAndAnonymous(t *testing.T) {
	require.NoError(t, facades.Orm().Query().Create(&models.Offer{ID: "offer-safety-decision-orphan", ProductID: "nonexistent-product", Status: models.OfferStatusPending}))
	require.Error(t, services.DecideOffer("offer-safety-decision-orphan", "accepted", &models.User{ID: "offer-safety-decision-admin", Role: models.RoleAdmin}))
	require.NotPanics(t, func() { require.Error(t, services.DecideOffer("offer-safety-decision-orphan", "accepted", nil)) })
	var offer models.Offer
	require.NoError(t, facades.Orm().Query().Where("id = ?", "offer-safety-decision-orphan").First(&offer))
	require.Equal(t, models.OfferStatusPending, offer.Status)
}
func TestOfferDecisionConcurrentOneWinner(t *testing.T) {
	actor := &models.User{ID: "offer-safety-decision-vendor", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&models.Product{ID: "offer-safety-decision-product", VendorID: actor.ID}))
	require.NoError(t, facades.Orm().Query().Create(&models.Offer{ID: "offer-safety-decision-race", ProductID: "offer-safety-decision-product", Status: models.OfferStatusPending}))
	require.Error(t, services.DecideOffer("offer-safety-decision-race", "accepted", &models.User{ID: "foreign-vendor", Role: models.RoleMerchant}))
	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make(chan error, 2)
	for _, status := range []string{"accepted", "declined"} {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			<-start
			results <- services.DecideOffer("offer-safety-decision-race", s, actor)
		}(status)
	}
	close(start)
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		}
	}
	require.Equal(t, 1, success)
}

func TestOfferDecisionChecksOwnershipAfterCompetingTransfer(t *testing.T) {
	actor := &models.User{ID: "offer-safety-decision-old-owner", Role: models.RoleMerchant}
	require.NoError(t, facades.Orm().Query().Create(&models.Product{ID: "offer-safety-decision-transfer-product", VendorID: actor.ID}))
	require.NoError(t, facades.Orm().Query().Create(&models.Offer{ID: "offer-safety-decision-transfer", ProductID: "offer-safety-decision-transfer-product", Status: models.OfferStatusPending}))
	result := make(chan error, 1)
	started := make(chan struct{})
	require.NoError(t, facades.Orm().Transaction(func(tx orm.Query) error {
		if _, err := tx.Model(&models.Product{}).Where("id = ?", "offer-safety-decision-transfer-product").Update("vendor_id", "offer-safety-decision-new-owner"); err != nil {
			return err
		}
		go func() {
			close(started)
			result <- services.DecideOffer("offer-safety-decision-transfer", "accepted", actor)
		}()
		<-started
		return nil
	}))
	require.Error(t, <-result)
	var offer models.Offer
	require.NoError(t, facades.Orm().Query().Where("id = ?", "offer-safety-decision-transfer").First(&offer))
	require.Equal(t, models.OfferStatusPending, offer.Status)
}
