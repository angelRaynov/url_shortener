package repository

import (
	"github.com/gocql/gocql"
	"time"
	"url_shortener/infrastructure/database"
	"url_shortener/internal/model"
)

const (
	QueryStore            = `INSERT INTO short_url (uid, short_url, long_url, owner_uid, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	QueryFindShort        = `SELECT short_url FROM short_url WHERE long_url = ?`
	QueryFindLong         = `SELECT long_url FROM short_url WHERE short_url = ?`
	QueryIsUnique         = `SELECT COUNT(*) FROM short_url WHERE short_url = ?`
	QueryFindLinksPerUser = `SELECT uid, short_url, long_url, owner_uid, created_at, updated_at FROM short_url WHERE owner_uid = ?`
)

type urlRepo struct {
	db *database.DB
}

func NewURLRepository(db *database.DB) *urlRepo {
	return &urlRepo{
		db: db,
	}
}

func (ur *urlRepo) Store(uid, ownerUID, short, long string) error {
	query := ur.db.Session.Query(QueryStore, uid, short, long, ownerUID, time.Now(), time.Now())

	return query.Exec()
}

func (ur *urlRepo) FindShort(long string) (string, error) {
	q := ur.db.Session.Query(QueryFindShort, long)

	var res string
	err := q.Consistency(gocql.One).Scan(&res)
	return res, err

}

func (ur *urlRepo) FindLong(short string) (string, error) {
	q := ur.db.Session.Query(QueryFindLong, short)

	var res string
	err := q.Consistency(gocql.One).Scan(&res)
	return res, err
}

func (ur *urlRepo) IsShortURLUnique(short string) (bool, error) {
	var count int
	err := ur.db.Session.Query(QueryIsUnique, short).Consistency(gocql.One).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil

}

func (ur *urlRepo) FindLinksPerUser(ownerID string) ([]model.URL, error) {
	iter := ur.db.Session.Query(QueryFindLinksPerUser, ownerID).Iter()

	var url model.URL
	var res []model.URL

	for iter.Scan(&url.UID, &url.ShortURL, &url.LongURL, &url.OwnerID, &url.CreatedAt, &url.UpdatedAt) {
		res = append(res, url)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return res, nil
}
