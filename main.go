package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"url_shortener/config"
	"url_shortener/internal/pkg/database"
	"url_shortener/internal/repository"
	"url_shortener/internal/url/delivery/http"
	"url_shortener/internal/url/usecase"
)

func main() {
	cfg := config.New()
	db := database.InitDB(cfg)
	repo := repository.NewURLRepository(db)
	useCase := usecase.NewURLUseCase(cfg, repo)
	handler := http.NewURLHandler(cfg, useCase)

	router := gin.Default()
	router.POST("/", handler.ShortenURL)
	router.GET("/:shortened", handler.ExpandURL)

	log.Println("listening on port :" + cfg.AppPort)

	log.Fatal(router.Run(":" + cfg.AppPort))

}
