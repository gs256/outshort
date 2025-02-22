package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	AnyError           = 0
	AliasAlreadyExists = 1
	NotFound           = 2
)

type StorageError struct {
	code    int
	message string
}

func NewStorageError(code int, message string) *StorageError {
	return &StorageError{code, message}
}

func (this *StorageError) Error() string {
	return this.message
}

const DbFile string = "database.db"

type Storage struct {
	db *sql.DB
}

func (this *Storage) Initialize() {
	db, err := sql.Open("sqlite3", DbFile)
	if err != nil {
		panic(err.Error())
	}
	this.db = db
}

func (this *Storage) Dispose() {
	defer this.db.Close()
}

func (this *Storage) CreateLink(originalUrl string, alias string) (int64, *StorageError) {
	alreadyExists, err := aliasAlreadyExists(this.db, alias)
	if err != nil {
		return -1, NewStorageError(AnyError, err.Error())
	}
	if alreadyExists {
		return -1, NewStorageError(AliasAlreadyExists, fmt.Sprintf("link with alias '%s' already exists", alias))
	}

	id, err := insertLink(this.db, originalUrl, alias)
	if err != nil {
		return -1, NewStorageError(AnyError, err.Error())
	}
	return id, nil
}

func (this *Storage) GetOriginalUrl(alias string) (string, *StorageError) {
	var originalURL string
	err := this.db.QueryRow("SELECT original_url FROM links WHERE alias = ? LIMIT 1", alias).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", NewStorageError(NotFound, fmt.Sprintf("No record with alias = '%s'", alias))
		} else {
			return "", NewStorageError(AnyError, err.Error())
		}
	}
	return originalURL, nil
}

func aliasAlreadyExists(db *sql.DB, alias string) (bool, error) {
	var exists int
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM links WHERE alias = ?)", alias).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}

func insertLink(db *sql.DB, originalUrl string, alias string) (int64, error) {
	res, err := db.Exec("INSERT INTO links (alias, original_url) VALUES(?, ?)", alias, originalUrl)
	if err != nil {
		return -1, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return -1, err
	}
	return id, nil
}
