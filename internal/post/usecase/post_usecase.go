package usecase

import (
	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/post"
)

type PostUsecase struct {
	postRepo post.Repository
}

func NewPostUsecase(pr post.Repository) post.Usecase {
	return &PostUsecase{
		postRepo: pr,
	}
}

func (pUC *PostUsecase) InsertPosts(posts []*models.Post, thread uint64, forum string) error {
	if err := pUC.postRepo.Insert(posts, thread, forum); err != nil {
		return err
	}
	return nil
}

func (pUC *PostUsecase) CheckPosts(posts []*models.Post, thread uint64) error {
	if err := pUC.postRepo.CheckPostsByParentAndAuthor(posts, thread); err != nil {
		return err
	}
	return nil
}

func (pUC *PostUsecase) SelectPost(id uint64) (*models.Post, error) {
	p, err := pUC.postRepo.SelectPost(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (pUC *PostUsecase) UpdatePost(post *models.Post) error {
	if err := pUC.postRepo.UpdatePost(post); err != nil {
		return err
	}
	return nil
}

func (pUC *PostUsecase) SelectSortesPosts(threadId uint64, sortType string, desc bool, since uint64, limit uint64) ([]*models.Post, error) {
	var err error
	posts := []*models.Post{}

	switch sortType {
	case "flat":
		posts, err = pUC.postRepo.GetPostsByFlat(threadId, desc, since, limit)
	case "tree":
		posts, err = pUC.postRepo.GetPostsByTree(threadId, desc, since, limit)
	case "parent_tree":
		posts, err = pUC.postRepo.GetPostsByParentTree(threadId, desc, since, limit)
	default:
		posts, err = pUC.postRepo.GetPostsByFlat(threadId, desc, since, limit)
	}

	if err != nil {
		return nil, err
	}

	return posts, err
}