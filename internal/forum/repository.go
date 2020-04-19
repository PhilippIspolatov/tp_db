package forum

import "github.com/ifo16u375/tp_db/internal/models"

type Repository interface{
	Insert(forum *models.Forum) error
	Select(slug string) (*models.Forum, error)
}

