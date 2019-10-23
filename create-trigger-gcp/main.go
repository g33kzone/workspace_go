package main

import (
	"fmt"
	"odj-deliver-cloudbuild/config"
	"odj-deliver-cloudbuild/constants"
	"odj-deliver-cloudbuild/db"
	"odj-deliver-cloudbuild/docs"

	//"odj-deliver-cloudbuild/docs"
	"odj-deliver-cloudbuild/route"

	_ "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/lib/pq"
)

// swagger comments

// @title One developer journey template
// @version 0.1
// @description This API is used to Create Trigger and store Build Details
// @contact.name Babita Gurjar
// @contact.url https://infosys.com
// @contact.email babita.gurjar@infosys.com
// @BasePath /
func main() {
	//initialize config
	conf := &config.Config{Title: "ODJ"}
	conf.InitConfig()

	//initialize db
	db := &db.Db{}
	db.InitDBConnection(conf.PostgresConfig, constants.DBConnectionAttempts)

	//set up routes
	route.InitGinEngine(conf.GinConfig, db)

	//set up swagger
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%v", conf.SwaggerConfig.Host, conf.SwaggerConfig.Port)

}
