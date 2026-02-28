package main

import (
	"database/sql"
	"fmt"
)

func openSQLite(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := ensureSchema(db); err != nil {
		return nil, err
	}
	return db, nil
}

func ensureSchema(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  markdown TEXT NOT NULL,
  html TEXT NOT NULL,
  created_at TEXT NOT NULL
);
`
	_, err := db.Exec(schema)
	return err
}

func insertPost(db *sql.DB, markdown string, html string, createdAt string) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO posts (markdown, html, created_at) VALUES (?, ?, ?)`,
		markdown,
		html,
		createdAt,
	)
	if err != nil {
		return 0, fmt.Errorf("insert post: %w", err)
	}
	return result.LastInsertId()
}
