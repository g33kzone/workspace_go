package config

import (
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
)

// Conf -
type Conf struct {
	GinServer
	DBConfig
	DBGorm *gorm.DB
	*Gin
}

// GinServer -
type GinServer struct {
	bindIPAddress string
	bindPort      int
	mode          string
}

// DBConfig -
type DBConfig struct {
	database string
	server   string
	port     int
	username string
	password string
}

// Init -
func Init() *Conf {
	var conf Conf

	conf.GinServer.bindIPAddress = os.Getenv("GO_GIN_IP_ADDR")
	conf.GinServer.bindPort, _ = strconv.Atoi(os.Getenv("GO_GIN_IP_PORT"))

	conf.DBConfig.database = os.Getenv("POSTGRES_DATABASE")
	conf.DBConfig.server = os.Getenv("POSTGRES_SERVER")
	conf.DBConfig.port, _ = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	conf.DBConfig.username = os.Getenv("POSTGRES_USER")
	conf.DBConfig.password = os.Getenv("POSTGRES_PASS")

	return &conf
}
