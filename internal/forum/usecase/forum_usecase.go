package usecase

import (
	"github.com/ifo16u375/tp_db/internal/forum"
	"github.com/ifo16u375/tp_db/internal/models"
)

type ForumUsecase struct {
	forumRepo forum.Repository
}

func NewForumUsecase(fr forum.Repository) forum.Usecase {
	return &ForumUsecase{
		forumRepo:fr,
	}
}

func (fUC *ForumUsecase) CreateForum(forum *models.Forum) error {
	if err := fUC.forumRepo.Insert(forum); err != nil {
		return err
	}
	return nil
}

func (fUC *ForumUsecase) GetForum(slug string) (*models.Forum, error) {
	f, err := fUC.forumRepo.Select(slug)
	if err != nil {
		return nil, err
	}
	return f, nil
}
