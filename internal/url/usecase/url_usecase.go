package usecase

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"math/rand"
	"url_shortener/infrastructure/config"
	"url_shortener/internal/model"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

type getCacher interface {
	GetShort(long string) (string, error)
	GetLong(short string) (string, error)
	Cache(short, long string) error
}

type storeFinder interface {
	Store(uid string, ownerUID string, short string, sr model.ShortenRequest) error
	FindShort(long string) (string, error)
	FindLong(short string) (string, error)
	IsShortURLUnique(short string) (bool, error)
	FindLinksPerUser(ownerID string) ([]model.URL, error)
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
func (uc *urlUseCase) Shorten(ownerID string, ur model.ShortenRequest) (string, error) {
	//check redis
	short, err := uc.cache.GetShort(ur.LongURL)
	if err == nil {
		uc.logger.Debugw("found short url in cache", "short_url", short, "long_url", ur.LongURL)
		return short, nil
	}

	//check db
	short, err = uc.repo.FindShort(ur.LongURL)
	if err == nil {
		uc.logger.Debugw("found short url in db", "short_url", short, "long_url", ur.LongURL)
		return short, nil
	}

	id, err := uc.generateShortURL()
	if err != nil {
		return "", fmt.Errorf("assuring url uniqueness:%w", err)
	}

	short = uc.cfg.AppURL + id
	//use custom domain if present
	if ur.Domain != "" {
		short = fmt.Sprintf("%s/%s", ur.Domain, id)
	}

	if ur.Title == "" {
		ur.Title = short
	}
	uid := uuid.New()
	err = uc.repo.Store(uid.String(), ownerID, short, ur)
	if err != nil {
		return "", fmt.Errorf("storing url:%w", err)
	}

	err = uc.cache.Cache(short, ur.LongURL)
	if err != nil {
		return "", fmt.Errorf("caching url:%w", err)
	}

	uc.logger.Infow("url shortened", "short_url", short, "long_url", ur.LongURL)
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

func (uc *urlUseCase) FetchLinksPerUser(ownerID string) ([]model.URL, error) {
	//check db
	links, err := uc.repo.FindLinksPerUser(ownerID)
	if err != nil {
		return nil, fmt.Errorf("finding links per user:%w", err)
	}

	uc.logger.Debugw("links per user fetched", "count", len(links), "owner_id", ownerID)
	return links, nil

}
