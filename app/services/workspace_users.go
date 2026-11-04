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

type WorkspaceUserFilters struct {
	Query        string `json:"q"`
	Role         string `json:"role"`
	Verification string `json:"verification"`
	Sort         string `json:"sort"`
	Page         int    `json:"-"`
	PerPage      int    `json:"-"`
	SelectedUser string `json:"-"`
}
type WorkspaceUserPageMeta struct {
	Page          int                  `json:"page"`
	PerPage       int                  `json:"per_page"`
	Total         int64                `json:"total"`
	FilteredTotal int64                `json:"filtered_total"`
	LastPage      int                  `json:"last_page"`
	From          int                  `json:"from"`
	To            int                  `json:"to"`
	Filters       WorkspaceUserFilters `json:"filters"`
}
type WorkspaceUserPage struct {
	Users              []models.User                         `json:"users"`
	UserPage           WorkspaceUserPageMeta                 `json:"userPage"`
	SelectedUser       *models.User                          `json:"selectedUser"`
	SelectedUserStatus string                                `json:"selectedUserStatus"`
	RecordCapabilities map[string]map[string]map[string]bool `json:"recordCapabilities"`
}

func ParseWorkspaceUserFilters(input map[string]string) (WorkspaceUserFilters, error) {
	f := WorkspaceUserFilters{Query: strings.TrimSpace(input["q"]), Role: input["role"], Verification: input["verification"], Sort: input["sort"], Page: 1, PerPage: 25, SelectedUser: strings.TrimSpace(input["selected_user"])}
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
	return f, validateWorkspaceUserFilters(f)
}
func validateWorkspaceUserFilters(f WorkspaceUserFilters) error {
	if f.Page < 1 || f.Page > 1000000 || f.PerPage < 1 || f.PerPage > 100 {
		return errors.New("invalid pagination")
	}
	if f.Sort != "newest" && f.Sort != "oldest" {
		return errors.New("invalid user sort")
	}
	switch f.Role {
	case "", "admin", "merchant", "buyer":
	default:
		return errors.New("invalid user role")
	}
	if f.Verification != "" && f.Verification != "verified" && f.Verification != "unverified" {
		return errors.New("invalid verification filter")
	}
	if !utf8.ValidString(f.Query) || utf8.RuneCountInString(f.Query) > 120 || !utf8.ValidString(f.SelectedUser) || utf8.RuneCountInString(f.SelectedUser) > 64 || strings.ContainsAny(f.Query+f.SelectedUser, "\x00\r\n") {
		return errors.New("invalid user search or selection")
	}
	return nil
}

// LoadWorkspaceUsers keeps the account loader unchanged. All counts and details
// share one read snapshot. Account directory access always requires users.read/all.
func LoadWorkspaceUsers(actor *models.User, f WorkspaceUserFilters) (*WorkspaceUserPage, error) {
	if actor == nil || !actor.CanAccessWorkspace() || !actor.Can("users.read", "") {
		return nil, errors.New("user workspace access denied")
	}
	if err := validateWorkspaceUserFilters(f); err != nil {
		return nil, err
	}
	out := &WorkspaceUserPage{Users: []models.User{}, SelectedUserStatus: "none", UserPage: WorkspaceUserPageMeta{Page: f.Page, PerPage: f.PerPage, Filters: f}, RecordCapabilities: map[string]map[string]map[string]bool{"users": {}}}
	err := facades.Orm().Transaction(func(tx orm.Query) error {
		scoped := func() orm.Query { return tx.Model(&models.User{}) }
		safeColumns := []string{"id", "name", "email", "role", "company_name", "country", "account_type", "verified", "avatar", "created_at", "updated_at"}
		filtered := func() orm.Query {
			q := scoped()
			if f.Role != "" {
				q = q.Where("role = ?", f.Role)
			}
			if f.Verification != "" {
				q = q.Where("verified = ?", f.Verification == "verified")
			}
			if strings.EqualFold(f.Query, "verified") || strings.EqualFold(f.Query, "unverified") {
				q = q.Where("verified = ?", strings.EqualFold(f.Query, "verified"))
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
		out.UserPage.Total, err = scoped().Count()
		if err != nil {
			return err
		}
		var ids []string
		if f.Query == "" || strings.EqualFold(f.Query, "verified") || strings.EqualFold(f.Query, "unverified") {
			out.UserPage.FilteredTotal, err = filtered().Count()
		} else {
			ids, out.UserPage.FilteredTotal, err = scanWorkspaceUserMatches(ordered(filtered()), f)
		}
		if err != nil {
			return err
		}
		out.UserPage.LastPage = max(1, int((out.UserPage.FilteredTotal+int64(f.PerPage)-1)/int64(f.PerPage)))
		out.UserPage.Page = min(f.Page, out.UserPage.LastPage)
		out.UserPage.Filters.Page = out.UserPage.Page
		offset := (out.UserPage.Page - 1) * f.PerPage
		if f.Query == "" || strings.EqualFold(f.Query, "verified") || strings.EqualFold(f.Query, "unverified") {
			err = ordered(filtered()).Select(strings.Join(safeColumns, ",")).Offset(offset).Limit(f.PerPage).Get(&out.Users)
		} else if len(ids) > 0 {
			err = ordered(filtered()).Select(strings.Join(safeColumns, ",")).Where("id IN (?)", ids).Limit(f.PerPage).Get(&out.Users)
		}
		if err != nil {
			return err
		}
		if len(out.Users) > 0 {
			out.UserPage.From = offset + 1
			out.UserPage.To = offset + len(out.Users)
		}
		if f.SelectedUser != "" {
			out.SelectedUserStatus = "unavailable"
			var selected models.User
			if err := scoped().Select(strings.Join(safeColumns, ",")).Where("id = ?", f.SelectedUser).First(&selected); err != nil {
				return err
			}
			if selected.ID != "" {
				out.SelectedUser = &selected
				out.SelectedUserStatus = "available"
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
func scanWorkspaceUserMatches(q orm.Query, f WorkspaceUserFilters) ([]string, int64, error) {
	lower := cases.Lower(language.Und)
	needle := lower.String(f.Query)
	requested := []string{}
	last := []string{}
	var matched int64
	var scanErr error
	cursor := q.Select("id", "name", "email", "company_name", "country", "role").Cursor()
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
		var user models.User
		if err := row.Scan(&user); err != nil {
			scanErr = err
			continue
		}
		found := false
		for _, value := range []string{user.Name, user.Email, user.CompanyName, user.Country, string(user.Role)} {
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
		last = append(last, user.ID)
		page := int(matched/int64(f.PerPage)) + 1
		if page == f.Page {
			requested = append(requested, user.ID)
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
