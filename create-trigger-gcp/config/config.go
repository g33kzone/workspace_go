package config

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

// Config stores configuration details for gin, swagger and postgres db
type Config struct {
	Title          string
	GinConfig      GinConfig
	SwaggerConfig  SwaggerConfig
	PostgresConfig PostgresConfig
	RepoDetails    RepoDetails
}

// GinConfig stores configuration details for gin
type GinConfig struct {
	Mode string
	Host string
	Port uint16
}

// SwaggerConfig stores configuration details for swagger
type SwaggerConfig struct {
	Host string
	Port uint16
}

// PostgresConfig stores configuration details for postgres
type PostgresConfig struct {
	Server              string
	Port                uint16
	Database            string
	User                string
	Password            string
	DBRetryAttempts     uint8
	DBReadTimeout       uint32
	DBWriteTimeout      uint32
	DBMaxConnections    uint8
	DBConnectionMaxLife uint32
}

//RepoDetails store Repository details
type RepoDetails struct {
	Branch   string
	UserName string
	FileName string
}

//InitConfig initialises static config object
func (config *Config) InitConfig() {
	log.Println("Loading config file")

	if _, err := toml.DecodeFile("config/config.toml", &config); err != nil {
		log.Println(fmt.Sprintf("CRASH! Failed to load the config! : %v", err.Error()))
		panic(err)
	}
}
