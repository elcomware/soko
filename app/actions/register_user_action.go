package actions

import (
	"errors"
	"fmt"
	"time"

	"elcomware/app/dto"
	"elcomware/app/facades"
	"elcomware/app/mappers"
	"elcomware/app/models"
)

type RegisterUserAction struct {
	userMapper *mappers.UserMapper
}

func NewRegisterUserAction() *RegisterUserAction {
	return &RegisterUserAction{
		userMapper: mappers.NewUserMapper(),
	}
}

func (a *RegisterUserAction) Execute(input dto.RegisterInput) (dto.UserOutput, error) {
	count, _ := facades.Orm().Query().Model(&models.User{}).Where("email = ?", input.Email).Count()
	if count > 0 {
		return dto.UserOutput{}, errors.New("email address is already registered")
	}

	hashedPassword, err := facades.Hash().Make(input.Password)
	if err != nil {
		hashedPassword = input.Password
	}

	role := models.RoleBuyer
	if input.Role == "merchant" || input.Role == "supplier" {
		role = models.RoleMerchant
	} else if input.Role == "admin" {
		role = models.RoleAdmin
	}

	accountType := input.AccountType
	if accountType == "" {
		if role == models.RoleMerchant {
			accountType = "supplier"
		} else {
			accountType = "individual"
		}
	}

	user := models.User{
		ID:          fmt.Sprintf("usr_%d", time.Now().UnixNano()),
		Name:        input.Name,
		Email:       input.Email,
		Password:    hashedPassword,
		Role:        role,
		CompanyName: input.CompanyName,
		Country:     input.Country,
		AccountType: accountType,
		Verified:    false,
		Avatar:      "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=200&q=80",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := facades.Orm().Query().Create(&user); err != nil {
		return dto.UserOutput{}, err
	}

	return a.userMapper.ToDTO(user), nil
}
