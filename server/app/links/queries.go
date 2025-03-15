package links

import (
	"database/sql"
	"fmt"
)

const LinkColumns = "id, uid, alias, original_url, name, lifetime_sec, created_at, owner_id"

func InsertQuickLink(db *sql.DB, originalUrl string, alias string) (int64, error) {
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

func InsertLink(db *sql.DB, uid string, originalUrl string, name string, alias string, lifetime int, ownerId int64) (*LinkModel, error) {
	res, err := db.Exec("INSERT INTO links (uid, alias, original_url, name, lifetime_sec, owner_id) VALUES(?, ?, ?, ?, ?, ?)", uid, alias, originalUrl, name, lifetime, ownerId)
	if err != nil {
		return nil, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return nil, err
	}
	linkModel, err := GetLinkById(db, id)
	if err != nil {
		return nil, err
	}
	return linkModel, nil
}

func UpdateLink(db *sql.DB, uid string, originalUrl string, name string, alias string, lifetime int) (*LinkModel, error) {
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
	linkModel, err := GetLinkByUid(db, uid)
	if err != nil {
		return nil, err
	}
	return linkModel, nil
}

func GetLinkById(db *sql.DB, id int64) (*LinkModel, error) {
	query := fmt.Sprintf("SELECT %s FROM links WHERE id = ?", LinkColumns)
	row := db.QueryRow(query, id)
	linkModel, err := ScanLinkModel(row)
	if err != nil {
		return nil, err
	}
	return linkModel, nil
}

func GetLinkByUid(db *sql.DB, uid string) (*LinkModel, error) {
	query := fmt.Sprintf("SELECT %s FROM links WHERE uid = ?", LinkColumns)
	row := db.QueryRow(query, uid)
	linkModel, err := ScanLinkModel(row)
	if err != nil {
		return nil, err
	}
	return linkModel, nil
}

func GetLinkByAlias(db *sql.DB, alias string) (*LinkModel, error) {
	query := fmt.Sprintf("SELECT %s FROM links WHERE alias = ?", LinkColumns)
	row := db.QueryRow(query, alias)
	linkModel, err := ScanLinkModel(row)
	if err != nil {
		return nil, err
	}
	return linkModel, nil
}

func ScanLinkModel(row *sql.Row) (*LinkModel, error) {
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
