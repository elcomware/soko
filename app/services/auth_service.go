package services

import (
	"elcomware/app/actions"
	"elcomware/app/dto"
	"elcomware/app/facades"
	"elcomware/app/mappers"
	"elcomware/app/models"
)

type AuthService struct {
	authenticateAction *actions.AuthenticateUserAction
	registerAction     *actions.RegisterUserAction
	userMapper         *mappers.UserMapper
}

func NewAuthService() *AuthService {
	return &AuthService{
		authenticateAction: actions.NewAuthenticateUserAction(),
		registerAction:     actions.NewRegisterUserAction(),
		userMapper:         mappers.NewUserMapper(),
	}
}

func (s *AuthService) Login(input dto.LoginInput) (dto.UserOutput, error) {
	return s.authenticateAction.Execute(input)
}

func (s *AuthService) Register(input dto.RegisterInput) (dto.UserOutput, error) {
	return s.registerAction.Execute(input)
}

func (s *AuthService) GetUserByID(id string) (*dto.UserOutput, error) {
	var user models.User
	if err := facades.Orm().Query().Where("id = ?", id).First(&user); err != nil {
		return nil, err
	}
	dto := s.userMapper.ToDTO(user)
	return &dto, nil
}
