package repository

import (
	"database/sql"

	"github.com/ifo16u375/tp_db/internal/models"
	"github.com/ifo16u375/tp_db/internal/tools"
	"github.com/ifo16u375/tp_db/internal/user"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Insert(user *models.User) error {
	if _, err := ur.db.Exec("INSERT INTO users VALUES($1, $2, $3, $4)",
		user.Nickname, user.FullName, user.Email, user.About); err != nil {
			logrus.Info(err)
			return err
	}
	return nil
}

func (ur *UserRepository) Update(user *models.User) error {
	err := ur.db.QueryRow("UPDATE users SET full_name=coalesce(nullif($1, ''), full_name), "+
		"email=coalesce(nullif($2, ''), email), about=coalesce(nullif($3, ''), about) "+
		"WHERE lower(nickname)=lower($4) RETURNING full_name, email, about, nickname",
		user.FullName, user.Email, user.About, user.Nickname).Scan(&user.FullName, &user.Email,
			&user.About, &user.Nickname)

	if err == sql.ErrNoRows {
		logrus.Error(err)
		return tools.ErrNotFound
	}

	if err != nil {
		logrus.Info(err)
		return err
	}

	return nil
}


func (ur *UserRepository) SelectByNickname(nickname string) (*models.User, error) {
	u := &models.User{}
	if err := ur.db.QueryRow("SELECT nickname, full_name, email, about FROM users " +
		"WHERE lower(nickname)=lower($1)", nickname).Scan(&u.Nickname, &u.FullName, &u.Email, &u.About); err != nil {
		logrus.Info(err)
		return nil, err

	}
	return u, nil
}

func (ur *UserRepository) SelectByEmail(email string) (*models.User, error) {
	u := &models.User{}
	if err := ur.db.QueryRow("SELECT nickname, full_name, email, about FROM users " +
		"WHERE lower(email)=lower($1)", email).Scan(&u.Nickname, &u.FullName, &u.Email, &u.About); err != nil {
		logrus.Info(err)
		return nil, err
	}
	return u, nil
}