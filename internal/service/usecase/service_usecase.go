package usecase

import (
	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/service"
)

type ServiceUsecase struct {
	serviceRepo service.Repository
}

func NewServiceUsecase(sr service.Repository) service.Usecase {
	return &ServiceUsecase{
		serviceRepo:sr,
	}
}

func (sUC *ServiceUsecase) ClearAllDB() error {
	if err := sUC.serviceRepo.ClearAllDB(); err != nil {
		return err
	}
	return nil
}

func (sUC *ServiceUsecase) GetInfoDB() (*models.Service, error) {
	s, err := sUC.serviceRepo.SelectStatusDB()
	if err != nil {
		return nil, err
	}
	return s, nil
}