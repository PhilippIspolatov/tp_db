package vote

import "github.com/ifo16u375/tp_db/internal/models"

type Repository interface {
	Insert(vote *models.Vote) error
	Select(vote *models.Vote) error
	Update(vote *models.Vote) error
}
