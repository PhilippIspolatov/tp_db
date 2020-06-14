package post

import "github.com/PhilippIspolatov/tp_db/internal/models"

type Repository interface {
	Insert(posts []*models.Post, thread uint64, forum string) error
	CheckPostsByParentAndAuthor(posts []*models.Post, id uint64) error
	SelectPost(id uint64) (*models.Post, error)
	UpdatePost(post *models.Post) error
	GetPostsByFlat(threadId uint64, desc bool, since uint64, limit uint64) ([]*models.Post, error)
	GetPostsByTree(threadId uint64, desc bool, since uint64, limit uint64) ([]*models.Post, error)
	GetPostsByParentTree(threadId uint64, desc bool, since uint64, limit uint64) ([]*models.Post, error)
}
