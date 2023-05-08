package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"url_shortener/config"
	"url_shortener/internal/pkg/cache"
	"url_shortener/internal/pkg/database"
	"url_shortener/internal/url/delivery/http"
	"url_shortener/internal/url/repository"
	"url_shortener/internal/url/usecase"
)

func main() {
	cfg := config.New()
	db := database.InitDB(cfg)
	redis := cache.NewURLCache(cfg)
	repo := repository.NewURLRepository(db)
	useCase := usecase.NewURLUseCase(cfg, repo, redis)
	handler := http.NewURLHandler(cfg, useCase)

	router := gin.Default()
	router.POST("/", handler.ShortenURL)
	router.GET("/:shortened", handler.ExpandURL)

	log.Println("listening on port :" + cfg.AppPort)

	log.Fatal(router.Run(":" + cfg.AppPort))

}
