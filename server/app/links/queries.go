package links

import (
	"database/sql"
	"fmt"
)

type Scannable interface {
	Scan(...any) error
}

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

func SelectLinksByOwner(db *sql.DB, ownerId int64) ([]LinkModel, error) {
	query := fmt.Sprintf("SELECT %s FROM links WHERE owner_id = ?", LinkColumns)
	rows, err := db.Query(query, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var links []LinkModel
	for rows.Next() {
		link, err := ScanLinkModel(rows)
		if err != nil {
			return nil, err
		}
		links = append(links, *link)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}

func SelectOriginalUrl(db *sql.DB, alias string) (string, error) {
	var originalURL string
	err := db.QueryRow("SELECT original_url FROM links WHERE alias = ? LIMIT 1", alias).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func SelectExistsAlias(db *sql.DB, alias string) (bool, error) {
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

func ScanLinkModel(scannable Scannable) (*LinkModel, error) {
	var linkModel LinkModel
	err := scannable.Scan(
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
