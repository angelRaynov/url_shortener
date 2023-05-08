package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"url_shortener/internal/model"
	"url_shortener/internal/url"
)

type urlHandler struct {
	urlUseCase url.ShortExpander
}

func NewURLHandler(uc url.ShortExpander) url.ShortenExpandHandler {
	return &urlHandler{urlUseCase: uc}
}

func (uh *urlHandler) ShortenURL(c *gin.Context)  {
	var urlRequest model.ShortenRequest

	if err := c.BindJSON(&urlRequest); err != nil {
		log.Printf("binding request params:%v", err)
		c.JSON(http.StatusBadRequest, "invalid payload")
		return
	}

	short := uh.urlUseCase.Shorten(urlRequest.LongURL)

	c.IndentedJSON(http.StatusCreated, model.ShortenResponse{
		Link:    short,
		LongURL: urlRequest.LongURL,
	})
	return

}

func (uh *urlHandler) ExpandURL(c *gin.Context)  {
	shortenedURL := c.Param("shortened")

	if shortenedURL != "" {
		shortenedURL = strings.TrimSpace(shortenedURL)
		//redirect
			long, err := uh.urlUseCase.Expand("http://localhost:1234/"+shortenedURL)
			if err != nil {
				err = fmt.Errorf("expanding url:%w", err)
				log.Printf("%s", err)
				c.IndentedJSON(http.StatusNotFound,"corresponding long url not found")
				return
			}

			c.Redirect(http.StatusMovedPermanently, long)
			return
	}
	c.JSON(http.StatusBadRequest,"unable to process request")


}
