package routes

import (
	"github.com/gin-gonic/gin"
)

// GinRoutes -
type GinRoutes struct {
	GinEngine *gin.Engine
}

// InitializeRoutes -
func (routes GinRoutes) InitializeRoutes() {
	routerGroup := routes.GinEngine.Group("test-bank")
	initCustomerRoutes(routerGroup)
}
