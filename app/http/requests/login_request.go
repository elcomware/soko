package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type LoginRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	Remember bool   `form:"remember" json:"remember"`
}

func (r *LoginRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *LoginRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"email":    "required|email",
		"password": "required|min_len:6",
	}
}

func (r *LoginRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"email.required":    "Email address is required.",
		"email.email":       "Please provide a valid email address.",
		"password.required": "Password is required.",
		"password.min_len":  "Password must be at least 6 characters.",
	}
}
