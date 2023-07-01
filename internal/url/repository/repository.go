package repository

import (
	"github.com/gocql/gocql"
	"time"
	"url_shortener/infrastructure/database"
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
	query := ur.db.Session.Query(`
			INSERT INTO short_url (uid, short_url, long_url, owner_uid, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, uid, short, long, ownerUID, time.Now(), time.Now())

	return query.Exec()
}

func (ur *urlRepo) FindShort(long string) (string, error) {
	q := ur.db.Session.Query(`
			SELECT short_url FROM short_url WHERE long_url = ?
		`, long)

	var res string
	err := q.Consistency(gocql.One).Scan(&res)
	return res, err

}

func (ur *urlRepo) FindLong(short string) (string, error) {
	q := ur.db.Session.Query(`
			SELECT long_url FROM short_url WHERE short_url = ?
		`, short)

	var res string
	err := q.Consistency(gocql.One).Scan(&res)
	return res, err
}

func (ur *urlRepo) IsShortURLUnique(short string) (bool, error) {
	var count int
	err := ur.db.Session.Query(`SELECT COUNT(*) FROM short_url WHERE short_url = ?`, short).Consistency(gocql.One).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil

}
