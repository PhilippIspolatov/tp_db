package forum

import "github.com/PhilippIspolatov/tp_db/internal/models"

type Repository interface{
	Insert(forum *models.Forum) error
	Select(slug string) (*models.Forum, error)
}

