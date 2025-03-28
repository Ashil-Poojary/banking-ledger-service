package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitPostgres() {
	var err error
	DB, err = sqlx.Connect("postgres", "host=localhost port=5432 user=user password=password dbname=banking sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to PostgreSQL:", err)
	}
	log.Println("Connected to PostgreSQL")
}
