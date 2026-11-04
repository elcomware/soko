package actions

import (
	"errors"
	"strings"

	"elcomware/app/facades"
	"elcomware/app/models"
)

type AuthenticateUserAction struct{}

func NewAuthenticateUserAction() *AuthenticateUserAction {
	return &AuthenticateUserAction{}
}

func (a *AuthenticateUserAction) Execute(email, password string) (*models.User, error) {
	var user models.User
	if err := facades.Orm().Query().Where("LOWER(email) = ?", strings.ToLower(strings.TrimSpace(email))).First(&user); err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.ID == "" || !facades.Hash().Check(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}
