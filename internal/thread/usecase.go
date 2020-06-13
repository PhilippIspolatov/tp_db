package thread

import "github.com/ifo16u375/tp_db/internal/models"

type Usecase interface {
	InsertThread(thread *models.Thread) error
	SelectThreadBySlug(slug string) (*models.Thread, error)
	SelectThreadById(id uint64) (*models.Thread, error)
	SelectThreads(slug string, limit uint64, since string, desc bool) ([]*models.Thread, error)
	UpdateThreadById(thread *models.Thread) error
	UpdateThreadBySlug(thread *models.Thread) error
}