package usecase

import (
	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/vote"
)

type VoteUsecase struct {
	voteRepo vote.Repository
}

func NewVoteUsecase(vr vote.Repository) vote.Usecase {
	return &VoteUsecase{
		voteRepo:vr,
	}
}

func (vUC *VoteUsecase) InsertVote(vote *models.Vote) error {
	if err := vUC.voteRepo.Insert(vote); err != nil {
		return err
	}
	return nil
}

func (vUC *VoteUsecase) SelectVote(vote *models.Vote) error {
	if err := vUC.voteRepo.Select(vote); err != nil {
		return err
	}
	return nil
}

func (vUC *VoteUsecase) UpdateVote(vote *models.Vote) error {
	if err := vUC.voteRepo.Update(vote); err != nil {
		return err
	}
	return nil
}
