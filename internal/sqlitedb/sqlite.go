package sqlitedb

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	CreateTables(db)
	return db, err
}

func CreateTables(DB *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			hashed_password CHAR(60) NOT NULL,
			token TEXT,
			created DATETIME NOT NULL
		);
		CREATE TABLE IF NOT EXISTS snippets (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			author_name TEXT NOT NULL,
			title VARCHAR(100) NOT NULL,
			content TEXT NOT NULL,
			created DATETIME NOT NULL
		);
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category TEXT NOT NULL UNIQUE
		);
		INSERT OR IGNORE INTO categories (category) VALUES
		("HTML/CSS"),
		("Golang"),
		("Rust"),
		("JavaScript"),
		("Other");
		CREATE TABLE IF NOT EXISTS post_cat (
			cat_id INTEGER,
			post_id INTEGER		
		);
		CREATE TABLE IF NOT EXISTS LikeReact (
			post_id INTEGER,
			user_id INTEGER NOT NULL,
			reaction INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS ComentReact (
			comment_id INTEGER,
			user_id INTEGER NOT NULL,
			type INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			author_name TEXT NOT NULL,
			content TEXT NOT NULL,
			created DATETIME NOT NULL		
		);
		`
	_, err := DB.Exec(query)
	return err
}
