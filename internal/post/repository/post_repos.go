package repository

import (
	"fmt"

	"github.com/PhilippIspolatov/tp_db/internal/models"
	"github.com/PhilippIspolatov/tp_db/internal/post"
	"github.com/PhilippIspolatov/tp_db/internal/tools"
	"github.com/jackc/pgx"
)

type PostRepository struct {
	db *pgx.ConnPool
}

func NewPostRepository(db *pgx.ConnPool) post.Repository {
	return &PostRepository{
		db: db,
	}
}

func (pr *PostRepository) Insert(posts []*models.Post, thread uint64, forum string) error {
	QueryString := "INSERT INTO posts (author, forum, message, parent, thread) VALUES "

	for _, post := range posts {
		QueryString += fmt.Sprintf("('%s', '%s', '%s', '%d', '%d'), ",
			post.Author, forum, post.Message, post.Parent, thread)
	}

	qr := []rune(QueryString)
	qr[len(qr)-2] = ' '
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
	for _, post := range posts {
		QueryString1 += fmt.Sprintf("'%s',", post.Author)
		a[post.Author]++
		if post.Parent > 0 {
			QueryString2 += fmt.Sprintf("%d,", post.Parent)
			p[post.Parent]++
		}
	}

	qr := []rune(QueryString1)
	qr[len(qr)-1] = ')'
	QueryString1 = string(qr)

	qr = []rune(QueryString2)
	qr[len(qr)-1] = ')'
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

	if err := pr.db.QueryRow("SELECT author, created, forum, id, isEdited, message, " +
		"parent, thread FROM posts WHERE id = $1", id).Scan(
		&p.Author, &p.Created, &p.Forum, &p.Id, &p.IsEdited,
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

func (pr *PostRepository) GetPostsByFlat(threadId uint64, desc bool, since uint64, limit uint64) ([]*models.Post, error) {
	posts := []*models.Post{}

	QueryString := fmt.Sprintf("SELECT author, created, forum, id, isEdited, message, " +
		"parent, thread FROM posts WHERE thread = %d ", threadId)

	if since > 0 {
		if desc {
			QueryString += fmt.Sprintf("AND id < %d ", since)
		} else {
			QueryString += fmt.Sprintf("AND id > %d ", since)
		}
	}

	QueryString += "ORDER BY ID "

	if desc {
		QueryString += "DESC "
	}

	if limit > 0 {
		QueryString += fmt.Sprintf("LIMIT %d ", limit)
	}

	rows, err := pr.db.Query(QueryString)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := &models.Post{}

		if err := rows.Scan(&p.Author, &p.Created, &p.Forum, &p.Id, &p.IsEdited, &p.Message, &p.Parent, &p.Thread); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func (pr *PostRepository) GetPostsByTree(threadId uint64, desc bool, since uint64, limit uint64) ([]*models.Post, error) {
	posts := []*models.Post{}

	QueryString := fmt.Sprintf("SELECT author, created, forum, id, isEdited, message, " +
		"parent, thread FROM posts WHERE thread = %d ", threadId)

	if since > 0 {
		if desc {
			QueryString += fmt.Sprintf("AND PATH < (SELECT path FROM posts WHERE id = %d) ", since)
		} else {
			QueryString += fmt.Sprintf("AND PATH > (SELECT path FROM posts WHERE id = %d) ", since)
		}
	}

	if desc {
		QueryString += "ORDER BY path[1] DESC, path DESC "
	} else {
		QueryString += "ORDER BY path[1], path "
	}

	if limit > 0 {
		QueryString += fmt.Sprintf("LIMIT %d", limit)
	}

	rows, err := pr.db.Query(QueryString)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := &models.Post{}

		if err := rows.Scan(&p.Author, &p.Created, &p.Forum, &p.Id, &p.IsEdited, &p.Message,
			&p.Parent, &p.Thread); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func (pr *PostRepository) GetPostsByParentTree(threadId uint64, desc bool, since uint64, limit uint64) ([]*models.Post, error) {
	posts := []*models.Post{}

	QueryString := fmt.Sprintf("SELECT author, created, forum, id, isEdited, message, " +
	"parent, thread FROM posts WHERE path[1] IN (SELECT id FROM posts WHERE thread = %d AND parent = 0 ", threadId)

	if since > 0 {
		if desc {
			QueryString += fmt.Sprintf("AND path[1] < (SELECT path[1] FROM posts WHERE id = %d) ", since)
		} else {
			QueryString += fmt.Sprintf("AND path[1] > (SELECT path[1] FROM posts WHERE id = %d) ", since)
		}
	}

	if desc {
		QueryString += "ORDER BY id DESC "
	} else {
		QueryString += "ORDER BY id "
	}

	if limit > 0 {
		QueryString += fmt.Sprintf("LIMIT %d) ", limit)
	}

	if desc {
		QueryString += "ORDER BY path[1] DESC, path, id "
	} else {
		QueryString += "ORDER BY path "
	}

	rows, err := pr.db.Query(QueryString)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := &models.Post{}

		if err := rows.Scan(&p.Author, &p.Created, &p.Forum, &p.Id, &p.IsEdited, &p.Message,
			&p.Parent, &p.Thread); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, nil
}
