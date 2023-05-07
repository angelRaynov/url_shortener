package repository

import (
	"fmt"
)

type URLRepo struct {
	shortLong map[string]string
	longShort map[string]string
	unique map[string]bool
}

type StoreFinder interface {
	Store(short,long string)
	FindShort(long string) (string,error)
	FindLong(short string) (string,error)
	IsShortURLUnique(short string) bool
}

func NewURLRepository() *URLRepo {
	return &URLRepo{
		shortLong: make(map[string]string),
		longShort: make(map[string]string),
		unique: make(map[string]bool),
	}
}

func (ur *URLRepo) Store(short,long string) {
	ur.shortLong[short] = long
	ur.longShort[long] = short
}

func (ur *URLRepo) FindShort(long string) (string,error) {
	if short, ok := ur.longShort[long]; ok {
		return short,nil
	}
	return "",fmt.Errorf("short url not found")

}

func (ur *URLRepo) FindLong(short string) (string,error) {
	if long, ok := ur.shortLong[short]; ok {
		return long, nil
	}
	return "",fmt.Errorf("long url not found")
}

func (ur *URLRepo) IsShortURLUnique(short string) bool {
	if ur.unique[short] {
		return false
	}
	return true

}

