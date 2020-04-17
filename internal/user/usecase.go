package user

import "github.com/ifo16u375/tp_db/internal/models"

type Usecase interface {
	CreateUser(user *models.User) ([]*models.User, error)
	UpdateUser(user *models.User) error
	GetUser(nickname string) (*models.User, error)
}
