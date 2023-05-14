package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"url_shortener/config"
	"url_shortener/internal/model"
	"url_shortener/internal/url"
)

type urlHandler struct {
	urlUseCase url.ShortExpander
	cfg        *config.Application
	logger     *zap.SugaredLogger
}

func NewURLHandler(cfg *config.Application, uc url.ShortExpander, logger *zap.SugaredLogger) url.ShortenExpandHandler {
	return &urlHandler{
		urlUseCase: uc,
		cfg:        cfg,
		logger:     logger,
	}
}

func (uh *urlHandler) ShortenURL(c *gin.Context) {
	var urlRequest model.ShortenRequest

	if err := c.BindJSON(&urlRequest); err != nil {
		uh.logger.Warnw("binding request params", "error", err)
		c.JSON(http.StatusBadRequest, "invalid payload")
		return
	}

	short, err := uh.urlUseCase.Shorten(urlRequest.LongURL)
	if err != nil {
		uh.logger.Warnw("shortening url", "error", err)
		c.JSON(http.StatusInternalServerError, "unable to shorten url at the moment, please try again later")
	}

	c.IndentedJSON(http.StatusCreated, model.ShortenResponse{
		Link:    short,
		LongURL: urlRequest.LongURL,
	})
	return

}

func (uh *urlHandler) ExpandURL(c *gin.Context) {
	shortenedURL := c.Param("shortened")

	if shortenedURL == "" {
		uh.logger.Warnw("unable to expand empty url")
		c.JSON(http.StatusBadRequest, "shortened url can not be empty")
	}

	shortenedURL = uh.cfg.AppURL + strings.TrimSpace(shortenedURL)

	long, err := uh.urlUseCase.Expand(shortenedURL)
	if err != nil {
		uh.logger.Warnw("expanding url", "short_url", shortenedURL, "error", err)
		c.IndentedJSON(http.StatusNotFound, "corresponding long url not found")
		return
	}

	c.Redirect(http.StatusMovedPermanently, long)
	return
}
