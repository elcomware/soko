package services

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"errors"
	"github.com/goravel/framework/contracts/database/orm"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
	"strings"
	"unicode/utf8"
)

type WorkspaceRfqFilters struct {
	Query       string `json:"q"`
	Status      string `json:"status"`
	Sort        string `json:"sort"`
	Page        int    `json:"-"`
	PerPage     int    `json:"-"`
	SelectedRfq string `json:"-"`
}
type WorkspaceRfqPageMeta struct {
	Page          int                 `json:"page"`
	PerPage       int                 `json:"per_page"`
	Total         int64               `json:"total"`
	FilteredTotal int64               `json:"filtered_total"`
	LastPage      int                 `json:"last_page"`
	From          int                 `json:"from"`
	To            int                 `json:"to"`
	Filters       WorkspaceRfqFilters `json:"filters"`
}
type WorkspaceRfqPage struct {
	Rfqs               []models.Rfq                          `json:"rfqs"`
	RfqPage            WorkspaceRfqPageMeta                  `json:"rfqPage"`
	SelectedRfq        *models.Rfq                           `json:"selectedRfq"`
	SelectedRfqStatus  string                                `json:"selectedRfqStatus"`
	RfqVisibility      string                                `json:"rfqVisibility"`
	RecordCapabilities map[string]map[string]map[string]bool `json:"recordCapabilities"`
}

func ParseWorkspaceRfqFilters(input map[string]string) (WorkspaceRfqFilters, error) {
	f := WorkspaceRfqFilters{Query: strings.TrimSpace(input["q"]), Status: input["status"], Sort: input["sort"], Page: 1, PerPage: 25, SelectedRfq: strings.TrimSpace(input["selected_rfq"])}
	if f.Sort == "" {
		f.Sort = "newest"
	}
	for key, target := range map[string]*int{"page": &f.Page, "per_page": &f.PerPage} {
		if raw := input[key]; raw != "" {
			v, err := strconv.Atoi(raw)
			if err != nil {
				return f, errors.New("invalid " + key)
			}
			*target = v
		}
	}
	return f, validateWorkspaceRfqFilters(f)
}
func validateWorkspaceRfqFilters(f WorkspaceRfqFilters) error {
	if f.Page < 1 || f.Page > 1000000 || f.PerPage < 1 || f.PerPage > 100 {
		return errors.New("invalid pagination")
	}
	if f.Sort != "newest" && f.Sort != "oldest" {
		return errors.New("invalid rfq sort")
	}
	switch f.Status {
	case "", "open", "awarded", "closed":
	default:
		return errors.New("invalid rfq status")
	}
	if !utf8.ValidString(f.Query) || utf8.RuneCountInString(f.Query) > 120 || !utf8.ValidString(f.SelectedRfq) || utf8.RuneCountInString(f.SelectedRfq) > 64 || strings.ContainsAny(f.Query+f.SelectedRfq, "\x00\r\n") {
		return errors.New("invalid rfq search or selection")
	}
	return nil
}

// LoadWorkspaceRfqs keeps the account loader unchanged. All counts and details
// share one read snapshot; supplier views never hydrate other vendors' quotes.
func LoadWorkspaceRfqs(actor *models.User, f WorkspaceRfqFilters) (*WorkspaceRfqPage, error) {
	if actor == nil || !actor.CanAccessWorkspace() || (!actor.HasPermission("rfqs.read") && !actor.HasPermission("rfqs.quote")) {
		return nil, errors.New("rfq workspace access denied")
	}
	if err := validateWorkspaceRfqFilters(f); err != nil {
		return nil, err
	}
	full := actor.Can("rfqs.read", "")
	out := &WorkspaceRfqPage{Rfqs: []models.Rfq{}, SelectedRfqStatus: "none", RfqVisibility: "supplier", RfqPage: WorkspaceRfqPageMeta{Page: f.Page, PerPage: f.PerPage, Filters: f}, RecordCapabilities: map[string]map[string]map[string]bool{"rfqs": {}, "quotes": {}}}
	if full {
		out.RfqVisibility = "full"
	}
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		scoped := func() orm.Query {
			q := tx.Model(&models.Rfq{})
			if !full {
				q = q.Where("(status = ? OR id IN (SELECT rfq_id FROM rfq_quotes WHERE vendor_id = ?))", models.RfqStatusOpen, actor.ID)
			}
			return q
		}
		filtered := func() orm.Query {
			q := scoped()
			if f.Status != "" {
				q = q.Where("status = ?", f.Status)
			}
			return q
		}
		ordered := func(q orm.Query) orm.Query {
			if f.Sort == "oldest" {
				return q.OrderBy("created_at").OrderBy("id")
			}
			return q.OrderByDesc("created_at").OrderByDesc("id")
		}
		var err error
		out.RfqPage.Total, err = scoped().Count()
		if err != nil {
			return err
		}
		var ids []string
		if f.Query == "" {
			out.RfqPage.FilteredTotal, err = filtered().Count()
		} else {
			ids, out.RfqPage.FilteredTotal, err = scanWorkspaceRfqMatches(ordered(filtered()), f)
		}
		if err != nil {
			return err
		}
		out.RfqPage.LastPage = max(1, int((out.RfqPage.FilteredTotal+int64(f.PerPage)-1)/int64(f.PerPage)))
		out.RfqPage.Page = min(f.Page, out.RfqPage.LastPage)
		out.RfqPage.Filters.Page = out.RfqPage.Page
		offset := (out.RfqPage.Page - 1) * f.PerPage
		if f.Query == "" {
			err = ordered(filtered()).Offset(offset).Limit(f.PerPage).Get(&out.Rfqs)
		} else if len(ids) > 0 {
			err = ordered(filtered()).Where("id IN (?)", ids).Limit(f.PerPage).Get(&out.Rfqs)
		}
		if err != nil {
			return err
		}
		if len(out.Rfqs) > 0 {
			out.RfqPage.From = offset + 1
			out.RfqPage.To = offset + len(out.Rfqs)
		}
		if f.SelectedRfq != "" {
			out.SelectedRfqStatus = "unavailable"
			var selected models.Rfq
			if err := scoped().Where("id = ?", f.SelectedRfq).First(&selected); err != nil {
				return err
			}
			if selected.ID != "" {
				out.SelectedRfq = &selected
				out.SelectedRfqStatus = "available"
			}
		}
		records := make([]*models.Rfq, 0, len(out.Rfqs)+1)
		ids = []string{}
		for i := range out.Rfqs {
			records = append(records, &out.Rfqs[i])
			ids = append(ids, out.Rfqs[i].ID)
		}
		if out.SelectedRfq != nil {
			records = append(records, out.SelectedRfq)
			ids = append(ids, out.SelectedRfq.ID)
		}
		if len(ids) == 0 {
			return nil
		}
		quotesQuery := tx.Model(&models.RfqQuote{}).Where("rfq_id IN (?)", ids)
		if !full {
			quotesQuery = quotesQuery.Where("vendor_id = ?", actor.ID)
		}
		var quotes []models.RfqQuote
		if err := quotesQuery.OrderByDesc("created_at").OrderByDesc("id").Get(&quotes); err != nil {
			return err
		}
		byRfq := map[string][]models.RfqQuote{}
		for _, quote := range quotes {
			byRfq[quote.RfqID] = append(byRfq[quote.RfqID], quote)
		}
		for _, rfq := range records {
			rfq.Quotes = byRfq[rfq.ID]
			if rfq.Quotes == nil {
				rfq.Quotes = []models.RfqQuote{}
			}
			open := rfq.Status == models.RfqStatusOpen
			award := open && actor.Can("rfqs.award", rfq.BuyerID)
			out.RecordCapabilities["rfqs"][rfq.ID] = map[string]bool{"quote": open && rfq.BuyerID != actor.ID && actor.HasPermission("rfqs.quote"), "award": award, "close": open && actor.Can("rfqs.close", rfq.BuyerID)}
			for _, quote := range rfq.Quotes {
				out.RecordCapabilities["quotes"][quote.ID] = map[string]bool{"accept": award && quote.Status == "pending" && quote.RfqID == rfq.ID}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
func scanWorkspaceRfqMatches(q orm.Query, f WorkspaceRfqFilters) ([]string, int64, error) {
	lower := cases.Lower(language.Und)
	needle := lower.String(f.Query)
	requested := []string{}
	last := []string{}
	var matched int64
	var scanErr error
	cursor := q.Select("id", "title", "category", "rfq_number", "buyer_company", "buyer_name", "destination_country", "status").Cursor()
	for row := range cursor {
		// Always drain: the ORM cursor's producer cannot be left blocked on send.
		if row.Err() != nil {
			if scanErr == nil {
				scanErr = row.Err()
			}
			continue
		}
		if scanErr != nil {
			continue
		}
		var rfq models.Rfq
		if err := row.Scan(&rfq); err != nil {
			scanErr = err
			continue
		}
		found := false
		for _, value := range []string{rfq.Title, rfq.Category, rfq.RfqNumber, rfq.BuyerCompany, rfq.BuyerName, rfq.DestinationCountry, string(rfq.Status)} {
			if strings.Contains(lower.String(value), needle) {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		if matched%int64(f.PerPage) == 0 {
			last = last[:0]
		}
		last = append(last, rfq.ID)
		page := int(matched/int64(f.PerPage)) + 1
		if page == f.Page {
			requested = append(requested, rfq.ID)
		}
		matched++
	}
	if scanErr != nil {
		return nil, 0, scanErr
	}
	if len(requested) == 0 && matched > 0 {
		return last, matched, nil
	}
	return requested, matched, nil
}
