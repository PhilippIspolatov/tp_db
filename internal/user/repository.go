package user

import "github.com/ifo16u375/tp_db/internal/models"

type Repository interface {
	Insert(user *models.User) error
	Update(user *models.User) error
	SelectByEmail(email string) (*models.User, error)
	SelectByNickname(nickname string) (*models.User, error)
	SelectAllUsers(slug string, limit uint64, since string, desc bool) ([]*models.User, error)
}