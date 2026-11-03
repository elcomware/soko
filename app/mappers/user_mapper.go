package mappers

import (
	"elcomware/app/dto"
	"elcomware/app/models"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) ToDTO(u models.User) dto.UserOutput {
	return dto.UserOutput{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Role:        string(u.Role),
		CompanyName: u.CompanyName,
		Country:     u.Country,
		AccountType: u.AccountType,
		Verified:    u.Verified,
		Avatar:      u.Avatar,
		CreatedAt:   u.CreatedAt,
	}
}
