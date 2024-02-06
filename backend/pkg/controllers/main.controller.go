package controllers

import (
	"fmt"
	_ "signalone/cmd/docs"
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

// LogAnalysisTask godoc
// @Summary Perform log analysis and generate solutions.
// @Description Perform log analysis based on the provided logs and generate solutions.
// @Tags analysis
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param logAnalysisPayload body LogAnalysisPayload true "Log analysis payload"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /issues/analysis [put]
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
		Severtiy:                  strings.ToUpper("Critical"),                             // TODO: Implement severity detection
		Title:                     "Sample issue title from 8 to 15 words. Quick summary.", // TODO: Produce title
		TimeStamp:                 time.Now(),
		IsResolved:                false,
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

// IssuesSearch godoc
// @Summary Search for issues based on specified criteria.
// @Description Search for issues based on specified criteria.
// @Tags issues
// @Accept json
// @Produce json
// @Param offset query int false "Offset for paginated results"
// @Param limit query int false "Maximum number of results per page (default: 30, max: 100)"
// @Param searchString query string false "Search string for filtering issues"
// @Param container query string false "Filter by container name"
// @Param issueSeverity query string false "Filter by issue severity"
// @Param issueType query string false "Filter by issue type"
// @Param startTimestamp query string false "Filter issues starting from this timestamp (RFC3339 format)"
// @Param endTimestamp query string false "Filter issues until this timestamp (RFC3339 format)"
// @Param isResolved query bool false "Filter resolved or unresolved issues"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /issues [get]
func (c *MainController) IssuesSearch(ctx *gin.Context) {
	issues := make([]models.IssueSearchResult, 0)
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
		fmt.Print("Error: ", err)
		startTimestamp = time.Time{}.UTC()
	}
	endTimestamp, err := time.Parse(time.RFC3339, endTimestampQuery)
	if err != nil || endTimestampQuery == "" {
		fmt.Print("Error: ", err)
		endTimestamp = time.Now().UTC()
	}

	qOpts := options.Find()
	qOpts.SetLimit(int64(limit))
	qOpts.SetSkip(int64(offset))
	qOpts.SetSort(bson.M{"timestamp": -1})
	qOpts.SetProjection(bson.M{
		"_id":           1,
		"containerName": 1,
		"severity":      1,
		"title":         1,
		"isResolved":    1,
		"timestamp":     1,
	})
	fmt.Print("startTimestamp: ", startTimestamp.UTC())
	fmt.Print("endTimestamp: ", endTimestamp.UTC())
	filter := bson.M{
		"isResolved": isResolved,
		"timestamp": bson.M{
			"$gte": startTimestamp.UTC(),
			"$lte": endTimestamp.UTC(),
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

// GetIssue godoc
// @Summary Get information about a specific issue.
// @Description Get information about a specific issue by providing its ID.
// @Tags issues
// @Accept json
// @Produce json
// @Param id path string true "ID of the issue"
// @Success 200 {object} models.Issue
// @Failure 404 {object} gin.H
// @Router /issues/{id} [get]
func (c *MainController) GetIssue(ctx *gin.Context) {
	id := ctx.Param("id")
	var issue models.Issue
	if err := c.issuesCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&issue); err != nil {
		ctx.JSON(404, gin.H{"error": "Not found"})
		return
	}
	ctx.JSON(200, issue)
}

// ResolveIssue godoc
// @Summary Resolve an issue by setting its status to resolved.
// @Description Resolve an issue by providing its ID and updating its status to resolved.
// @Tags issues
// @Accept json
// @Produce json
// @Param id path string true "ID of the issue to be resolved"
// @Success 200 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /issues/resolve/{id} [post]
// @RequestBody application/json ResolveIssueRequest true "Issue resolution request"
func (c *MainController) ResolveIssue(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.issuesCollection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"isResolved": true,
			},
		})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if res.MatchedCount == 0 {
		ctx.JSON(404, gin.H{"error": "Not found"})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "Success",
	})
}

// DeleteIssues godoc
// @Summary Delete issues based on the provided container name.
// @Description Delete issues based on the provided container name.
// @Tags issues
// @Accept json
// @Produce json
// @Param container query string true "Container name to delete issues from"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /issues [delete]
func (c *MainController) DeleteIssues(ctx *gin.Context) {
	container := ctx.Query("container")
	fmt.Print("Container: ", container)
	res, err := c.issuesCollection.DeleteMany(ctx, bson.M{"containerName": container})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "Success",
		"count":   res.DeletedCount,
	})
}

// GetContainers godoc
// @Summary Get a list of containers based on the provided user ID.
// @Description Get a list of containers based on the provided user ID.
// @Tags containers
// @Accept json
// @Produce json
// @Param userId query string true "User ID to filter containers"
// @Success 200 {array} string
// @Failure 500 {object} gin.H
// @Router /containers [get]
func (c *MainController) GetContainers(ctx *gin.Context) {
	containers := make([]string, 0)
	userId := ctx.Query("userId")
	results, err := c.issuesCollection.Distinct(ctx, "containerName", bson.M{"userId": userId})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	for _, r := range results {
		if container, ok := r.(string); ok {
			containers = append(containers, container)
		}
	}
	ctx.JSON(200, containers)
}
