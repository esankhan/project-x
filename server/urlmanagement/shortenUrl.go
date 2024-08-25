package urlmanagement

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type shortenUrlRequest struct {
	Url string `json:"url"`
	Email string `json:"email"`
}

func ShortenUrlHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
var req shortenUrlRequest

if err := c.BindJSON(&req); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Invalid request",
	})
	return
}

if(req.Url == "" || req.Email == ""){
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Invalid request",
	})
	return
}

// check if the user exists

id, err2 := FindUserByEmail(req.Email);

if(err2 != nil || id == -1){
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "User not found",
	})
	return
}

// is the url valid
if(!IsValidUrl(req.Url)){
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Invalid URL",
	})
	return
}
// shorten the url
shortUrl := ShortenUrl(req.Url)

// save the url to the database
SaveUrl(id, req.Url, shortUrl)

c.JSON(http.StatusOK, gin.H{
	"message": "URL shortened successfully",
	"shortUrl": shortUrl,
	"email":req.Email,
	"url":req.Url,
})

	}
}