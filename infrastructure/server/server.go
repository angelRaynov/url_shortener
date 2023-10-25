package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
	"url_shortener/infrastructure/config"
)

type urlHandler interface {
	ShortenURL(c *gin.Context)
	Redirect(c *gin.Context)
	MyLinks(c *gin.Context)
}

type userHandler interface {
	Authenticate(c *gin.Context)
	Register(c *gin.Context)
	Edit(c *gin.Context)
}

type logger interface {
	Fatalf(string, ...interface{})
	Infof(string, ...interface{})
	Info(args ...interface{})
}

type server struct {
	log         logger
	urlHandler  urlHandler
	userHandler userHandler
	cfg         *config.Application
}

func NewServer(cfg *config.Application, log logger, uh urlHandler, ush userHandler) *server {
	return &server{
		log:         log,
		urlHandler:  uh,
		userHandler: ush,
		cfg:         cfg,
	}
}

func (s *server) Run() {
	router := gin.Default()
	router.POST("/shorten", AuthMiddleware(s.cfg), s.urlHandler.ShortenURL)
	router.GET("/my", AuthMiddleware(s.cfg), s.urlHandler.MyLinks)

	router.POST("/authenticate", s.userHandler.Authenticate)
	router.POST("/register", s.userHandler.Register)
	router.PATCH("/edit", AuthMiddleware(s.cfg), s.userHandler.Edit)
	router.GET("/:short", AuthMiddleware(s.cfg), s.urlHandler.Redirect)

	port := fmt.Sprintf(":%s", s.cfg.AppPort)

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		s.log.Infof("listening on port :%s", s.cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		s.log.Fatalf("Server shutdown failed: %v", err)
	}

	s.log.Info("Server stopped gracefully")
}
