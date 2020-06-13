package repository

import (
	"database/sql"

	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/vote"
)

type VoteRepository struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) vote.Repository {
	return &VoteRepository{
		db:db,
	}
}

func (vr *VoteRepository) Insert(vote *models.Vote) error {
	if _, err := vr.db.Exec("INSERT INTO votes (nickname, thread, voice) VALUES "+
		"($1, $2, $3)", vote.Nickname, vote.Thread, vote.Voice); err != nil {
		return err
	}
	return nil
}

func (vr *VoteRepository) Select(vote *models.Vote) error {
	if err := vr.db.QueryRow("SELECT thread, voice FROM votes WHERE thread=$1 "+
		"AND lower(nickname)=lower($2)", vote.Thread, vote.Nickname).Scan(&vote.Thread, &vote.Voice); err != nil {
		return err
	}
	return nil
}

func (vr *VoteRepository) Update(vote *models.Vote) error {
	if err := vr.db.QueryRow("UPDATE votes SET voice = $1 WHERE lower(nickname)=lower($2) "+
		" AND thread = $3 RETURNING voice",
		vote.Voice, vote.Nickname, vote.Thread).Scan(&vote.Voice); err != nil {
		return err
	}
	return nil
}