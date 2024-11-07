package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbuname = "chatuser"
	dbpass  = "qwerty"
	dbname  = "chat"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	dbcon := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbuname, dbpass, dbname)

	db, err := sql.Open("mysql", dbcon)
	if err != nil {
		log.Println("error connect to db: ", err)
		return nil, err
	}

	log.Println("DB connected")
	fmt.Println("DB connected")

	return db, nil
}
