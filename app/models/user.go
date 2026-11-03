package models

import "time"

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleMerchant UserRole = "merchant"
	RoleBuyer    UserRole = "buyer"
)

type User struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Role        UserRole  `json:"role"`
	CompanyName string    `json:"company_name"`
	Country     string    `json:"country"`
	AccountType string    `json:"account_type"` // enterprise, supplier, individual
	Verified    bool      `json:"verified"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
}

