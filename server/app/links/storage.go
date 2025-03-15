package links

import (
	"database/sql"
	"fmt"
	"outshort/app/common"
	"outshort/app/users"
)

type Storage struct {
	db *sql.DB
}

func generateLinkUid() string {
	return common.RandomString(16)
}

func (this *Storage) Initialize(dbConnection *common.DbConnection) {
	this.db = dbConnection.Database()
}

func (this *Storage) CreateQuickLink(originalUrl string, alias string) (int64, *common.StorageError) {
	alreadyExists, err := this.AliasAlreadyExists(alias)
	if err != nil {
		return -1, common.NewStorageError(common.ErrorAny, err.Error())
	}
	if alreadyExists {
		return -1, common.NewStorageError(common.ErrorUniqueViolation, fmt.Sprintf("link with alias '%s' already exists", alias))
	}

	id, err := InsertQuickLink(this.db, originalUrl, alias)
	if err != nil {
		return -1, common.NewStorageError(common.ErrorAny, err.Error())
	}
	return id, nil
}

func (this *Storage) CreateLink(originalUrl string, name string, alias string, lifetime int, ownerId int64) (*LinkModel, *common.StorageError) {
	alreadyExists, err := this.AliasAlreadyExists(alias)
	if err != nil {
		return nil, common.NewStorageError(common.ErrorAny, err.Error())
	}
	if alreadyExists {
		return nil, common.NewStorageError(common.ErrorUniqueViolation, fmt.Sprintf("link with alias '%s' already exists", alias))
	}
	uid := generateLinkUid()
	linkModel, err := InsertLink(this.db, uid, originalUrl, name, alias, lifetime, ownerId)
	if err != nil {
		return nil, common.NewStorageError(common.ErrorAny, err.Error())
	}
	return linkModel, nil
}

func (this *Storage) UpdateLink(uid string, originalUrl string, name string, alias string, lifetime int, ownerId int64) (*LinkModel, *common.StorageError) {
	linkModel, err := UpdateLink(this.db, uid, originalUrl, name, alias, lifetime)
	if err != nil {
		return nil, common.NewStorageError(common.ErrorAny, err.Error())
	}
	return linkModel, nil
}

func (this *Storage) GetAllLinks(ownerId int64) ([]LinkModel, *common.StorageError) {
	rows, err := this.db.Query("SELECT id, uid, alias, original_url, name, lifetime_sec, created_at, owner_id FROM links WHERE owner_id = ?", ownerId)
	if err != nil {
		return nil, common.NewStorageError(common.ErrorAny, "Query error")
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
			return nil, common.NewStorageError(common.ErrorAny, "Failed to scan row")
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, common.NewStorageError(common.ErrorAny, "Unknown error")
	}
	return links, nil
}

func (this *Storage) GetOriginalUrl(alias string) (string, *common.StorageError) {
	var originalURL string
	err := this.db.QueryRow("SELECT original_url FROM links WHERE alias = ? LIMIT 1", alias).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", common.NewStorageError(common.ErrorNotFound, fmt.Sprintf("No record with alias = '%s'", alias))
		} else {
			return "", common.NewStorageError(common.ErrorAny, err.Error())
		}
	}
	return originalURL, nil
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

func (this *Storage) FindLinkByUid(uid string) (*LinkModel, *common.StorageError) {
	link, err := GetLinkByUid(this.db, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.NewStorageError(common.ErrorNotFound, err.Error())
		}
		return nil, common.NewStorageError(common.ErrorAny, err.Error())
	}
	return link, nil
}

func (this *Storage) FindLinkByAlias(alias string) (*LinkModel, *common.StorageError) {
	link, err := GetLinkByAlias(this.db, alias)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.NewStorageError(common.ErrorNotFound, err.Error())
		}
		return nil, common.NewStorageError(common.ErrorAny, err.Error())
	}
	return link, nil
}

// TODO: fix duplication
func (this *Storage) GetUserInfo(authToken string) (*users.UserModel, *common.StorageError) {
	var user users.UserModel
	err := this.db.QueryRow("SELECT users.* FROM users JOIN auth_tokens ON users.id = auth_tokens.user_id WHERE auth_tokens.token = ?", authToken).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.NewStorageError(common.ErrorNotFound, "User not found")
		}
		return nil, common.NewStorageError(common.ErrorAny, "Unknown error")
	}
	return &user, nil
}
