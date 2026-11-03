package actions

import (
	"errors"

	"elcomware/app/dto"
	"elcomware/app/facades"
	"elcomware/app/mappers"
	"elcomware/app/models"
)

type AuthenticateUserAction struct {
	userMapper *mappers.UserMapper
}

func NewAuthenticateUserAction() *AuthenticateUserAction {
	return &AuthenticateUserAction{
		userMapper: mappers.NewUserMapper(),
	}
}

func (a *AuthenticateUserAction) Execute(input dto.LoginInput) (dto.UserOutput, error) {
	var user models.User
	if err := facades.Orm().Query().Where("email = ?", input.Email).First(&user); err != nil {
		return dto.UserOutput{}, errors.New("invalid email or password")
	}

	// Verify hashed password
	if !facades.Hash().Check(input.Password, user.Password) {
		// Fallback check for plain dev passwords
		if user.Password != "plain:"+input.Password && user.Password != input.Password {
			return dto.UserOutput{}, errors.New("invalid email or password")
		}
	}

	return a.userMapper.ToDTO(user), nil
}
