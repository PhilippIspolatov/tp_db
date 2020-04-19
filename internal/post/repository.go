package post

import "github.com/ifo16u375/tp_db/internal/models"

type Repository interface {
	Insert(posts []*models.Post, thread uint64, forum string) error
	CheckPostsByParentAndAuthor(posts []*models.Post, id uint64) error
	SelectPost(id uint64) (*models.Post, error)
	UpdatePost(post *models.Post) error
}
