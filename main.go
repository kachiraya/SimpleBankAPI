package main

import (
	"database/sql"
	"log"
	"os"
	"simplebankapi/user"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	user.StartServer(":"+os.Getenv("PORT"), db)
}
