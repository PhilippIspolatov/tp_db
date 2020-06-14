package service

import "github.com/PhilippIspolatov/tp_db/internal/models"

type Repository interface {
	ClearAllDB() error
	SelectStatusDB() (*models.Service, error)
}
