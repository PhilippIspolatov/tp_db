package usecase

import (
	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/tools"
	"github.com/ifo16u375/tp_db/internal/user"
)

type UserUsecase struct {
	userRepo user.Repository
}

func NewUserUsecase(ur user.Repository) user.Usecase {
	return &UserUsecase{
		userRepo: ur,
	}
}

func (uUC *UserUsecase) CreateUser(user *models.User) ([]*models.User, error) {
	var ConflictData []*models.User
	res1, err1 := uUC.userRepo.SelectByNickname(user.Nickname)
	res2, err2 := uUC.userRepo.SelectByEmail(user.Email)
	if err1 == nil && err2 == nil {
		if res1.Nickname == res2.Nickname {
			ConflictData = append(ConflictData, res1)
			return ConflictData, tools.ErrConflict
		}
		ConflictData = append(ConflictData, res1)
		ConflictData = append(ConflictData, res2)
		return ConflictData, tools.ErrConflict
	}

	if err1 == nil {
		ConflictData = append(ConflictData, res1)
		return ConflictData, tools.ErrConflict
	}
	if err2 == nil {
		ConflictData = append(ConflictData, res2)
		return ConflictData, tools.ErrConflict
	}

	if err := uUC.userRepo.Insert(user); err != nil {
		return nil, err
	}
	return nil, nil
}

func (uUC *UserUsecase) UpdateUser(user *models.User) error {
	if _, err := uUC.userRepo.SelectByEmail(user.Email); err == nil {
		return tools.ErrConflict
	}
	err := uUC.userRepo.Update(user); if err != nil {
		return err
	}
	return nil
}

func (uUC *UserUsecase) GetUser(nickname string) (*models.User, error) {
	u, err := uUC.userRepo.SelectByNickname(nickname)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (uUC *UserUsecase) GetAllUsers(slug string, limit uint64, since string, desc bool) ([]*models.User, error) {
	res, err := uUC.userRepo.SelectAllUsers(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}
	return res, nil
}
