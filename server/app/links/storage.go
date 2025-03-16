package links

import (
	"database/sql"
	"fmt"
	"outshort/app/common"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dbConnection *common.DbConnection) *Storage {
	return &Storage{
		db: dbConnection.Database(),
	}
}

func (this *Storage) CreateQuickLink(originalUrl string, alias string) (int64, *common.StorageError) {
	alreadyExists, err := this.AliasAlreadyExists(alias)
	if err != nil {
		return -1, common.NewStorageError(common.ErrorAny, err.Error())
	}
	if alreadyExists {
		return -1, common.NewStorageError(common.ErrorUniqueViolation, fmt.Sprintf("link with alias '%s' already exists", alias))
	}

	id, err_ := InsertQuickLink(this.db, originalUrl, alias)
	if err_ != nil {
		return -1, common.NewStorageError(common.ErrorAny, err_.Error())
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
	uid := common.GenerateLinkUid()
	linkModel, err_ := InsertLink(this.db, uid, originalUrl, name, alias, lifetime, ownerId)
	if err_ != nil {
		return nil, common.NewStorageError(common.ErrorAny, err_.Error())
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
	links, err := SelectLinksByOwner(this.db, ownerId)
	if err != nil {
		return nil, common.NewStorageError(common.ErrorAny, err.Error())
	}
	return links, nil
}

func (this *Storage) GetOriginalUrl(alias string) (string, *common.StorageError) {
	originalUrl, err := SelectOriginalUrl(this.db, alias)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", common.NewStorageError(common.ErrorNotFound, fmt.Sprintf("no record with alias = '%s'", alias))
		} else {
			return "", common.NewStorageError(common.ErrorAny, err.Error())
		}
	}
	return originalUrl, nil
}

func (this *Storage) AliasAlreadyExists(alias string) (bool, *common.StorageError) {
	exists, err := SelectExistsAlias(this.db, alias)
	if err != nil {
		return false, common.NewStorageError(common.ErrorAny, err.Error())
	}
	return exists, nil
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
