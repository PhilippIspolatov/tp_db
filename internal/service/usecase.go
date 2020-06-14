package service

import "github.com/PhilippIspolatov/tp_db/internal/models"

type Usecase interface {
	ClearAllDB() error
	GetInfoDB() (*models.Service, error)
}
