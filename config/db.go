package config

import (
	"database/sql"
	//так надо
	_ "github.com/lib/pq"
)

//Env db
type Env struct {
	DB *sql.DB
}

//NewDB content
func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
