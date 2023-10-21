package repository

import (
	"database/sql"
	"url_shortener/internal/model"
)

const (
	QueryStore            = `INSERT INTO url (uid, title, domain, short_url, long_url, owner_uid) VALUES (?, ?, ?, ?, ?, ?)`
	QueryFindShort        = `SELECT short_url FROM url WHERE long_url = ?`
	QueryFindLong         = `SELECT long_url FROM url WHERE short_url = ?`
	QueryIsUnique         = `SELECT COUNT(*) FROM url WHERE short_url = ?`
	QueryFindLinksPerUser = `SELECT uid, title, domain, short_url, long_url, owner_uid, created_at, updated_at FROM url WHERE owner_uid = ?`
)

type urlRepo struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *urlRepo {
	return &urlRepo{
		db: db,
	}
}

// TODO: ADD edit link functionality
func (ur *urlRepo) Store(uid string, ownerUID string, short string, sr model.ShortenRequest) error {
	_, err := ur.db.Exec(QueryStore, uid, sr.Title, sr.Domain, short, sr.LongURL, ownerUID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *urlRepo) FindShort(long string) (string, error) {
	row := ur.db.QueryRow(QueryFindShort, long)

	var res string
	err := row.Scan(&res)
	return res, err

}

func (ur *urlRepo) FindLong(short string) (string, error) {
	row := ur.db.QueryRow(QueryFindLong, short)

	var res string
	err := row.Scan(&res)
	return res, err
}

func (ur *urlRepo) IsShortURLUnique(short string) (bool, error) {
	var count int
	err := ur.db.QueryRow(QueryIsUnique, short).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil

}

func (ur *urlRepo) FindLinksPerUser(ownerID string) ([]model.URL, error) {
	rows, err := ur.db.Query(QueryFindLinksPerUser, ownerID)
	if err != nil {
		return nil, err
	}
	var url model.URL
	var res []model.URL

	for rows.Next() {
		err = rows.Scan(&url.UID, &url.ShortURL, &url.LongURL, &url.OwnerID, &url.CreatedAt, &url.UpdatedAt)
		if err != nil {
			return res, err
		}
		res = append(res, url)
	}

	if err := rows.Close(); err != nil {
		return res, err
	}

	return res, nil
}
