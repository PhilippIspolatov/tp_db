package repository

import (
	"database/sql"

	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/service"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) service.Repository {
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