package repository

import (
	"github.com/PhilippIspolatov/tp_db/internal/models"
	"github.com/PhilippIspolatov/tp_db/internal/service"
	"github.com/jackc/pgx"
)

type ServiceRepository struct {
	db *pgx.ConnPool
}

func NewServiceRepository(db *pgx.ConnPool) service.Repository {
	return &ServiceRepository{
		db:db,
	}
}

func (sr *ServiceRepository) ClearAllDB() error {
	if _, err := sr.db.Exec("" +
		"TRUNCATE TABLE forums_users; " +
		"TRUNCATE TABLE forums CASCADE; " +
		"TRUNCATE TABLE posts CASCADE; " +
		"TRUNCATE TABLE threads CASCADE; " +
		"TRUNCATE TABLE users CASCADE; " +
		"TRUNCATE TABLE votes CASCADE; "); err != nil {
		return err
	}
	return nil
}

func (sr *ServiceRepository) SelectStatusDB() (*models.Service, error) {
	s := &models.Service{}

	if err := sr.db.QueryRow("SELECT * FROM (SELECT COUNT(id) FROM posts) as posts, "+
		"(SELECT COUNT(slug) FROM forums) as forums, " +
		"(SELECT COUNT(id) FROM threads) as threads, " +
		"(SELECT COUNT(nickname) FROM users) as users").Scan(&s.Post,
			&s.Forum, &s.Thread, &s.User); err != nil {
		return nil, err
	}
	return s, nil
}