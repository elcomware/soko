package queries

import (
	"elcomware/app/facades"
	"elcomware/app/models"
	"errors"
)

type UserQueries struct{}

func NewUserQueries() *UserQueries {
	return &UserQueries{}
}

func (q *UserQueries) FindByID(id string) (*models.User, error) {
	var user models.User
	if err := facades.Orm().Query().Where("id = ?", id).First(&user); err != nil {
		return nil, err
	}
	if user.ID == "" {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (q *UserQueries) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := facades.Orm().Query().Where("email = ?", email).First(&user); err != nil {
		return nil, err
	}
	if user.ID == "" {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (q *UserQueries) FindAll() ([]models.User, error) {
	var users []models.User
	err := facades.Orm().Query().OrderByDesc("created_at").Get(&users)
	return users, err
}

func (q *UserQueries) Count() (int64, error) {
	return facades.Orm().Query().Model(&models.User{}).Count()
}
