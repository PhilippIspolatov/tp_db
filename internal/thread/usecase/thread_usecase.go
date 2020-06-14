package usecase

import (
	"github.com/PhilippIspolatov/tp_db/internal/models"
	"github.com/PhilippIspolatov/tp_db/internal/thread"
)

type ThreadUsecase struct {
	threadRepo thread.Repository
}

func NewThreadUsecase(tr thread.Repository) thread.Usecase {
	return &ThreadUsecase{
		threadRepo: tr,
	}
}

func (tUC *ThreadUsecase) InsertThread(thread *models.Thread) error {
	if err := tUC.threadRepo.Insert(thread); err != nil {
		return err
	}
	return nil
}

func (tUC *ThreadUsecase) SelectThreadBySlug(slug string) (*models.Thread, error) {
	t, err := tUC.threadRepo.SelectBySlug(slug)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (tUC *ThreadUsecase) SelectThreadById(id uint64) (*models.Thread, error) {
	t, err := tUC.threadRepo.SelectById(id)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (tUC *ThreadUsecase) SelectThreads(slug string, limit uint64, since string, desc bool) ([]*models.Thread, error) {
	return tUC.threadRepo.Select(slug, limit, since, desc)
}

func (tUC *ThreadUsecase) UpdateThreadById(thread *models.Thread) error {
	if err := tUC.threadRepo.UpdateById(thread); err != nil {
		return err
	}
	return nil
}

func (tUC *ThreadUsecase) UpdateThreadBySlug(thread *models.Thread) error {
	if err := tUC.threadRepo.UpdateBySlug(thread); err != nil {
		return err
	}
	return nil
}