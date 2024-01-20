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
	authorizationRouterGroup := rg.Group("/auth")
	authorizationRouterGroup.POST("/user/login", func(c *gin.Context) {})
	authorizationRouterGroup.PUT("/user/register", func(c *gin.Context) {})

	userRouterGroup := rg.Group("/user")
	userRouterGroup.GET("/settings", func(c *gin.Context) {})
	userRouterGroup.POST("/settings", func(c *gin.Context) {})
	userRouterGroup.GET("/issues", mr.mainController.IssuesSearch)
	userRouterGroup.GET("/issues/:id", mr.mainController.GetIssue)
	userRouterGroup.POST("/issues/:id", mr.mainController.ResolveIssue)
	userRouterGroup.POST("/agent/authenticate", func(c *gin.Context) {})

	agentRouterGroup := rg.Group("/agent")
	agentRouterGroup.DELETE("/issues", mr.mainController.DeleteIssues)
	agentRouterGroup.PUT("/issues/analysis", mr.mainController.LogAnalysisTask)

}
