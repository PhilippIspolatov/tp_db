package repository

import (
	"database/sql"
	"fmt"

	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/post"
	"github.com/ifo16u375/tp_db/internal/tools"
)

type PostRepository struct{
	db *sql.DB
}

func NewPostRepository(db *sql.DB) post.Repository {
	return &PostRepository{
		db:db,
	}
}

func (pr* PostRepository) Insert(posts []*models.Post, thread uint64, forum string) error {
	QueryString := "INSERT INTO posts (author, forum, message, parent, thread) VALUES "

	for _, post := range posts {
		QueryString += fmt.Sprintf("('%s', '%s', '%s', '%d', '%d'), ",
			post.Author, forum, post.Message, post.Parent, thread)
	}

	qr := []rune(QueryString)
	qr[len(qr) - 2] = ' '
	QueryString = string(qr)



	QueryString += "RETURNING author, created, forum, id, isEdited, message, parent, thread"

	res, err := pr.db.Query(QueryString)

	if err != nil {
		return err
	}

	i := 0
	for res.Next() {
		if err := res.Scan(&posts[i].Author, &posts[i].Created, &posts[i].Forum,
			&posts[i].Id, &posts[i].IsEdited, &posts[i].Message, &posts[i].Parent,
			&posts[i].Thread); err != nil {
			return err
		}
		i++
	}

	return nil
}

func (pr *PostRepository) CheckPostsByParentAndAuthor(posts []*models.Post, id uint64) error {
	p := map[uint64]uint64{}
	a := map[string]uint64{}

	QueryString1 := "SELECT COUNT (nickname) FROM users WHERE nickname IN ("
	QueryString2 := "SELECT COUNT(id) FROM posts WHERE thread = $1 AND id IN ("
	for _, post := range posts{
		QueryString1 += fmt.Sprintf("'%s',", post.Author)
		a[post.Author]++
		if post.Parent > 0 {
			QueryString2 += fmt.Sprintf("%d,", post.Parent)
			p[post.Parent]++
		}
	}

	qr := []rune(QueryString1)
	qr[len(qr) - 1] = ')'
	QueryString1 = string(qr)

	qr = []rune(QueryString2)
	qr[len(qr) - 1] = ')'
	QueryString2 = string(qr)

	count := 0
	if err := pr.db.QueryRow(QueryString1).Scan(&count); err != nil {
		return err
	}

	if count != len(a) {
		return tools.ErrNotFound
	}

	if len(p) == 0 {
		return nil
	}

	if err := pr.db.QueryRow(QueryString2, id).Scan(&count); err != nil {
		return err
	}

	if count != len(p) {
		return tools.ErrConflict
	}

	return nil
}

func (pr *PostRepository) SelectPost(id uint64) (*models.Post, error) {
	p := &models.Post{}
	if err := pr.db.QueryRow("SELECT * FROM posts WHERE id = $1", id).Scan(
		&p.Author, &p.Created, &p.Forum, &p.Id,  &p.IsEdited,
		&p.Message, &p.Parent, &p.Thread); err != nil {
		return nil, err
	}
	return p, nil
}

func (pr *PostRepository) UpdatePost(post *models.Post) error {
	if err := pr.db.QueryRow("UPDATE posts SET message=coalesce(nullif($1, ''), message) WHERE id=$2 "+
		"RETURNING author, created, forum, id, isEdited, message, parent, thread",
		post.Message, post.Id).Scan(&post.Author, &post.Created, &post.Forum,
		&post.Id, &post.IsEdited, &post.Message, &post.Parent,
		&post.Thread); err != nil {
		return err
	}
	return nil
}
