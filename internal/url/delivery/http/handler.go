package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"url_shortener/helper"
	"url_shortener/infrastructure/config"
	"url_shortener/internal/model"
)

type shortExpander interface {
	Shorten(ownerID string, urlRequest model.ShortenRequest) (string, error)
	Expand(short string) (string, error)
	FetchLinksPerUser(ownerID string) ([]model.URL, error)
}

type urlHandler struct {
	urlUseCase shortExpander
	cfg        *config.Application
	logger     *zap.SugaredLogger
}

func NewURLHandler(cfg *config.Application, uc shortExpander, logger *zap.SugaredLogger) *urlHandler {
	return &urlHandler{
		urlUseCase: uc,
		cfg:        cfg,
		logger:     logger,
	}
}

// TODO: Add redirect functionality. i.e when someone click shortened url to redirect to the long one
func (uh *urlHandler) ShortenURL(c *gin.Context) {
	var urlRequest model.ShortenRequest
	if err := c.BindJSON(&urlRequest); err != nil {
		uh.logger.Warnw("binding request params", "error", err)
		c.JSON(http.StatusBadRequest, "invalid payload")
		return
	}

	urlRequest.LongURL = helper.LinkPreconditioning(urlRequest.LongURL)

	if !helper.IsValidURL(urlRequest.LongURL) {
		uh.logger.Debugw("invalid url", "long_url", urlRequest.LongURL)
		c.JSON(http.StatusBadRequest, "invalid url")
		return
	}

	contextUser, err := helper.GetUserFromContext(c)
	if err != nil {
		uh.logger.Warnw("fetching user from context", "error", err)
		c.JSON(http.StatusInternalServerError, "unable to shorten url at the moment, please try again later")
		return
	}

	short, err := uh.urlUseCase.Shorten(contextUser.UID, urlRequest)
	if err != nil {
		uh.logger.Warnw("shortening url", "error", err)
		c.JSON(http.StatusInternalServerError, "unable to shorten url at the moment, please try again later")
		return
	}

	c.IndentedJSON(http.StatusCreated, model.ShortenResponse{
		Link:    short,
		LongURL: urlRequest.LongURL,
	})
	return

}

func (uh *urlHandler) ExpandURL(c *gin.Context) {
	var urlRequest model.ExpandRequest
	if err := c.BindJSON(&urlRequest); err != nil {
		uh.logger.Warnw("binding request params", "error", err)
		c.JSON(http.StatusBadRequest, "invalid payload")
		return
	}

	urlRequest.ShortURL = strings.TrimSpace(urlRequest.ShortURL)

	if !helper.IsValidURL(urlRequest.ShortURL) {
		uh.logger.Debugw("invalid url", "short_url", urlRequest.ShortURL)
		c.JSON(http.StatusBadRequest, "invalid url")
		return
	}

	long, err := uh.urlUseCase.Expand(urlRequest.ShortURL)
	if err != nil {
		uh.logger.Warnw("expanding url", "short_url", urlRequest.ShortURL, "error", err)
		c.IndentedJSON(http.StatusNotFound, "corresponding long url not found")
		return
	}

	c.Redirect(http.StatusMovedPermanently, long)
	return
}

func (uh *urlHandler) MyLinks(c *gin.Context) {
	contextUser, err := helper.GetUserFromContext(c)
	if err != nil {
		uh.logger.Warnw("fetching user from context", "error", err)
		c.JSON(http.StatusInternalServerError, "unable to get your links, please try again later")
		return
	}

	links, err := uh.urlUseCase.FetchLinksPerUser(contextUser.UID)
	if err != nil {
		uh.logger.Warnw("getting my links", "owner_id", contextUser.UID, "error", err)
		c.IndentedJSON(http.StatusNotFound, "unable to get your links, please try again later")
		return
	}

	c.IndentedJSON(http.StatusOK, links)
	return
}
