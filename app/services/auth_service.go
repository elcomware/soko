package services

import (
	"elcomware/app/actions"
	"elcomware/app/models"
	"elcomware/app/queries"
)

type AuthService struct {
	userQueries        *queries.UserQueries
	authenticateAction *actions.AuthenticateUserAction
	registerAction     *actions.RegisterUserAction
}

func NewAuthService() *AuthService {
	return &AuthService{
		userQueries:        queries.NewUserQueries(),
		authenticateAction: actions.NewAuthenticateUserAction(),
		registerAction:     actions.NewRegisterUserAction(),
	}
}

func (s *AuthService) Authenticate(email, password string) (*models.User, error) {
	return s.authenticateAction.Execute(email, password)
}

func (s *AuthService) Register(user *models.User) (*models.User, error) {
	return s.registerAction.Execute(user)
}

func (s *AuthService) GetUserByID(id string) (*models.User, error) {
	return s.userQueries.FindByID(id)
}

func (s *AuthService) GetAllUsers() ([]models.User, error) {
	return s.userQueries.FindAll()
}
