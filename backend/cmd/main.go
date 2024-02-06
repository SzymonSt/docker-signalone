package main

import (
	"context"
	"net/http"
	"signalone/cmd/config"
	"signalone/pkg/controllers"
	"signalone/pkg/routers"
	"signalone/pkg/utils"

	_ "signalone/cmd/docs" // Import the generated docs package

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var InferenceHyperParameters = map[string]interface{}{
	"temperature":    0.7,
	"top_k":          20,
	"top_p":          0.9,
	"do_sample":      true,
	"max_new_tokens": 160,
}

var RAGHyperParameters = map[string]interface{}{
	"limit": 3,
}

//	@title			SignalOne API
//	@version		1.0
//	@description	API for SignalOne application
//	@host			localhost:8080
//	@BasePath		/api
func main() {
	var (
		server = gin.Default()
	)
	cfg := config.New()
	if cfg == nil {
		panic("critical: unable to load config")
	}

	// if cfg.Mode == "local" {
	// 	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// }

	ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	appDbClient, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(cfg.ApplicationDbUrl),
	)
	if err != nil {
		panic(err)
	}
	issuesCollectionClient := appDbClient.Database(cfg.ApplicationDbName).Collection(cfg.ApplicationIssuesCollectionName)
	usersCollectionClient := appDbClient.Database(cfg.ApplicationDbName).Collection(cfg.ApplicationUsersCollectionName)

	savedAnalysisDbClient, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(cfg.SavedAnalysisDbUrl),
	)
	if err != nil {
		panic(err)
	}
	savedAnalysisCollectionClient := savedAnalysisDbClient.Database(cfg.SavedAnalysisDbName).Collection(cfg.SavedAnalysisCollectionName)

	hfwrapper := utils.NewHfWrapper(
		cfg.InferenceApiUrl,
		"models",
		cfg.InferenceBaseModel,
		cfg.InferenceApiKey,
		InferenceHyperParameters["temperature"].(float64),
		InferenceHyperParameters["top_k"].(int),
		InferenceHyperParameters["top_p"].(float64),
		InferenceHyperParameters["do_sample"].(bool),
		InferenceHyperParameters["max_new_tokens"].(int),
	)

	ragwrapper := utils.NewRagWrapper(
		cfg.SolutionDbHost,
		hfwrapper,
		cfg.SolutionCollectionName,
		uint64(RAGHyperParameters["limit"].(int)),
	)

	inferenceEngine := utils.NewInferenceEngine(
		hfwrapper,
		ragwrapper,
	)

	mainController := controllers.NewMainController(
		inferenceEngine,
		issuesCollectionClient,
		usersCollectionClient,
		savedAnalysisCollectionClient,
	)

	//authController TBD
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	//	@BasePath	/api
	router := server.Group("/api")
	router.GET("/healthz", func(ctx *gin.Context) {
		message := "signal api is up and running, operational subsystems: {}"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	// Add the following route for Swagger info
	//router.GET("/swagger-info", getSwaggerInfo)

	routeController := routers.NewMainRouter(mainController)
	routeController.RegisterRoutes(router)

	server.Run(":" + cfg.ServerPort)
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /hi [get]
func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	c.JSON(http.StatusOK, res)
}
