package config

import (
	"fmt"
	"go-test/routes"

	"github.com/gin-gonic/gin"
)

// Gin -
type Gin struct {
	GinEngine    *gin.Engine
	ServerConfig GinServer
}

// StartServer -
func (gin Gin) StartServer() {
	fmt.Println(fmt.Sprintf("Attempting to start Gin server - %s:%d",
		gin.ServerConfig.bindIPAddress, gin.ServerConfig.bindPort))

	router := routes.GinRoutes{GinEngine: gin.GinEngine}
	router.InitializeRoutes()

	err := gin.GinEngine.Run(fmt.Sprintf("%s:%d",
		gin.ServerConfig.bindIPAddress, gin.ServerConfig.bindPort))
	if err != nil {
		panic(fmt.Sprintf("Failed to start Gin server %s", err.Error()))
	}
}
