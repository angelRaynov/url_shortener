package main

import (
	"url_shortener/infrastructure/cache"
	"url_shortener/infrastructure/config"
	"url_shortener/infrastructure/database"
	"url_shortener/infrastructure/logger"
	"url_shortener/infrastructure/server"
	urlHTTP "url_shortener/internal/url/delivery/http"
	urlRepo "url_shortener/internal/url/repository"
	urlUseCase "url_shortener/internal/url/usecase"
	userHTTP "url_shortener/internal/user/delivery/http"
	userRepo "url_shortener/internal/user/repository"
	userUseCase "url_shortener/internal/user/usecase"
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

	ar := userRepo.NewAuthRepository(db)
	auc := userUseCase.NewUserUseCase(cfg, l, ar)
	ah := userHTTP.NewAuthHandler(l, auc)

	srv := server.NewServer(cfg, l, uh, ah)

	srv.Run()
}
