package routers

import (
	"signalone/pkg/controllers"
	middlewares "signalone/pkg/middleware"

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
	authorizationRouterGroup.POST("/login-with-github", mr.mainController.LoginWithGithubHandler)
	authorizationRouterGroup.POST("/login-with-google", mr.mainController.LoginWithGoogleHandler)
	authorizationRouterGroup.POST("/token/refresh", mr.mainController.RefreshTokenHandler)
	authorizationRouterGroup.POST("/user/login", func(c *gin.Context) {})
	authorizationRouterGroup.PUT("/user/register", func(c *gin.Context) {})

	userRouterGroup := rg.Group("/user", middlewares.CheckAuthorization)
	{
		userRouterGroup.POST("/agent/authenticate", func(c *gin.Context) {})
		userRouterGroup.GET("/containers", mr.mainController.GetContainers)
		userRouterGroup.GET("/issues", mr.mainController.IssuesSearch)
		userRouterGroup.GET("/issues/:id", mr.mainController.GetIssue)
		userRouterGroup.POST("/issues/:id", mr.mainController.ResolveIssue)
		userRouterGroup.PUT("/issues/:id/score", mr.mainController.RateIssue)
		userRouterGroup.GET("/settings", func(c *gin.Context) {})
		userRouterGroup.POST("/settings", func(c *gin.Context) {})
	}

	agentRouterGroup := rg.Group("/agent")
	agentRouterGroup.DELETE("/issues", mr.mainController.DeleteIssues)
	agentRouterGroup.PUT("/issues/analysis", mr.mainController.LogAnalysisTask)
}
