package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func connectDB() *sql.DB {
	connStr := "postgres://postgres:HariSatz@localhost:5432/audio_uploader?sslmode=disable"

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	fmt.Println("Connected to the database")
	db = database
	return db
}

func createTableIfNotExists(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS audio_files (
		id SERIAL PRIMARY KEY,
		filename TEXT,
		file_data BYTEA,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("audio_files table is ready")
}
