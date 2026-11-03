package dto

import "time"

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type RegisterInput struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	CompanyName string `json:"company_name"`
	Country     string `json:"country"`
	AccountType string `json:"account_type"`
}

type UserOutput struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	CompanyName string    `json:"company_name"`
	Country     string    `json:"country"`
	AccountType string    `json:"account_type"`
	Verified    bool      `json:"verified"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
}
