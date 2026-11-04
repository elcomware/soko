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
	if !u.CanManageCatalog() {
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
	if offer.ID == "" || !u.CanManageProduct(&p) {
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
	if u == nil {
		return errors.New("sign in required")
	}
	return facades.Orm().Transaction(func(tx orm.Query) error {
		var rfq models.Rfq
		if err := tx.Where("id = ?", id).First(&rfq); err != nil {
			return err
		}
		if rfq.ID == "" || (rfq.BuyerID != u.ID && u.Role != models.RoleAdmin) {
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
