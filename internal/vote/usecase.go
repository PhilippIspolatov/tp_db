package vote

import "github.com/ifo16u375/tp_db/internal/models"

type Usecase interface {
	InsertVote(vote *models.Vote) error
	SelectVote(vote *models.Vote) error
	UpdateVote(vote *models.Vote) error
}