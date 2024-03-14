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
	rg.POST("/contact", mr.mainController.ContactHandler)

	authorizationRouterGroup := rg.Group("/auth")
	authorizationRouterGroup.POST("/email-confirmation", mr.mainController.VerifyEmail)
	authorizationRouterGroup.POST("/email-confirmation-link-resend", mr.mainController.ResendConfirmationEmail)
	authorizationRouterGroup.POST("/login", mr.mainController.LoginHandler)
	authorizationRouterGroup.POST("/login-with-github", mr.mainController.LoginWithGithubHandler)
	authorizationRouterGroup.POST("/login-with-google", mr.mainController.LoginWithGoogleHandler)
	authorizationRouterGroup.POST("/register", mr.mainController.RegisterHandler)
	authorizationRouterGroup.POST("/token/refresh", mr.mainController.RefreshTokenHandler)

	userRouterGroup := rg.Group("/user", middlewares.CheckAuthorization)
	{
		userRouterGroup.POST("/agent/authenticate", mr.mainController.AuthenticateAgent)
		userRouterGroup.GET("/containers", mr.mainController.GetContainers)
		userRouterGroup.GET("/issues", mr.mainController.IssuesSearch)
		userRouterGroup.GET("/issues/:id", mr.mainController.GetIssue)
		userRouterGroup.PUT("/issues/:id/regenerate", mr.mainController.RegenerateSolution)
		userRouterGroup.PUT("/issues/:id/resolve", mr.mainController.ResolveIssue)
		userRouterGroup.PUT("/issues/:id/score", mr.mainController.RateIssue)
		userRouterGroup.GET("/settings", func(c *gin.Context) {})
		userRouterGroup.POST("/settings", func(c *gin.Context) {})
	}

	agentRouterGroup := rg.Group("/agent", mr.mainController.CheckAgentAuthorization)
	{
		agentRouterGroup.DELETE("/issues/:containerId", mr.mainController.DeleteIssues)
		agentRouterGroup.PUT("/issues/analysis", mr.mainController.LogAnalysisTask)
	}
}
