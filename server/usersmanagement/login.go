package usersmanagement

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req loginRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
			})
			return
		}

		if(req.Email == "" || req.Password == ""){
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
			})
			return
		}
		
		user := FindUserByEmail(req.Email)

		// create redis connection
		
		
		if(user.Email == ""){
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User not found",
			})
			return
		}
 matchPassword := ComparePassword(user.Password, req.Password)
 if(!matchPassword){
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email address or password is incorrect",
		})
		return
	}
		c.JSON(http.StatusOK, gin.H{
			"message": "Logged in Successfully",
			"email": req.Email,
			"password": req.Password,
		})
	}
}