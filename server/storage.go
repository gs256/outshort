package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

const (
	AnyError           = 0
	UniqueViolation    = 1
	NotFound           = 2
	InvalidCredentials = 3
)

type StorageError struct {
	code    int
	message string
}

type UserModel struct {
	Id       int64
	Username string
	Password string
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
		return -1, NewStorageError(UniqueViolation, fmt.Sprintf("link with alias '%s' already exists", alias))
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

func (this *Storage) CreateUser(username string, password string) (int64, *StorageError) {
	res, err := this.db.Exec("INSERT INTO users (username, password) VALUES(?, ?)", username, password)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
			return -1, NewStorageError(UniqueViolation, fmt.Sprintf("Username is already taken"))
		}
		return -1, NewStorageError(AnyError, "Failed to create user")
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return -1, NewStorageError(AnyError, "Failed to get user id")
	}
	return id, nil
}

func (this *Storage) CreateAuthToken(token string, userId int64, lifetimeSec int) *StorageError {
	_, err := this.db.Exec("INSERT INTO auth_tokens (token, user_id, lifetime_sec) VALUES(?, ?, ?)", token, userId, lifetimeSec)
	if err != nil {
		return NewStorageError(AnyError, "Failed to create auth token")
	}
	return nil
}

func (this *Storage) AuthenticateUser(username string, password string) (int64, *StorageError) {
	var user UserModel
	err := this.db.QueryRow("SELECT id, username, password FROM users WHERE username = ? LIMIT 1", username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, NewStorageError(NotFound, "User not found")
		}
		return -1, NewStorageError(AnyError, "Unknown error")
	}
	if user.Password == password {
		return user.Id, nil
	}
	return -1, NewStorageError(InvalidCredentials, "Invalid credentials")
}

func (this *Storage) GetUserInfo(authToken string) (*UserModel, *StorageError) {
	var user UserModel
	err := this.db.QueryRow("SELECT users.* FROM users JOIN auth_tokens ON users.id = auth_tokens.user_id WHERE auth_tokens.token = ?", authToken).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewStorageError(NotFound, "User not found")
		}
		return nil, NewStorageError(AnyError, "Unknown error")
	}
	return &user, nil
}

func (this *Storage) DeleteAuthToken(authToken string) *StorageError {
	_, err := this.db.Exec("DELETE FROM auth_tokens WHERE token = ?", authToken)
	if err != nil {
		return NewStorageError(AnyError, "Unknown error")
	}
	return nil
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
