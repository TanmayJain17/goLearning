package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tanmax2000"
	dbname   = "gofruitcart"
)

var db *sql.DB

func Connect() {

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	DB, err := sql.Open("postgres", connStr)
	if err != nil {

		panic(err)
	}
	db = DB

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected")

}

func GetDb() *sql.DB {
	return db
}
