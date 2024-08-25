package usersmanagement

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type registerRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`	
	Password string `json:"password"`
}


func RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest

		if err:= c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
			})
			return
		}
		if(FindUser(req.Email)){
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "use another email",
			})
			return
		}
		if(req.Email == "" || req.Username == "" || req.Password == ""){
				c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
			})
			return
		}
		req.Password = HashPassword(req.Password)
		id := InsertUser(req)
		c.JSON(http.StatusOK, gin.H{
			"message": "Register success with email",
			"email": req.Email,
			"id":id,
		})
	}
	
}



