package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //vai pro crl
)

// Connect conecta com o banco
func Connect() (*sql.DB, error) {
	strConn := "dev:rootroot@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", strConn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
