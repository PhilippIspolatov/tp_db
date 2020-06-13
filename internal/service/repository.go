package service

import "github.com/ifo16u375/tp_db/internal/models"

type Repository interface {
	ClearAllDB() error
	SelectStatusDB() (*models.Service, error)
}
