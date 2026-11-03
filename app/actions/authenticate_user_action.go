package actions

import (
	"errors"

	"elcomware/app/facades"
	"elcomware/app/models"
)

type AuthenticateUserAction struct{}

func NewAuthenticateUserAction() *AuthenticateUserAction {
	return &AuthenticateUserAction{}
}

func (a *AuthenticateUserAction) Execute(email, password string) (*models.User, error) {
	var user models.User
	if err := facades.Orm().Query().Where("email = ?", email).First(&user); err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !facades.Hash().Check(password, user.Password) {
		if user.Password != "plain:"+password && user.Password != password {
			return nil, errors.New("invalid email or password")
		}
	}

	return &user, nil
}

