package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB opens the SQLite database at path (creates file if not exists)
func InitDB(path string) error {
	dsn := fmt.Sprintf("%s?_foreign_keys=on", path)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}
	// verify
	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}
	DB = db
	log.Printf("SQLite DB opened: %s", path)
	return nil
}

// InsertModel inserts a model row and returns the inserted id
func InsertModel(name, description, fileName, fileURL string, fileSize int64, archiveID, uploadedBy *int64) (int64, error) {
	res, err := DB.Exec(`INSERT INTO models (name, description, file_name, file_url, file_size, archive_id, uploaded_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, name, description, fileName, fileURL, fileSize, archiveID, uploadedBy)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetModelRow fetches file_name, file_url and archive_id for a model id
func GetModelRow(id int64) (fileName string, fileURL string, archiveID sql.NullInt64, err error) {
	row := DB.QueryRow(`SELECT file_name, file_url, archive_id FROM models WHERE id = ?`, id)
	var arch sql.NullInt64
	if err = row.Scan(&fileName, &fileURL, &arch); err != nil {
		return
	}
	archiveID = arch
	return
}

// DeleteModelByID deletes a model row and returns its file_name and archive_id (if any)
func DeleteModelByID(id int64) (fileName string, archiveID sql.NullInt64, err error) {
	// fetch first
	row := DB.QueryRow(`SELECT file_name, archive_id FROM models WHERE id = ?`, id)
	var arch sql.NullInt64
	if err = row.Scan(&fileName, &arch); err != nil {
		return
	}
	// delete
	if _, err = DB.Exec(`DELETE FROM models WHERE id = ?`, id); err != nil {
		return
	}
	archiveID = arch
	return
}

// GetArchiveNameByID returns archive name for given id
func GetArchiveNameByID(id int64) (string, error) {
	var name string
	err := DB.QueryRow(`SELECT name FROM archives WHERE id = ?`, id).Scan(&name)
	return name, err
}
