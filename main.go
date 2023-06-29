package main

import (
	"url_shortener/infrastructure/cache"
	"url_shortener/infrastructure/config"
	"url_shortener/infrastructure/database"
	"url_shortener/infrastructure/logger"
	"url_shortener/infrastructure/server"
	"url_shortener/internal/url/delivery/http"
	"url_shortener/internal/url/repository"
	"url_shortener/internal/url/usecase"
)

func main() {
	cfg := config.New()

	l := logger.Init(cfg.AppMode)

	db := database.InitDB(cfg, l)
	repo := repository.NewURLRepository(db)

	c := cache.NewURLCache(cfg, l)
	cacheRepo := repository.NewCacheRepo(c)

	useCase := usecase.NewURLUseCase(cfg, repo, cacheRepo, l)

	handler := http.NewURLHandler(cfg, useCase, l)

	srv := server.NewServer(cfg, l, handler)

	srv.Run()
}
