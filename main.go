package main

import (
	"github.com/gin-gonic/gin"
	"url_shortener/config"
	"url_shortener/internal/pkg/cache"
	"url_shortener/internal/pkg/database"
	"url_shortener/internal/url/delivery/http"
	"url_shortener/internal/url/repository"
	"url_shortener/internal/url/usecase"
	"url_shortener/logger"
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

	router := gin.Default()
	router.POST("/", handler.ShortenURL)
	router.GET("/:shortened", handler.ExpandURL)

	l.Infof("listening on port :%s", cfg.AppPort)

	l.Fatal(router.Run(":" + cfg.AppPort))

}
