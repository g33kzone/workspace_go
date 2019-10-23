package main

import (
	"go-test/config"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

func main() {
	conf := config.Init()
	conf.DBGorm = config.DBGormConn(conf.DBConfig)

	ginServer := config.Gin{ServerConfig: conf.GinServer,
		GinEngine: gin.Default()}
	ginServer.StartServer()
}
