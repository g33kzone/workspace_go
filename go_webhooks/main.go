package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var router *gin.Engine

func main() {
	router = gin.Default()
	initializeRoutes()
	router.Run("0.0.0.0:8080")
}

func initializeRoutes() {
	router.GET("/api", fetchHelloWorld)
}

func fetchHelloWorld(c *gin.Context) {
	log.Info("GET API call received...")

	c.JSON(http.StatusOK, "Hello World !!")
}
