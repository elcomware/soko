package actions

import (
	"elcomware/app/facades"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"errors"
	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/database/orm"
	"math"
	"time"
)

type SubmitRfqQuoteAction struct{}

func NewSubmitRfqQuoteAction() *SubmitRfqQuoteAction { return &SubmitRfqQuoteAction{} }
func (a *SubmitRfqQuoteAction) ExecuteFromRequest(r *requests.SubmitRfqQuoteRequest, u *models.User) (*models.RfqQuote, error) {
	if !u.CanManageCatalog() {
		return nil, errors.New("only suppliers and administrators may quote")
	}
	if r.UnitPrice <= 0 || r.UnitPrice > 1e12 || math.IsNaN(r.UnitPrice) || math.IsInf(r.UnitPrice, 0) || r.LeadTimeDays < 0 || r.LeadTimeDays > 3650 || len(r.Notes) > 5000 {
		return nil, errors.New("provide a positive price and a valid lead time")
	}
	var quote models.RfqQuote
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		var rfq models.Rfq
		if err := tx.Where("id = ?", r.RfqID).First(&rfq); err != nil {
			return err
		}
		if r.UnitPrice*float64(rfq.Quantity) > 1e12 {
			return errors.New("quotation amount exceeds the supported limit")
		}
		if rfq.ID == "" || rfq.Status != models.RfqStatusOpen || rfq.BuyerID == u.ID {
			return errors.New("request is unavailable for quoting")
		}
		if err := tx.Where("rfq_id = ? AND vendor_id = ?", rfq.ID, u.ID).First(&quote); err != nil {
			return err
		}
		if quote.ID == "" {
			quote.ID = "quote_" + uuid.NewString()
			quote.CreatedAt = time.Now()
		}
		quote.RfqID = rfq.ID
		quote.VendorID = u.ID
		quote.VendorName = u.CompanyName
		if quote.VendorName == "" {
			quote.VendorName = u.Name
		}
		quote.VendorCountry = u.Country
		quote.UnitPrice = r.UnitPrice
		quote.TotalAmount = math.Round(r.UnitPrice*float64(rfq.Quantity)*100) / 100
		quote.LeadTimeDays = r.LeadTimeDays
		quote.Notes = r.Notes
		quote.Status = "pending"
		quote.UpdatedAt = time.Now()
		return tx.Save(&quote)
	})
	if err != nil {
		return nil, err
	}
	return &quote, nil
}
