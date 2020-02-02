package config

import (
	"fmt"

	"github.com/g33kzone/go-gin-hello-world/routes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Gin - Gin Server
type Gin struct {
	GinEngine    *gin.Engine
	ServerConfig GinServer
}

// StartServer - Initialize Gin Server
func (gin Gin) StartServer() {
	logrus.Infof("Attempting to start Gin server - %s:%d", gin.ServerConfig.bindIPAddress, gin.ServerConfig.bindPort)

	router := routes.GinRoutes{GinEngine: gin.GinEngine}
	router.InitializeRoutes()

	err := gin.GinEngine.Run(fmt.Sprintf("%s:%d", gin.ServerConfig.bindIPAddress, gin.ServerConfig.bindPort))
	if err != nil {
		panic(fmt.Sprintf("Failed to start Gin server %s", err.Error()))
	}
}
