package actions

import (
	"errors"
	"fmt"
	"time"

	"elcomware/app/facades"
	"elcomware/app/models"
)

type RegisterUserAction struct{}

func NewRegisterUserAction() *RegisterUserAction {
	return &RegisterUserAction{}
}

func (a *RegisterUserAction) Execute(user *models.User) (*models.User, error) {
	count, _ := facades.Orm().Query().Model(&models.User{}).Where("email = ?", user.Email).Count()
	if count > 0 {
		return nil, errors.New("email address is already registered")
	}

	hashedPassword, err := facades.Hash().Make(user.Password)
	if err == nil {
		user.Password = hashedPassword
	}

	if user.ID == "" {
		user.ID = fmt.Sprintf("usr_%d", time.Now().UnixNano())
	}
	if user.Avatar == "" {
		user.Avatar = "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=200&q=80"
	}
	user.Verified = false

	if err := facades.Orm().Query().Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

