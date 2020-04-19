package forum

import "github.com/ifo16u375/tp_db/internal/models"

type Usecase interface {
	CreateForum(forum *models.Forum) error
	GetForum(slug string) (*models.Forum, error)
}
