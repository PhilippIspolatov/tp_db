package repository

import (
	"database/sql"

	"github.com/ifo16u375/tp_db/internal/forum"
	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/sirupsen/logrus"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(db *sql.DB) forum.Repository {
	return &ForumRepository{
		db: db,
	}
}

func (fr *ForumRepository) Insert(forum *models.Forum) error {
	err := fr.db.QueryRow("INSERT INTO forums (slug, title, owner) VALUES ($1, $2, $3) "+
		"RETURNING posts, slug, threads, title, owner", forum.Slug, forum.Title,
		forum.User).Scan(&forum.Posts, &forum.Slug,
		&forum.Threads, &forum.Title, &forum.User)
	if err != nil{
		return err
	}
	return nil
}

func (fr *ForumRepository) Select(slug string) (*models.Forum, error) {
	f := &models.Forum{}

	if err := fr.db.QueryRow("SELECT posts, slug, threads, title, owner FROM forums " +
		"WHERE lower(slug)=lower($1)", slug).Scan(&f.Posts, &f.Slug, &f.Threads, &f.Title, &f.User); err != nil {
		logrus.Info(err)
		return nil, err
	}
	return f, nil
}
