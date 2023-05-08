package repository

import (
	"github.com/gocql/gocql"
	"time"
	"url_shortener/internal/pkg/database"
)

type URLRepo struct {
	db *database.DB
}

type StoreFinder interface {
	Store(short, long string) error
	FindShort(long string) (string,error)
	FindLong(short string) (string,error)
	IsShortURLUnique(short string) (bool, error)
}

func NewURLRepository(db *database.DB) StoreFinder {
	return &URLRepo{
		db: db,
	}
}

func (ur *URLRepo) Store(short, long string) error {
	query := ur.db.Session.Query(`
			INSERT INTO short_urls (short_url, long_url, created_at)
			VALUES (?, ?, ?)
		`, short, long, time.Now())

	return query.Exec()
}

func (ur *URLRepo) FindShort(long string) (string,error) {
	q := ur.db.Session.Query(`
			SELECT short_url FROM short_urls WHERE long_url = ?
		`, long)

	var res string
	err := q.Consistency(gocql.One).Scan(&res)
	return res,err

}

func (ur *URLRepo) FindLong(short string) (string,error) {
	q := ur.db.Session.Query(`
			SELECT long_url FROM short_urls WHERE short_url = ?
		`, short)

	var res string
	err := q.Consistency(gocql.One).Scan(&res)
	return res,err
}

func (ur *URLRepo) IsShortURLUnique(short string) (bool, error) {
	var count int
	err := ur.db.Session.Query(`SELECT COUNT(*) FROM short_urls WHERE short_url = ?`, short).Consistency(gocql.One).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil

}

