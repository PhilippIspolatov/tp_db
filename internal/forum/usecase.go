package forum

import "github.com/PhilippIspolatov/tp_db/internal/models"

type Usecase interface {
	CreateForum(forum *models.Forum) error
	GetForum(slug string) (*models.Forum, error)
}
