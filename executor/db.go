package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Initialize() {
	db, err := sql.Open("sqlite3", "executor.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = db
	migrateTables()
}

func migrateTables() {
	// Images table stores references to the images information, in case they are not yet built.
	DB.Exec(`
		CREATE TABLE IF NOT EXISTS images (
			id TEXT PRIMARY KEY,
			name TEXT,
			image_build TEXT,
			version TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)

	// Session is a single batch run of code. It can produce multiple run results.
	DB.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			image_tag TEXT,
			strictness TEXT,
			timeout INTEGER,
			memory_limit INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME DEFAULT NULL
		)
	`)

	// Runs store each run results. The same session ID runs the same code.
	DB.Exec(`
		CREATE TABLE IF NOT EXISTS runs (
			id TEXT PRIMARY KEY,
			session_id TEXT,
			stdin TEXT,
			stdout TEXT,
			stderr TEXT,
			judge TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
}

type Session struct {
	ID          string
	ImageTag    string
	Strictness  string
	Timeout     int
	MemoryLimit int
}

func CreateSession(session *Session) error {
	stmt, err := DB.Prepare("INSERT INTO sessions (id, image_tag, strictness, timeout, memory_limit) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.ID, session.ImageTag, session.Strictness, session.Timeout, session.MemoryLimit)
	if err != nil {
		return err
	}

	return nil
}
