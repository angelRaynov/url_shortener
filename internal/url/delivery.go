package url

import (
	"github.com/gin-gonic/gin"
)

type ShortenExpandHandler interface {
	ShortenURL(c *gin.Context)
	ExpandURL(c *gin.Context)
}
