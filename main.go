package main

import (
	"database/sql"
	"log"
	"os"
	"simplebankapi-heroku/user"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        first_name TEXT,
        last_name TEXT
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}
	createTable = `
	CREATE TABLE IF NOT EXISTS bankaccounts (
        id SERIAL PRIMARY KEY,
        userid SERIAL,
        account_number SERIAL,
        name TEXT,
        balance NUMBER
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	user.StartServer(":"+os.Getenv("PORT"), db)
}
