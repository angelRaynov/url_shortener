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

type ShortExpander interface {
	ShortenURL(c *gin.Context)
	ExpandURL(c *gin.Context)
	MyLinks(c *gin.Context)
}

type Authenticator interface {
	Authenticate(c *gin.Context)
	Register(c *gin.Context)
	Edit(c *gin.Context)
}

type Logger interface {
	Fatalf(string, ...interface{})
	Infof(string, ...interface{})
	Info(args ...interface{})
}

type server struct {
	log         Logger
	urlHandler  ShortExpander
	authHandler Authenticator
	cfg         *config.Application
}

func NewServer(cfg *config.Application, log Logger, uh ShortExpander, ah Authenticator) *server {
	return &server{
		log:         log,
		urlHandler:  uh,
		authHandler: ah,
		cfg:         cfg,
	}
}

func (s *server) Run() {
	router := gin.Default()
	router.POST("/shorten", AuthMiddleware(s.cfg), s.urlHandler.ShortenURL)
	router.POST("/expand", AuthMiddleware(s.cfg), s.urlHandler.ExpandURL)
	router.GET("/my", AuthMiddleware(s.cfg), s.urlHandler.MyLinks)

	router.POST("/authenticate", s.authHandler.Authenticate)
	router.POST("/register", s.authHandler.Register)
	router.PATCH("/edit", AuthMiddleware(s.cfg), s.authHandler.Edit)

	port := fmt.Sprintf(":%s", s.cfg.AppPort)

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatalf("Server failed to start: %v", err)
		}
		s.log.Infof("listening on port :%s", s.cfg.AppPort)
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
