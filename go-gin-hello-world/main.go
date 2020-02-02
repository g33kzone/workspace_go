package main

import (
	"github.com/g33kzone/go-gin-hello-world/config"
	"github.com/gin-gonic/gin"
)

func main() {

	conf := config.Init()

	ginServer := config.Gin{GinEngine: gin.Default(), ServerConfig: conf.GinServer}

	ginServer.StartServer()
}
