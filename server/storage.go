package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

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

type LinkModel struct {
	Id          int64
	Uid         string
	Alias       string
	OriginalUrl string
	Name        string
	LifetimeSec int
	CreatedAt   time.Time
	OwnerId     sql.NullInt64
}

const LinkColumns = "id, uid, alias, original_url, name, lifetime_sec, created_at, owner_id"

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

func (this *Storage) CreateQuickLink(originalUrl string, alias string) (int64, *StorageError) {
	alreadyExists, err := this.AliasAlreadyExists(alias)
	if err != nil {
		return -1, NewStorageError(AnyError, err.Error())
	}
	if alreadyExists {
		return -1, NewStorageError(UniqueViolation, fmt.Sprintf("link with alias '%s' already exists", alias))
	}

	id, err := insertQuickLink(this.db, originalUrl, alias)
	if err != nil {
		return -1, NewStorageError(AnyError, err.Error())
	}
	return id, nil
}

func (this *Storage) CreateLink(originalUrl string, name string, alias string, lifetime int, ownerId int64) (*LinkModel, *StorageError) {
	alreadyExists, err := this.AliasAlreadyExists(alias)
	if err != nil {
		return nil, NewStorageError(AnyError, err.Error())
	}
	if alreadyExists {
		return nil, NewStorageError(UniqueViolation, fmt.Sprintf("link with alias '%s' already exists", alias))
	}
	uid := generateLinkUid()
	linkModel, err := insertLink(this.db, uid, originalUrl, name, alias, lifetime, ownerId)
	if err != nil {
		return nil, NewStorageError(AnyError, err.Error())
	}
	return linkModel, nil
}

func (this *Storage) UpdateLink(uid string, originalUrl string, name string, alias string, lifetime int, ownerId int64) (*LinkModel, *StorageError) {
	linkModel, err := updateLink(this.db, uid, originalUrl, name, alias, lifetime)
	if err != nil {
		return nil, NewStorageError(AnyError, err.Error())
	}
	return linkModel, nil
}

func (this *Storage) GetAllLinks(ownerId int64) ([]LinkModel, *StorageError) {
	rows, err := this.db.Query("SELECT id, uid, alias, original_url, name, lifetime_sec, created_at, owner_id FROM links WHERE owner_id = ?", ownerId)
	if err != nil {
		return nil, NewStorageError(AnyError, "Query error")
	}
	defer rows.Close()
	var links []LinkModel
	for rows.Next() {
		var link LinkModel
		err := rows.Scan(
			&link.Id,
			&link.Uid,
			&link.Alias,
			&link.OriginalUrl,
			&link.Name,
			&link.LifetimeSec,
			&link.CreatedAt,
			&link.OwnerId,
		)
		if err != nil {
			return nil, NewStorageError(AnyError, "Failed to scan row")
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, NewStorageError(AnyError, "Unknown error")
	}
	return links, nil
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

func (this *Storage) AliasAlreadyExists(alias string) (bool, error) {
	var exists int
	err := this.db.QueryRow("SELECT EXISTS(SELECT 1 FROM links WHERE alias = ?)", alias).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}

func (this *Storage) FindLinkByUid(uid string) (*LinkModel, *StorageError) {
	link, err := getLinkByUid(this.db, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewStorageError(NotFound, err.Error())
		}
		return nil, NewStorageError(AnyError, err.Error())
	}
	return link, nil
}

func (this *Storage) FindLinkByAlias(alias string) (*LinkModel, *StorageError) {
	link, err := getLinkByAlias(this.db, alias)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewStorageError(NotFound, err.Error())
		}
		return nil, NewStorageError(AnyError, err.Error())
	}
	return link, nil
}

func insertQuickLink(db *sql.DB, originalUrl string, alias string) (int64, error) {
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

func insertLink(db *sql.DB, uid string, originalUrl string, name string, alias string, lifetime int, ownerId int64) (*LinkModel, error) {
	res, err := db.Exec("INSERT INTO links (uid, alias, original_url, name, lifetime_sec, owner_id) VALUES(?, ?, ?, ?, ?, ?)", uid, alias, originalUrl, name, lifetime, ownerId)
	if err != nil {
		return nil, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return nil, err
	}
	linkModel, err := getLinkById(db, id)
	if err != nil {
		return nil, err
	}
	return linkModel, nil
}

func updateLink(db *sql.DB, uid string, originalUrl string, name string, alias string, lifetime int) (*LinkModel, error) {
	res, err := db.Exec("UPDATE links SET alias = ?, original_url = ?, name = ?, lifetime_sec = ? WHERE uid = ?", alias, originalUrl, name, lifetime, uid)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no rows affected")
	}
	linkModel, err := getLinkByUid(db, uid)
	if err != nil {
		return nil, err
	}
	return linkModel, nil
}

func getLinkById(db *sql.DB, id int64) (*LinkModel, error) {
	query := fmt.Sprintf("SELECT %s FROM links WHERE id = ?", LinkColumns)
	row := db.QueryRow(query, id)
	linkModel, err := scanLinkModel(row)
	if err != nil {
		return nil, err
	}
	return linkModel, nil
}

func getLinkByUid(db *sql.DB, uid string) (*LinkModel, error) {
	query := fmt.Sprintf("SELECT %s FROM links WHERE uid = ?", LinkColumns)
	row := db.QueryRow(query, uid)
	linkModel, err := scanLinkModel(row)
	if err != nil {
		return nil, err
	}
	return linkModel, nil
}

func getLinkByAlias(db *sql.DB, alias string) (*LinkModel, error) {
	query := fmt.Sprintf("SELECT %s FROM links WHERE alias = ?", LinkColumns)
	row := db.QueryRow(query, alias)
	linkModel, err := scanLinkModel(row)
	if err != nil {
		return nil, err
	}
	return linkModel, nil
}

func scanLinkModel(row *sql.Row) (*LinkModel, error) {
	var linkModel LinkModel
	err := row.Scan(
		&linkModel.Id,
		&linkModel.Uid,
		&linkModel.Alias,
		&linkModel.OriginalUrl,
		&linkModel.Name,
		&linkModel.LifetimeSec,
		&linkModel.CreatedAt,
		&linkModel.OwnerId,
	)
	return &linkModel, err
}

func generateLinkUid() string {
	return RandomString(16)
}
