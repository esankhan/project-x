package urlmanagement

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type resolveUrlRequest struct {
	ShortUrl string `json:"shortUrl"`
}

func ResolveUrlHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req resolveUrlRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
			})
			return
		}
		if(req.ShortUrl == ""){
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
			})
			return
		}

		
		url,dbSource,err := ResolveUrl(req.ShortUrl)
		if(err != nil){
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "URL not found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "URL resolved successfully",
			"url": url,
			"dbSource":dbSource,
		})
	}
}