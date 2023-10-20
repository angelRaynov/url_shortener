package main

import (
	"url_shortener/infrastructure/cache"
	"url_shortener/infrastructure/config"
	"url_shortener/infrastructure/database"
	"url_shortener/infrastructure/logger"
	"url_shortener/infrastructure/server"
	authHTTP "url_shortener/internal/authentication/delivery/http"
	authRepo "url_shortener/internal/authentication/repository"
	authUseCase "url_shortener/internal/authentication/usecase"
	urlHTTP "url_shortener/internal/url/delivery/http"
	urlRepo "url_shortener/internal/url/repository"
	urlUseCase "url_shortener/internal/url/usecase"
)

func main() {
	cfg := config.New()

	l := logger.Init(cfg.AppMode)

	db := database.Init(cfg, l)
	repo := urlRepo.NewURLRepository(db)

	c := cache.NewURLCache(cfg, l)
	cacheRepo := urlRepo.NewCacheRepo(c)

	useCase := urlUseCase.NewURLUseCase(cfg, repo, cacheRepo, l)

	uh := urlHTTP.NewURLHandler(cfg, useCase, l)

	ar := authRepo.NewAuthRepository(db)
	auc := authUseCase.NewAuthUseCase(cfg, l, ar)
	ah := authHTTP.NewAuthHandler(l, auc)

	srv := server.NewServer(cfg, l, uh, ah)

	srv.Run()
}
