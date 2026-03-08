package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDB() {

	connStr := os.Getenv("SUPABASE_URL")

	var err error
	DB, err = sql.Open("pgx", connStr)

	if err != nil {
		log.Fatal("Unable to connect to Supabase:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}

	fmt.Println("Connected to Supabase")

	createTable()
}

func createTable() {

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT,
		email TEXT,
		avatar_url TEXT,
		last_login TIMESTAMP
	);`

	_, err := DB.Exec(query)

	if err != nil {
		log.Fatal("Could not create table:", err)
	}
}
