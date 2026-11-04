package services

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"errors"
	"github.com/goravel/framework/contracts/database/orm"
	"time"
)

func DecideOffer(id, status string, u *models.User) error {
	if status != "accepted" && status != "declined" {
		return errors.New("choose accepted or declined")
	}
	if !u.HasPermission("offers.decide") {
		return errors.New("access denied")
	}
	var offer models.Offer
	if err := facades.Orm().Query().Where("id = ?", id).First(&offer); err != nil {
		return err
	}
	var p models.Product
	if err := facades.Orm().Query().Where("id = ?", offer.ProductID).First(&p); err != nil {
		return err
	}
	if offer.ID == "" || !u.Can("offers.decide", p.VendorID) {
		return errors.New("offer not found or access denied")
	}
	result, err := facades.Orm().Query().Model(&models.Offer{}).Where("id = ? AND status = ?", id, models.OfferStatusPending).Update(map[string]any{"status": status, "updated_at": time.Now()})
	if err != nil {
		return err
	}
	if result.RowsAffected != 1 {
		return errors.New("offer has already been reviewed")
	}
	return nil
}
func ResolveRfq(id, quoteID string, closeOnly bool, u *models.User) error {
	return resolveRfq(id, quoteID, closeOnly, u, nil)
}

// ResolveRfqWithQuoteVersion awards only the exact quotation reviewed by the
// administrator. A stale or absent snapshot requires explicit refresh/review.
func ResolveRfqWithQuoteVersion(id, quoteID, version string, u *models.User) error {
	expected, err := time.Parse(time.RFC3339Nano, version)
	if err != nil || expected.IsZero() {
		return errors.New("Refresh the quotation details and review the terms before accepting.")
	}
	return resolveRfq(id, quoteID, false, u, &expected)
}

func resolveRfq(id, quoteID string, closeOnly bool, u *models.User, expected *time.Time) error {
	if u == nil {
		return errors.New("sign in required")
	}
	return facades.Orm().Transaction(func(tx orm.Query) error {
		// Lock the parent before any reads. Quote submission takes the same
		// write lock, preventing edits from interleaving with award/closure.
		// A no-op UPDATE is portable to SQLite, unlike SELECT FOR UPDATE.
		if _, err := tx.Exec("UPDATE rfqs SET id = id WHERE id = ?", id); err != nil {
			return err
		}
		var rfq models.Rfq
		if err := tx.Where("id = ?", id).First(&rfq); err != nil {
			return err
		}
		if rfq.ID == "" || !u.Can(map[bool]string{true: "rfqs.close", false: "rfqs.award"}[closeOnly], rfq.BuyerID) {
			return errors.New("request not found or access denied")
		}
		status := models.RfqStatusClosed
		var quote models.RfqQuote
		if !closeOnly {
			if err := tx.Where("id = ? AND rfq_id = ?", quoteID, id).First(&quote); err != nil {
				return err
			}
			if quote.ID == "" || quote.Status != "pending" {
				return errors.New("quotation is not available")
			}
			if expected != nil && !quote.UpdatedAt.Equal(*expected) {
				return errors.New("This quotation changed. Refresh the details and review the updated terms before accepting.")
			}
			status = models.RfqStatusAwarded
		}
		result, err := tx.Model(&models.Rfq{}).Where("id = ? AND status = ?", id, models.RfqStatusOpen).Update(map[string]any{"status": status, "updated_at": time.Now()})
		if err != nil {
			return err
		}
		if result.RowsAffected != 1 {
			return errors.New("request is no longer open")
		}
		if closeOnly {
			return nil
		}
		if _, err := tx.Model(&models.RfqQuote{}).Where("rfq_id = ?", id).Update(map[string]any{"status": "rejected", "updated_at": time.Now()}); err != nil {
			return err
		}
		_, err = tx.Model(&models.RfqQuote{}).Where("id = ?", quoteID).Update(map[string]any{"status": "accepted", "updated_at": time.Now()})
		return err
	})
}
