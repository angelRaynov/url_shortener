package server

import (
	"github.com/gin-gonic/gin"
	"url_shortener/infrastructure/config"
)

type ShortExpander interface {
	ShortenURL(c *gin.Context)
	ExpandURL(c *gin.Context)
}

type Logger interface {
	Infof(string, ...interface{})
	Fatal(...interface{})
	//Infow(string, ...interface{})
	//Errorw(string, ...interface{})
	//Warnw(string, ...interface{})
	//Debugw(string, ...interface{})
	//Info(...interface{})
	//Error(...interface{})
	//Warn(...interface{})
	//Debug(...interface{})
	//Errorf(string, ...interface{})
	//Warnf(string, ...interface{})
	//Debugf(string, ...interface{})
}

type server struct {
	log     Logger
	handler ShortExpander
	cfg     *config.Application
}

func NewServer(cfg *config.Application, log Logger, h ShortExpander) *server {
	return &server{
		log:     log,
		handler: h,
		cfg:     cfg,
	}
}

func (s *server) Run() {
	router := gin.Default()
	router.POST("/shorten", AuthMiddleware(s.cfg), s.handler.ShortenURL)
	router.POST("/expand", s.handler.ExpandURL)

	s.log.Infof("listening on port :%s", s.cfg.AppPort)

	s.log.Fatal(router.Run(":" + s.cfg.AppPort))
}
