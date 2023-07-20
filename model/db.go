package model

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func InitDB() (*sqlx.DB, error) {
	user := os.Getenv("NS_MARIADB_USER")
	if user == "" {
		user = "mehm8128"
	}

	pass := os.Getenv("NS_MARIADB_PASSWORD")
	if pass == "" {
		pass = "math8128"
	}

	host := os.Getenv("NS_MARIADB_HOSTNAME")
	if host == "" {
		host = "localhost"
	}

	dbname := os.Getenv("NS_MARIADB_DATABASE")
	if dbname == "" {
		dbname = "traQuest"
	}
	_db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true&loc=Local", user, pass, host, dbname))
	if err != nil {
		return nil, err
	}
	db = _db

	return db, err
}
