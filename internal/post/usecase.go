package post

import "github.com/ifo16u375/tp_db/internal/models"

type Usecase interface {
	InsertPosts(posts []*models.Post, thread uint64, forum string) error
	CheckPosts(posts []*models.Post, thread uint64) error
	SelectPost(id uint64) (*models.Post, error)
	UpdatePost(post *models.Post) error
}