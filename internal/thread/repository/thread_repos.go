package repository

import (
	"database/sql"

	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/thread"

)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) thread.Repository {
	return &ThreadRepository{
		db: db,
	}
}

func (tr *ThreadRepository) UpdateById(thread *models.Thread) error {
	if err := tr.db.QueryRow("UPDATE threads SET "+
		"message=coalesce(nullif($1, ''), message), "+
		"title=coalesce(nullif($2, ''), title) "+
		"WHERE id = $3 "+
		"RETURNING author, created, forum, id, message, slug, title, votes",
		thread.Message, thread.Title, thread.Id).Scan(&thread.Author,
			&thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug,
			&thread.Title, &thread.Votes); err != nil {
		return err
	}
	return nil
}

func (tr *ThreadRepository) UpdateBySlug(thread *models.Thread) error {
	if err := tr.db.QueryRow("UPDATE threads SET "+
		"message=coalesce(nullif($1, ''), message), "+
		"title=coalesce(nullif($2, ''), title) "+
		"WHERE lower(slug) = lower($3) "+
		"RETURNING author, created, forum, id, message, slug, title, votes",
		thread.Message, thread.Title, thread.Slug).Scan(&thread.Author,
		&thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug,
		&thread.Title, &thread.Votes); err != nil {
		return err
	}
	return nil
}

func (tr *ThreadRepository) Insert(thread *models.Thread) error {
	if err := tr.db.QueryRow("INSERT INTO threads (author, created, forum, message, slug, title, votes) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", thread.Author, thread.Created, thread.Forum, thread.Message,
		thread.Slug, thread.Title, thread.Votes).Scan(&thread.Id); err != nil {
		return err
	}
	return nil
}

func (tr *ThreadRepository) SelectBySlug(slug string) (*models.Thread, error) {
	t := &models.Thread{}

	if err := tr.db.QueryRow("SELECT author, created, forum, id, "+
		"message, slug, title, votes FROM threads WHERE lower(slug)=lower($1)", slug).Scan(&t.Author, &t.Created,
		&t.Forum, &t.Id, &t.Message, &t.Slug, &t.Title, &t.Votes); err != nil {
		return nil, err
	}
	return t, nil
}

func (tr *ThreadRepository) SelectById(id uint64) (*models.Thread, error) {
	t := &models.Thread{}

	if err := tr.db.QueryRow("SELECT author, created, forum, id, "+
		"message, slug, title, votes FROM threads WHERE id=$1", id).Scan(&t.Author, &t.Created,
		&t.Forum, &t.Id, &t.Message, &t.Slug, &t.Title, &t.Votes); err != nil {
		return nil, err
	}
	return t, nil
}

func (tr *ThreadRepository) Select(slug string, limit uint64, since string, desc bool) ([]*models.Thread, error) {
	t := []*models.Thread{}
	var err error
	var res *sql.Rows

	QueryString := "SELECT * FROM threads where lower(forum)=lower($1) "

	switch {
	case limit > 0 && since != "":
		{
			if desc {
				QueryString += "AND created <= $2 ORDER BY created DESC "
			} else {
				QueryString += "AND created >= $2 ORDER BY created ASC "
			}
			QueryString += "LIMIT $3"

			res, err = tr.db.Query(QueryString, slug, since, limit)
			if err != nil {
				return t, err
			}
		}
	case limit > 0:
		{
			if desc {
				QueryString += "ORDER BY created DESC "
			} else {
				QueryString += "ORDER BY created ASC "
			}
			QueryString += "LIMIT $2"

			res, err = tr.db.Query(QueryString, slug, limit)
			if err != nil {
				return t, err
			}
		}
	case since != "":
		{
			if desc {
				QueryString += "AND created <= $2 ORDER BY created DESC "
			} else {
				QueryString += "AND created <= $2 ORDER BY created ASC "
			}

			res, err = tr.db.Query(QueryString, slug, since)
			if err != nil {
				return t, err
			}
		}
	}

	for res.Next() {
		th := &models.Thread{}

		if err := res.Scan(&th.Author, &th.Created, &th.Forum,
			&th.Id, &th.Message, &th.Slug, &th.Title, &th.Votes); err != nil {
			return nil, err
		}
		t = append(t, th)
	}
	return t, nil
}
