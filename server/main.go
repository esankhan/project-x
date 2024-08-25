package main

import (
	"fmt"
	"os"

	"github.com/esankhan/project-x/urlmanagement"
	"github.com/esankhan/project-x/usersmanagement"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)



func main() {
	err := godotenv.Load()
if err != nil {
	fmt.Println("Error loading .env file",err)
}
  router := gin.Default();

  v1 := router.Group("/v1")
  {
	v1.POST("/register", usersmanagement.RegisterHandler())
	v1.POST("/login", usersmanagement.LoginHandler())
	v1.POST("/short", urlmanagement.ShortenUrlHandler())
	v1.POST("/resolve", urlmanagement.ResolveUrlHandler())
  }

  port:= ":"+os.Getenv("APP_PORT")

  fmt.Println("server running on", port)

  router.Run(port)

}