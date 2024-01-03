package controllers

import (
	"fmt"
	"signalone/pkg/components"
	"signalone/pkg/models"
	"signalone/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogAnalysisPayload struct {
	userId      string
	isPro       bool // JUST FOR TESTING PURPOSES
	containerId string
	logs        string
}

type MainController struct {
	iEngine                 *utils.InferenceEngine
	applicationCollection   *mongo.Collection
	analysisStoreCollection *mongo.Collection
}

func NewMainController(iEngine *utils.InferenceEngine,
	applicationCollection *mongo.Collection,
	analysisStoreCollection *mongo.Collection) *MainController {
	return &MainController{
		iEngine:                 iEngine,
		applicationCollection:   applicationCollection,
		analysisStoreCollection: analysisStoreCollection,
	}
}

func (c *MainController) LogAnalysisTask(ctx *gin.Context) {
	var generatedSummary string
	var proposedSolutions components.SolutionPredictionResult
	summarizationTaskPromptTemplate := `<|user|>
	Summarize these logs and generate a single paragraph summary of what is happening in these logs in high technical detail: %s</s>
	<|assistant|>`
	issueId := uuid.New().String()
	var logAnalysisPayload LogAnalysisPayload
	if err := ctx.ShouldBindJSON(&logAnalysisPayload); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	generatedSummary = c.iEngine.LogSummarization(fmt.Sprintf(summarizationTaskPromptTemplate, logAnalysisPayload.logs))
	if logAnalysisPayload.isPro {
		proposedSolutions = c.iEngine.PredictSolutions(generatedSummary)
	} else {
		proposedSolutions = components.SolutionPredictionResult{
			SolutionSummary: "",
			SolutionSources: []models.IssueSolutionPredictionSolutionSource{},
		}
		c.analysisStoreCollection.InsertOne(ctx, models.SavedAnalysis{
			Logs:       logAnalysisPayload.logs,
			LogSummary: generatedSummary,
		})
	}

	c.applicationCollection.InsertOne(ctx, models.Issue{
		Id:                        issueId,
		UserId:                    logAnalysisPayload.userId,
		Logs:                      logAnalysisPayload.logs,
		LogSummary:                generatedSummary,
		PredictedSolutionsSummary: proposedSolutions.SolutionSummary,
		PredictedSolutionsSources: proposedSolutions.SolutionSources,
	})
	ctx.JSON(200, gin.H{
		"message": "Success",
		"issueId": issueId,
	})
}
