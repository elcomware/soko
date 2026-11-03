package models

import "time"

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleMerchant UserRole = "merchant"
	RoleBuyer    UserRole = "buyer"
)

type User struct {
	ID          string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	Email       string    `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Password    string    `gorm:"type:varchar(255)" json:"-"`
	Role        UserRole  `gorm:"type:varchar(32);default:buyer" json:"role"`
	CompanyName string    `gorm:"type:varchar(255)" json:"company_name"`
	Country     string    `gorm:"type:varchar(64)" json:"country"`
	AccountType string    `gorm:"type:varchar(32);default:individual" json:"account_type"`
	Verified    bool      `gorm:"default:false" json:"verified"`
	Avatar      string    `gorm:"type:text" json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
