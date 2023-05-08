package usecase

import (
	"fmt"
	"log"
	"math/rand"
	"url_shortener/config"
	"url_shortener/internal/pkg/cache"
	"url_shortener/internal/repository"
	"url_shortener/internal/url"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

type urlUseCase struct {
	cfg  *config.Application
	repo repository.StoreFinder
	cache cache.GetCacher
}

func NewURLUseCase( cfg *config.Application, r repository.StoreFinder, cache cache.GetCacher) url.ShortExpander {
	return &urlUseCase{
		repo: r,
		cfg:  cfg,
		cache: cache,
	}
}
func (uc *urlUseCase) Shorten(long string) (string,error) {
	//check redis
	short, err := uc.cache.GetShort(long)
	if err == nil {
		log.Printf("found short url in cache %s", short)
		return short,nil
	}

	//check db
	short, err = uc.repo.FindShort(long)
	if err == nil {
		log.Printf("found short url in db %s", short)
		return short, nil
	}

	id,err := uc.generateShortURL()
	if err != nil {
		log.Printf("checking for uniqueness:%v", err)
	}

	short = uc.cfg.AppURL + id

	err = uc.repo.Store(short, long)
	if err != nil {
		return "", fmt.Errorf("storing url:%w",err)
	}
	err = uc.cache.Cache(short, long)
	if err != nil {
		return "", fmt.Errorf("caching url:%w",err)
	}
	log.Printf("url shortened %s\n", short)

	return short,nil
}

func (uc *urlUseCase) Expand(short string) (string, error) {
	//check redis
	long, err := uc.cache.GetLong(short)
	if err == nil {
		log.Printf("found short url in cache %s", short)
		return long,nil
	}

	//check db
	long, err = uc.repo.FindLong(short)
	if err != nil {
		return "", fmt.Errorf("finding cached long url:%w", err)
	}
	log.Printf("found cached long url %s", long)
	return long, nil

}

func (uc *urlUseCase) generateShortURL() (string,error) {
	lengthConstraint := 5
	idBytes := make([]byte, lengthConstraint)

	for {
		for i := 0; i < lengthConstraint; i++ {
			idBytes[i] = chars[rand.Intn(len(chars))]
		}
		short := string(idBytes)

		isUnique, err := uc.repo.IsShortURLUnique(short)
		if err != nil {
			return "",fmt.Errorf("generating short url:%w",err)
		}

		if isUnique {
			return short, nil
		}
		log.Printf("short url %s is not unique, generating new one\n", short)
	}

}
