package usecase

import (
	"fmt"
	"log"
	"math/rand"
	"url_shortener/internal/repository"
)

const 	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

type urlUseCase struct {
	repo repository.StoreFinder
}

func NewURLUseCase(r repository.StoreFinder) *urlUseCase {
	return &urlUseCase{
		repo: r,
	}
}
func (uc *urlUseCase) Shorten(long string) string {
	short, err := uc.repo.FindShort(long)
	if err == nil {
		log.Printf("found cached short url %s",short)
		return short
	}
	log.Printf("finding cached short url:%v",err)

	short = uc.generateShortURL()

	uc.repo.Store(short, long)
	log.Printf("url shortened %s\n", short)

	return short
}

func (uc *urlUseCase) Expand(short string) (string,error) {
	long, err := uc.repo.FindLong(short)
	if err != nil {
		return "",fmt.Errorf("finding cached long url:%w",err)
	}
	log.Printf("found cached long url %s",long)
	return long, nil

}

func (uc *urlUseCase) generateShortURL() string {
	lengthConstraint := 5
	idBytes := make([]byte, lengthConstraint)


	for {
		for i := 0; i < lengthConstraint; i++ {
			idBytes[i] = chars[rand.Intn(len(chars))]
		}
		short := string(idBytes)

		if uc.repo.IsShortURLUnique(short) {
			return short
		}
		log.Printf("short url %s is not unique, generating new one\n", short)
	}

}