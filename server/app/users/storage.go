package users

import (
	"database/sql"
	"errors"
	"fmt"
	"outshort/app/common"

	"github.com/mattn/go-sqlite3"
)

const DbFile string = "database.db"

type UserModel struct {
	Id       int64
	Username string
	Password string
}

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

func (this *Storage) CreateUser(username string, password string) (int64, *common.StorageError) {
	res, err := this.db.Exec("INSERT INTO users (username, password) VALUES(?, ?)", username, password)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
			return -1, common.NewStorageError(common.ErrorUniqueViolation, fmt.Sprintf("Username is already taken"))
		}
		return -1, common.NewStorageError(common.ErrorAny, "Failed to create user")
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return -1, common.NewStorageError(common.ErrorAny, "Failed to get user id")
	}
	return id, nil
}

func (this *Storage) CreateAuthToken(token string, userId int64, lifetimeSec int) *common.StorageError {
	_, err := this.db.Exec("INSERT INTO auth_tokens (token, user_id, lifetime_sec) VALUES(?, ?, ?)", token, userId, lifetimeSec)
	if err != nil {
		return common.NewStorageError(common.ErrorAny, "Failed to create auth token")
	}
	return nil
}

func (this *Storage) AuthenticateUser(username string, password string) (int64, *common.StorageError) {
	var user UserModel
	err := this.db.QueryRow("SELECT id, username, password FROM users WHERE username = ? LIMIT 1", username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, common.NewStorageError(common.ErrorNotFound, "User not found")
		}
		return -1, common.NewStorageError(common.ErrorAny, "Unknown error")
	}
	if user.Password == password {
		return user.Id, nil
	}
	return -1, common.NewStorageError(common.ErrorInvalidCredentials, "Invalid credentials")
}

// TODO: fix duplication
func (this *Storage) GetUserInfo(authToken string) (*UserModel, *common.StorageError) {
	var user UserModel
	err := this.db.QueryRow("SELECT users.* FROM users JOIN auth_tokens ON users.id = auth_tokens.user_id WHERE auth_tokens.token = ?", authToken).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.NewStorageError(common.ErrorNotFound, "User not found")
		}
		return nil, common.NewStorageError(common.ErrorAny, "Unknown error")
	}
	return &user, nil
}

func (this *Storage) DeleteAuthToken(authToken string) *common.StorageError {
	_, err := this.db.Exec("DELETE FROM auth_tokens WHERE token = ?", authToken)
	if err != nil {
		return common.NewStorageError(common.ErrorAny, "Unknown error")
	}
	return nil
}
