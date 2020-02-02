package config

import (
	"os"
	"strconv"
)

// Conf - Main Struct for all Configurations
type Conf struct {
	GinServer
	*Gin
}

// GinServer - Stores Server details
type GinServer struct {
	bindIPAddress string
	bindPort      int
}

// Init - Initialize config properties
func Init() *Conf {
	var conf Conf

	conf.GinServer.bindIPAddress = os.Getenv("GO_GIN_IP_ADDR")
	conf.GinServer.bindPort, _ = strconv.Atoi(os.Getenv("GO_GIN_IP_PORT"))

	return &conf
}
