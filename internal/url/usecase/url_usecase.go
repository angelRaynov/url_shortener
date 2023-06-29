package usecase

import (
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"url_shortener/infrastructure/config"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

type getCacher interface {
	GetShort(long string) (string, error)
	GetLong(short string) (string, error)
	Cache(short, long string) error
}

type storeFinder interface {
	Store(short, long string) error
	FindShort(long string) (string, error)
	FindLong(short string) (string, error)
	IsShortURLUnique(short string) (bool, error)
}

type urlUseCase struct {
	cfg    *config.Application
	repo   storeFinder
	cache  getCacher
	logger *zap.SugaredLogger
}

func NewURLUseCase(cfg *config.Application, r storeFinder, cache getCacher, logger *zap.SugaredLogger) *urlUseCase {
	return &urlUseCase{
		repo:   r,
		cfg:    cfg,
		cache:  cache,
		logger: logger,
	}
}
func (uc *urlUseCase) Shorten(long string) (string, error) {
	//check redis
	short, err := uc.cache.GetShort(long)
	if err == nil {
		uc.logger.Debugw("found short url in cache", "short_url", short, "long_url", long)
		return short, nil
	}

	//check db
	short, err = uc.repo.FindShort(long)
	if err == nil {
		uc.logger.Debugw("found short url in db", "short_url", short, "long_url", long)
		return short, nil
	}

	id, err := uc.generateShortURL()
	if err != nil {
		return "", fmt.Errorf("assuring url uniqueness:%w", err)
	}

	short = uc.cfg.AppURL + id

	err = uc.repo.Store(short, long)
	if err != nil {
		return "", fmt.Errorf("storing url:%w", err)
	}

	err = uc.cache.Cache(short, long)
	if err != nil {
		return "", fmt.Errorf("caching url:%w", err)
	}

	uc.logger.Infow("url shortened", "short_url", short, "long_url", long)
	return short, nil
}

func (uc *urlUseCase) Expand(short string) (string, error) {
	//check redis
	long, err := uc.cache.GetLong(short)
	if err == nil {
		uc.logger.Debugw("found long url in cache ", "long_url", long, "short_url", short)
		return long, nil
	}

	//check db
	long, err = uc.repo.FindLong(short)
	if err != nil {
		return "", fmt.Errorf("finding cached long url:%w", err)
	}

	uc.logger.Debugw("found long url in db ", "long_url", long, "short_url", short)
	return long, nil

}

func (uc *urlUseCase) generateShortURL() (string, error) {
	lengthConstraint := 7
	idBytes := make([]byte, lengthConstraint)

	for {
		for i := 0; i < lengthConstraint; i++ {
			idBytes[i] = chars[rand.Intn(len(chars))]
		}
		short := string(idBytes)

		isUnique, err := uc.repo.IsShortURLUnique(short)
		if err != nil {
			return "", fmt.Errorf("checking for uniqueness:%w", err)
		}

		if isUnique {
			return short, nil
		}
		uc.logger.Warnw("short url is not unique, generating new one", "short_url", short)
	}

}
