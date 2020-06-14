package post

import "github.com/PhilippIspolatov/tp_db/internal/models"

type Usecase interface {
	InsertPosts(posts []*models.Post, thread uint64, forum string) error
	CheckPosts(posts []*models.Post, thread uint64) error
	SelectPost(id uint64) (*models.Post, error)
	UpdatePost(post *models.Post) error
	SelectSortesPosts(threadId uint64, sortType string, desc bool, since uint64, limit uint64) ([]*models.Post, error)
}