package controllers

import (
	"fmt"
	"signalone/pkg/components"
	"signalone/pkg/models"
	"signalone/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogAnalysisPayload struct {
	UserId      string `json:"userId"`
	ContainerId string `json:"containerId"`
	Logs        string `json:"logs"`
}

type MainController struct {
	iEngine                 *utils.InferenceEngine
	issuesCollection        *mongo.Collection
	usersCollection         *mongo.Collection
	analysisStoreCollection *mongo.Collection
}

func NewMainController(iEngine *utils.InferenceEngine,
	issuesCollection *mongo.Collection, usersCollection *mongo.Collection,
	analysisStoreCollection *mongo.Collection) *MainController {
	return &MainController{
		iEngine:                 iEngine,
		issuesCollection:        issuesCollection,
		usersCollection:         usersCollection,
		analysisStoreCollection: analysisStoreCollection,
	}
}

func (c *MainController) LogAnalysisTask(ctx *gin.Context) {
	var generatedSummary string
	var proposedSolutions components.SolutionPredictionResult
	var user models.User
	summarizationTaskPromptTemplate := `<|user|>
	Summarize these logs and generate a single paragraph summary of what is happening in these logs in high technical detail: %s</s>
	<|assistant|>`

	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		ctx.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	bearerToken = strings.TrimPrefix(bearerToken, "Bearer ")
	var logAnalysisPayload LogAnalysisPayload
	if err := ctx.ShouldBindJSON(&logAnalysisPayload); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userResult := c.usersCollection.FindOne(ctx, bson.M{"userId": logAnalysisPayload.UserId})
	err := userResult.Decode(&user)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//Commented our for testing purposes
	// if user.AgentBearerToken != bearerToken {
	// 	ctx.JSON(401, gin.H{
	// 		"message": "Unauthorized",
	// 	})
	// 	return
	// }
	issueId := uuid.New().String()
	generatedSummary = c.iEngine.LogSummarization(
		fmt.Sprintf(summarizationTaskPromptTemplate, logAnalysisPayload.Logs),
	)
	generatedSummary = strings.Split(generatedSummary, "<|assistant|>")[1]
	proposedSolutions = c.iEngine.PredictSolutions(generatedSummary)
	if !user.IsPro {
		c.analysisStoreCollection.InsertOne(ctx, models.SavedAnalysis{
			Logs:       logAnalysisPayload.Logs,
			LogSummary: generatedSummary,
		})
	}
	fmt.Println("Generated Summary: ", generatedSummary)
	fmt.Println("Proposed Solutions: ", proposedSolutions)
	fmt.Printf("Soultion Sources: %+v\n", proposedSolutions.SolutionSources)

	c.issuesCollection.InsertOne(ctx, models.Issue{
		Id:                        issueId,
		UserId:                    logAnalysisPayload.UserId,
		Logs:                      logAnalysisPayload.Logs,
		LogSummary:                generatedSummary,
		PredictedSolutionsSummary: proposedSolutions.SolutionSummary,
		PredictedSolutionsSources: proposedSolutions.SolutionSources,
	})
	ctx.JSON(200, gin.H{
		"message": "Success",
		"issueId": issueId,
	})
}
