package route

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"odj-deliver-cloudbuild/config"
	"odj-deliver-cloudbuild/db"
)

//InitGinEngine is used to initialize and then start gin web server
func InitGinEngine(ginConfig config.GinConfig, db *db.Db) {
	log.Println(fmt.Sprintf("Starting web server on -> %s:%v", ginConfig.Host, ginConfig.Port))

	router := setUpRouter(ginConfig, db)

	setGinMode(strings.ToUpper(ginConfig.Mode))
	startEngine(router, ginConfig)
}

//setUpRouter sets up routes, handle end points to be exposed
func setUpRouter(ginConfig config.GinConfig, db *db.Db) *gin.Engine {
	router := gin.Default()

	routerGroup := router.Group("odj-deliver")                                     //common base path for all routes
	routerGroup.Use(setCors())                                                     // set cors headers
	routerGroup.Use(globalRecover)                                                 // recover on panic with json and a 500 status
	routerGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // adding in the swagger UI
	routerGroup.POST("/create-trigger", func(c *gin.Context) {
		handleClouddBuild(c, db)
	})
	routerGroup.POST("/webhook/github", func(c *gin.Context) {
		handleWebHooksGithub(c, db)
	})
	routerGroup.POST("/webhook/bitbucket", func(c *gin.Context) {
		handleWebHooksBitbucket(c, db)
	})

	return router
}

//setGinMode sets gin mode as per defined in config
func setGinMode(mode string) {
	switch strings.ToUpper(mode) {
	case "DEBUG":
		gin.SetMode(gin.DebugMode)
	case "TEST":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)

	}
}

//startEngine starts routes via gin engine
func startEngine(router *gin.Engine, ginConfig config.GinConfig) {
	err := router.Run(fmt.Sprintf("%s:%v", ginConfig.Host, ginConfig.Port))

	if err != nil {
		log.Println(fmt.Sprintf("CRASH! Failed to start web server : %v", err.Error()))
		panic(err)
	}
}
