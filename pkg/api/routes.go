package api

import (
	"github.com/gin-gonic/gin"
)

func LoginRoutes(router *gin.Engine) {

	// Grouping routes
	apiGroup := router.Group("/api")
    {
    	apiGroup.POST("/v1/tasks/", UserAuthMiddleware(), CreateTaskController)
    }
}
