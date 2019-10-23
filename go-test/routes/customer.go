package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitCustomerRoutes -
func initCustomerRoutes(customerRoute *gin.RouterGroup) {
	customerRoute.GET("/customers", helloWorld)
}

func helloWorld(c *gin.Context) {
	fmt.Println("Hello World")
	c.JSON(http.StatusOK, "Hello World")
}
