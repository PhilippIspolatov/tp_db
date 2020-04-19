package service

import "github.com/ifo16u375/tp_db/internal/models"

type Usecase interface {
	ClearAllDB() error
	GetInfoDB() (*models.Service, error)
}
