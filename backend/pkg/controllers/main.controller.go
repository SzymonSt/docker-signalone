package controllers

import (
	"fmt"
	"signalone/pkg/components"
	"signalone/pkg/models"
	"signalone/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogAnalysisPayload struct {
	UserId        string `json:"userId"`
	ContainerName string `json:"containerName"`
	Logs          string `json:"logs"`
}

type GetIssuesPayload struct {
	UserId string `json:"userId"`
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
		ContainerName:             logAnalysisPayload.ContainerName,
		Severtiy:                  "Critical",                                              // TODO: Implement severity detection
		Title:                     "Sample issue title from 8 to 15 words. Quick summary.", // TODO: Produce title
		TimeStamp:                 time.Now(),
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

func (c *MainController) IssuesSearch(ctx *gin.Context) {
	var issues []models.IssueSearchResult
	var max int64
	offsetQuery := ctx.Query("offset")
	limitQuery := ctx.Query("limit")
	_ = ctx.Query("searchString")
	container := ctx.Query("container")
	issueSeverity := ctx.Query("issueSeverity")
	issueType := ctx.Query("issueType")
	startTimestampQuery := ctx.Query("startTimestamp")
	endTimestampQuery := ctx.Query("endTimestamp")
	isResolved, err := strconv.ParseBool(ctx.Query("isResolved"))
	if err != nil {
		isResolved = false
	}

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil || offsetQuery == "" {
		offset = 0
	}
	limit, err := strconv.Atoi(limitQuery)
	if err != nil || limit > 100 || limitQuery == "" {
		limit = 30
	}
	startTimestamp, err := time.Parse(time.RFC3339, startTimestampQuery)
	if err != nil {
		startTimestamp = time.Time{}
	}
	endTimestamp, err := time.Parse(time.RFC3339, endTimestampQuery)
	if err != nil || endTimestampQuery == "" {
		endTimestamp = time.Now()
	}

	qOpts := options.Find()
	qOpts.SetLimit(int64(limit))
	qOpts.SetSkip(int64(offset))
	qOpts.SetSort(bson.M{"timestamp": -1})
	qOpts.SetProjection(bson.M{
		"_id":                       1,
		"userId":                    0,
		"containerName":             1,
		"severity":                  1,
		"title":                     1,
		"logs":                      0,
		"logSummary":                0,
		"predictedSolutionsSummary": 0,
		"predictedSolutionsSources": 0,
	})
	filter := bson.M{
		"isResolved": isResolved,
		"timestamp": bson.M{
			"$gte": startTimestamp,
			"$lte": endTimestamp,
		},
	}
	if container != "" {
		filter["containerName"] = container
	}
	if issueSeverity != "" {
		filter["severity"] = issueSeverity
	}
	if issueType != "" {
		filter["type"] = issueType
	}

	cursor, err := c.issuesCollection.Find(ctx, filter, qOpts)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var issue models.IssueSearchResult
		if err := cursor.Decode(&issue); err != nil {
			continue
		}
		issues = append(issues, issue)
	}
	max, _ = c.issuesCollection.CountDocuments(ctx, filter)

	ctx.JSON(200, gin.H{
		"issues": issues,
		"max":    max,
	})
}

func (c *MainController) GetIssue(ctx *gin.Context) {
	id := ctx.Param("id")
	var issue models.Issue
	if err := c.issuesCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&issue); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"issue": issue,
	})
}

func (c *MainController) ResolveIssue(ctx *gin.Context) {
	// id := ctx.Param("id")
}

func (c *MainController) DeleteIssues(ctx *gin.Context) {
	// container := ctx.Query("container")
}
