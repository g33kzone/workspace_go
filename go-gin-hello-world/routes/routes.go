package routes

import "github.com/gin-gonic/gin"

// GinRoutes -
type GinRoutes struct {
	GinEngine *gin.Engine
}

// InitializeRoutes - Initialize all Routes
func (routes GinRoutes) InitializeRoutes() {

	routerGroup := routes.GinEngine.Group("v1")
	initGreetingsRoute(routerGroup)
}
