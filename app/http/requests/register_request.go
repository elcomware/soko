package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type RegisterRequest struct {
	Name        string `form:"name" json:"name"`
	Email       string `form:"email" json:"email"`
	Password    string `form:"password" json:"password"`
	CompanyName string `form:"company_name" json:"company_name"`
	Country     string `form:"country" json:"country"`
	AccountType string `form:"account_type" json:"account_type"`
}

func (r *RegisterRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *RegisterRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"name":     "required|min_len:2",
		"email":    "required|email",
		"password": "required|min_len:6",
	}
}

func (r *RegisterRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":     "Full name or business contact name is required.",
		"email.required":    "Email address is required.",
		"email.email":       "Please provide a valid business email address.",
		"password.required": "Password is required.",
		"password.min_len":  "Password must be at least 6 characters.",
	}
}

