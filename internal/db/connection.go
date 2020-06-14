package db

import (
	"encoding/json"
	"log"
	"os"

	"github.com/jackc/pgx"
	"github.com/kr/pretty"
)

type db struct {
	Conn *pgx.ConnPool
}

type connInfo struct {
	Host     string
	Port     uint32
	Database string
	User     string
	Password string
}

func NewDataBase(confPath string) (*db, error) {
	conf, err := loadConf(confPath)

	if err != nil {
		return nil, err
	}

	connStr, err := pgx.ParseConnectionString(conf)

	if err != nil {
		return nil, err
	}

	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connStr,
		MaxConnections: 100,
	})

	for err != nil {
		conn, err = pgx.NewConnPool(
			pgx.ConnPoolConfig{
				ConnConfig:     connStr,
				MaxConnections: 100,
			})
	}

	return &db{
		Conn: conn,
	}, nil

}

func loadConf(path string) (string, error) {
	db := &connInfo{}

	connStr := ""

	f, err := os.Open(path)

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		return connStr, err
	}

	pars := json.NewDecoder(f)

	if err := pars.Decode(db); err != nil {
		return connStr, err
	}

	connStr = pretty.Sprintf("postgresql://%s:%s@%s:%d/%s", db.User, db.Password, db.Host, db.Port, db.Database)

	return connStr, nil
}
