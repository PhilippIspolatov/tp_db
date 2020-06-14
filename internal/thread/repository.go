package thread

import "github.com/PhilippIspolatov/tp_db/internal/models"

type Repository interface {
	Insert(thread *models.Thread) error
	SelectBySlug(slug string) (*models.Thread, error)
	SelectById(id uint64) (*models.Thread, error)
	Select(slug string, limit uint64, since string, desc bool) ([]*models.Thread, error)
	UpdateById(thread *models.Thread) error
	UpdateBySlug(thread *models.Thread) error
}
