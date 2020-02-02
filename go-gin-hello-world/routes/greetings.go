package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initGreetingsRoute(greetinsgRoute *gin.RouterGroup) {
	greetinsgRoute.GET("/greetings", helloWorld)
}

func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello World...")
}
