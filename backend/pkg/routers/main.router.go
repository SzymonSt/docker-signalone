package routers

import (
	"signalone/pkg/controllers"

	"github.com/gin-gonic/gin"
)

type MainRouter struct {
	mainController *controllers.MainController
}

func NewMainRouter(mainController *controllers.MainController) *MainRouter {
	return &MainRouter{
		mainController: mainController,
	}
}

func (mr *MainRouter) RegisterRoutes(rg *gin.RouterGroup) {
	authorizationRouterGroup := rg.Group("/api/auth")
	authorizationRouterGroup.POST("/user/login", func(c *gin.Context) {})
	authorizationRouterGroup.POST("/user/register", func(c *gin.Context) {})
	authorizationRouterGroup.POST("/agent/authenticate", func(c *gin.Context) {})

	userRouterGroup := rg.Group("/api/user")
	userRouterGroup.GET("/settings", func(c *gin.Context) {})
	userRouterGroup.POST("/settings", func(c *gin.Context) {})
	userRouterGroup.GET("/issues", func(c *gin.Context) {})
	userRouterGroup.GET("/issues/:id", func(c *gin.Context) {})

	agentRouterGroup := rg.Group("/api/agent")
	agentRouterGroup.PUT("/issues/analysis", mr.mainController.LogAnalysisTask)

}
