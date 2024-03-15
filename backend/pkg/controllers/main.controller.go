package controllers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"signalone/cmd/config"
	_ "signalone/docs"
	"signalone/pkg/models"
	"signalone/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	e "github.com/jordan-wright/email"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogAnalysisPayload struct {
	UserId        string `json:"userId"`
	ContainerName string `json:"containerName"`
	ContainerId   string `json:"containerId"`
	Severity      string `json:"severity"`
	Logs          string `json:"logs"`
}

type Log struct {
	Logs []string `bson:"logs"`
}

type EmailClientConfig struct {
	AuthData    smtp.Auth
	HostAddress string
	From        string
	TlsConfig   *tls.Config
}

type MainController struct {
	issuesCollection        *mongo.Collection
	usersCollection         *mongo.Collection
	analysisStoreCollection *mongo.Collection
	waitlistCollection      *mongo.Collection
	emailClientData         EmailClientConfig
}

const ACCESS_TOKEN_EXPIRATION_TIME = time.Minute * 10
const REFRESH_TOKEN_EXPIRATION_TIME = time.Hour * 24

func NewMainController(issuesCollection *mongo.Collection,
	usersCollection *mongo.Collection,
	analysisStoreCollection *mongo.Collection,
	waitlistCollection *mongo.Collection,
	emailClientData EmailClientConfig) *MainController {
	return &MainController{
		issuesCollection:        issuesCollection,
		usersCollection:         usersCollection,
		analysisStoreCollection: analysisStoreCollection,
		waitlistCollection:      waitlistCollection,
		emailClientData:         emailClientData,
	}
}

func (c *MainController) ContactHandler(ctx *gin.Context) {
	var emailReqBody models.Email

	err := ctx.ShouldBindJSON(&emailReqBody)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailObj := e.NewEmail()
	emailObj.From = c.emailClientData.From
	emailObj.To = []string{"contact@signaloneai.com"}
	emailObj.Subject = fmt.Sprintf("[CONTACT] %s", emailReqBody.MessageTitle)
	emailObj.Text = []byte(fmt.Sprintf("From: %s \nMessage: %s", emailReqBody.Email, emailReqBody.MessageContent))
	err = emailObj.SendWithStartTLS(c.emailClientData.HostAddress, c.emailClientData.AuthData, c.emailClientData.TlsConfig)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error sending the mail"})
		return
	}

	resEmailObj := e.NewEmail()
	resEmailObj.From = c.emailClientData.From
	resEmailObj.To = []string{emailReqBody.Email}
	resEmailObj.Subject = "Thank you for contacting us"
	resEmailObj.HTML = []byte(`<img alt="Signal0ne" title="Signal0ne Logo" width="196px" height="57px" src="https://signaloneai.com/online-assets/Signal0ne.jpg"
	style="margin-top: 40px;"><h1 style="color: black">Hello,</h1> <p style="color: black">Thank you for contacting us.</p> <p style="color: black">We will get back to you as soon as possible.</p><br><p style="color: black; margin-bottom: 0; margin-top: 4px;">Best regards,</p><p style="color: black; font-family: consolas; font-size: 15px; font-weight: bold; margin-top: 6px;";>Signal0ne Team</p>`)
	err = resEmailObj.SendWithStartTLS(c.emailClientData.HostAddress, c.emailClientData.AuthData, c.emailClientData.TlsConfig)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error sending the mail"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Email has been sent successfully",
	})
}

func (c *MainController) WaitlistHandler(ctx *gin.Context) {
	var waitlistEntry models.WaitlistEntry

	err := ctx.ShouldBindJSON(&waitlistEntry)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	waitlistQueryResult := c.waitlistCollection.FindOne(ctx, bson.M{"email": waitlistEntry.Email})
	if waitlistQueryResult.Err() != mongo.ErrNoDocuments {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists in the waitlist"})
		return
	}

	emailObj := e.NewEmail()
	emailObj.From = c.emailClientData.From
	emailObj.To = []string{waitlistEntry.Email}
	emailObj.Subject = "Thank you for joining the waitlist!"
	emailObj.HTML = []byte(utils.WaitlistEntryConfirmationEmail)
	err = emailObj.SendWithStartTLS(c.emailClientData.HostAddress, c.emailClientData.AuthData, c.emailClientData.TlsConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error sending the mail"})
		return
	}

	_, err = c.waitlistCollection.InsertOne(ctx, waitlistEntry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /issues/analysis [put]
func (c *MainController) LogAnalysisTask(ctx *gin.Context) {
	var user models.User
	var analysisResponse models.IssueAnalysis

	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		ctx.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	var logAnalysisPayload LogAnalysisPayload
	if err := ctx.ShouldBindJSON(&logAnalysisPayload); err != nil {
		fmt.Printf("Error: %s", err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userResult := c.usersCollection.FindOne(ctx, bson.M{"userId": logAnalysisPayload.UserId})
	err := userResult.Decode(&user)
	if err != nil {
		fmt.Printf("Error: %s", err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	issueId := uuid.New().String()
	go func() {
		var issueLogs = make([]string, 0)
		var issueLog Log
		var isNewIssue = true

		formattedAnalysisLogs := strings.Split(logAnalysisPayload.Logs, "\n")
		formattedAnalysisRelevantLogs := utils.FilterForRelevantLogs(formattedAnalysisLogs)

		qOpts := options.Find()
		qOpts.Projection = bson.M{"logs": 1}

		cursor, err := c.issuesCollection.Find(ctx, bson.M{
			"containerId": logAnalysisPayload.ContainerId,
			"isResolved":  false,
		}, qOpts)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			if err := cursor.Decode(&issueLog); err != nil {
				continue
			}
			issueLogs = append(issueLogs, issueLog.Logs...)
		}

		//Compare logs with previous logs and if they are similar enough, don't call the prediction agent
		if len(issueLogs) > 0 {
			isNewIssue = utils.CompareLogs(formattedAnalysisRelevantLogs, issueLogs)
			if !isNewIssue {
				return
			}
		}

		initialInsertResult, _ := c.issuesCollection.InsertOne(ctx, models.Issue{
			Id:                        issueId,
			UserId:                    "",
			ContainerName:             logAnalysisPayload.ContainerName,
			ContainerId:               logAnalysisPayload.ContainerId,
			Score:                     0,
			Severity:                  logAnalysisPayload.Severity,
			Title:                     analysisResponse.Title,
			TimeStamp:                 time.Now(),
			IsResolved:                false,
			Logs:                      formattedAnalysisLogs,
			LogSummary:                "",
			PredictedSolutionsSummary: "",
			PredictedSolutionsSources: []string{},
		})

		data := map[string]string{"logs": strings.Join(formattedAnalysisRelevantLogs, "\n")}
		jsonData, _ := json.Marshal(data)
		analysisResponse, err = utils.CallPredictionAgentService(jsonData)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		if !user.IsPro {
			c.analysisStoreCollection.InsertOne(ctx, models.SavedAnalysis{
				Logs:       logAnalysisPayload.Logs,
				LogSummary: analysisResponse.LogSummary,
			})
		}

		_, err = c.issuesCollection.UpdateOne(ctx,
			bson.M{
				"_id":         initialInsertResult.InsertedID,
				"containerId": logAnalysisPayload.ContainerId,
			},
			bson.M{"$set": bson.M{
				"userId":                         logAnalysisPayload.UserId,
				"title":                          analysisResponse.Title,
				"timestamp":                      time.Now(),
				"predictedSolutionsSummary":      analysisResponse.PredictedSolutions,
				"issuePredictedSolutionsSources": analysisResponse.Sources,
				"logSummary":                     analysisResponse.LogSummary,
			},
			})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}()

	ctx.JSON(200, gin.H{
		"message": "Acknowledged",
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
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /issues [get]
func (c *MainController) IssuesSearch(ctx *gin.Context) {
	var max int64
	issues := make([]models.IssueSearchResult, 0)

	userId, err := getUserIdFromToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	container := ctx.Query("container")
	endTimestampQuery := ctx.Query("endTimestamp")
	issueSeverity := ctx.Query("issueSeverity")
	issueType := ctx.Query("issueType")
	limitQuery := ctx.Query("limit")
	offsetQuery := ctx.Query("offset")
	startTimestampQuery := ctx.Query("startTimestamp")
	_ = ctx.Query("searchString")

	isResolved, err := strconv.ParseBool(ctx.Query("isResolved"))
	if err != nil {
		isResolved = true
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

	filter := bson.M{
		"userId": userId,
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

	if !isResolved {
		filter["isResolved"] = isResolved
	}

	cursor, err := c.issuesCollection.Find(ctx, filter, qOpts)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	ctx.JSON(http.StatusOK, gin.H{
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
// @Failure 404 {object} map[string]any
// @Router /issues/{id} [get]
func (c *MainController) GetIssue(ctx *gin.Context) {
	var issue models.Issue
	id := ctx.Param("id")

	if err := c.issuesCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&issue); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	ctx.JSON(http.StatusOK, issue)
}

func (c *MainController) RateIssue(ctx *gin.Context) {
	var issue models.Issue
	var issueRateReq models.IssueRateRequest
	var user models.User

	userId, err := getUserIdFromToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	//TODO: Remove hardcoded userId
	// userId = "4c78e05c-2f83-4e6e-b4c1-8721618a1c89"

	err = ctx.ShouldBindJSON(&issueRateReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if *issueRateReq.Score != -1 && *issueRateReq.Score != 0 && *issueRateReq.Score != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Score must be one of: -1, 0, 1"})
		return
	}

	userResult := c.usersCollection.FindOne(ctx, bson.M{"userId": userId})

	err = userResult.Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	issueConditions := bson.M{
		"_id":    id,
		"userId": userId,
	}

	filter := utils.GenerateFilter(issueConditions, "$and")
	issueResult := c.issuesCollection.FindOne(ctx, filter)

	err = issueResult.Decode(&issue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var currentIssueScore = issue.Score

	if currentIssueScore == *issueRateReq.Score {
		ctx.JSON(http.StatusOK, gin.H{"message": "Issue already rated with the same score"})
		return
	}

	updatedIssueResult, err := c.issuesCollection.UpdateOne(ctx,
		filter,
		bson.M{
			"$set": bson.M{
				"score": issueRateReq.Score,
			},
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if updatedIssueResult.MatchedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Issue cannot be found"})
		return
	}

	counter := user.Counter
	counter = utils.CalculateNewCounter(currentIssueScore, *issueRateReq.Score, counter)

	updatedUserResult, err := c.usersCollection.UpdateOne(ctx,
		bson.M{"userId": userId},
		bson.M{
			"$set": bson.M{
				"counter": counter,
			},
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if updatedUserResult.MatchedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User cannot be found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

func (c *MainController) RegenerateSolution(ctx *gin.Context) {
	var analysisResponse models.IssueAnalysis
	var issue models.Issue
	var user models.User

	userId, err := getUserIdFromToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	issueResult := c.issuesCollection.FindOne(ctx, bson.M{"_id": id, "userId": userId})

	err = issueResult.Decode(&issue)
	if err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResult := c.usersCollection.FindOne(ctx, bson.M{"userId": userId})

	err = userResult.Decode(&user)
	if err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var formattedAnalysisRelevantLogs = utils.FilterForRelevantLogs(issue.Logs)
	data := map[string]string{"logs": strings.Join(formattedAnalysisRelevantLogs, "\n")}
	jsonData, _ := json.Marshal(data)

	analysisResponse, err = utils.CallPredictionAgentService(jsonData)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	if !user.IsPro {
		c.analysisStoreCollection.InsertOne(ctx, models.SavedAnalysis{
			Logs:       strings.Join(issue.Logs, "\n"),
			LogSummary: analysisResponse.LogSummary,
		})
	}
	_, err = c.issuesCollection.UpdateOne(ctx, bson.M{"_id": id, "userId": userId}, bson.M{"$set": bson.M{
		"title":                          analysisResponse.Title,
		"timestamp":                      time.Now(),
		"predictedSolutionsSummary":      analysisResponse.PredictedSolutions,
		"issuePredictedSolutionsSources": analysisResponse.Sources,
		"logSummary":                     analysisResponse.LogSummary,
	}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	issueResult = c.issuesCollection.FindOne(ctx, bson.M{"_id": id, "userId": userId})

	err = issueResult.Decode(&issue)
	if err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}

	ctx.JSON(http.StatusOK, issue)
}

// ResolveIssue godoc
// @Summary Mark issue as resolved/unresolved.
// @Description Resolve an issue by providing its ID and resolve state of the issue.
// @Tags issues
// @Accept json
// @Produce json
// @Param id path string true "ID of the issue to be resolved"
// @Success 200 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /issues/{id}/resolve [put]
// @RequestBody application/json isResolved boolean
func (c *MainController) ResolveIssue(ctx *gin.Context) {
	var requestData models.IssueResolveRequest

	userId, err := getUserIdFromToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	issueResult, err := c.issuesCollection.UpdateOne(ctx, bson.M{"_id": id, "userId": userId}, bson.M{"$set": bson.M{"isResolved": *requestData.IsResolved}})
	if issueResult.MatchedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
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
// @Success 200 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /issues [delete]
func (c *MainController) DeleteIssues(ctx *gin.Context) {
	containerId := ctx.Param("containerId")
	res, err := c.issuesCollection.DeleteMany(ctx, bson.M{"containerId": containerId})
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
// @Failure 500 {object} map[string]any
// @Router /containers [get]
func (c *MainController) GetContainers(ctx *gin.Context) {
	containers := make([]string, 0)

	userId, err := getUserIdFromToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	results, err := c.issuesCollection.Distinct(ctx, "containerName", bson.M{"userId": userId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, result := range results {
		if container, ok := result.(string); ok {
			containers = append(containers, container)
		}
	}
	ctx.JSON(http.StatusOK, containers)
}

// Auth Handlers
func (c *MainController) LoginWithGithubHandler(ctx *gin.Context) {
	var requestData models.GithubTokenRequest
	var user models.User

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userData, err = getGithubData(requestData.Code)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userResult := c.usersCollection.FindOne(ctx, bson.M{"userId": strconv.Itoa(userData.Id)})

	err = userResult.Decode(&user)

	if err != nil && err.Error() != mongo.ErrNoDocuments.Error() {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
		user = models.User{
			UserId:           strconv.Itoa(userData.Id),
			UserName:         userData.Login,
			IsPro:            false,
			AgentBearerToken: "",
			Counter:          0,
			Type:             "github",
		}

		_, err = c.usersCollection.InsertOne(ctx, user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	accessTokenString, err := createToken(user.UserId, user.UserName, "access")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't make authentication token"})
		return
	}

	refreshTokenString, err := createToken(user.UserId, user.UserName, "refresh")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't make authentication token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Success",
		"accessToken":  accessTokenString,
		"expiresIn":    int64(ACCESS_TOKEN_EXPIRATION_TIME) / int64(time.Second),
		"refreshToken": refreshTokenString,
	})
}

func (c *MainController) LoginWithGoogleHandler(ctx *gin.Context) {
	var requestData models.GoogleTokenRequest
	var user models.User

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := validateGoogleJWT(requestData.IdToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userResult := c.usersCollection.FindOne(ctx, bson.M{"userId": claims.Subject})

	err = userResult.Decode(&user)

	if err != nil && err.Error() != mongo.ErrNoDocuments.Error() {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
		user = models.User{
			UserId:           claims.Subject,
			UserName:         claims.FirstName,
			IsPro:            false,
			AgentBearerToken: "",
			Counter:          0,
			Type:             "google",
		}

		_, err = c.usersCollection.InsertOne(ctx, user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	accessTokenString, err := createToken(user.UserId, user.UserName, "access")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't make authentication token"})
		return
	}

	refreshTokenString, err := createToken(user.UserId, user.UserName, "refresh")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't make authentication token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Success",
		"accessToken":  accessTokenString,
		"expiresIn":    int64(ACCESS_TOKEN_EXPIRATION_TIME) / int64(time.Second),
		"refreshToken": refreshTokenString,
	})
}

func (c *MainController) LoginHandler(ctx *gin.Context) {
	var loginData models.SignalAccountPayload
	var user models.User

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	loginData.Email = strings.ToLower(loginData.Email)

	userResult := c.usersCollection.FindOne(ctx, bson.M{"userName": loginData.Email, "type": "signalone"})

	err := userResult.Decode(&user)
	if err == mongo.ErrNoDocuments {
		ctx.JSON(http.StatusUnauthorized, gin.H{"descriptionKey": "INVALID_CREDENTIALS"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	if !utils.ComparePasswordHashes(user.PasswordHash, loginData.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"descriptionKey": "INVALID_CREDENTIALS"})
		return
	}

	if !user.EmailConfirmed {
		ctx.JSON(http.StatusUnauthorized, gin.H{"descriptionKey": "ACCOUNT_NOT_ACTIVE"})
		return
	}

	accessTokenString, err := createToken(user.UserId, user.UserName, "access")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	refreshTokenString, err := createToken(user.UserId, user.UserName, "refresh")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Success",
		"accessToken":  accessTokenString,
		"expiresIn":    int64(ACCESS_TOKEN_EXPIRATION_TIME) / int64(time.Second),
		"refreshToken": refreshTokenString,
	})

}

func (c *MainController) RegisterHandler(ctx *gin.Context) {
	var loginData models.SignalAccountPayload
	var user models.User

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	loginData.Email = strings.ToLower(loginData.Email)

	if !utils.PasswordValidation(loginData.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "INVALID_PASSWORD"})
		return
	}

	userResult := c.usersCollection.FindOne(ctx, bson.M{"userName": loginData.Email, "type": "signalone"})

	err := userResult.Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		ctx.JSON(http.StatusInternalServerError, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	if err == nil {
		if !user.EmailConfirmed {
			ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "DUPLICATE_NOT_CONFIRMED_USER_EMAIL"})
			return
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "EMAIL_ALREADY_IN_USE"})
			return
		}
	}

	hashedPassword, err := utils.HashPassword(loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	userId := uuid.New().String()
	confirmationToken, err := createToken(userId, loginData.Email, "refresh")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	confirmationLink := fmt.Sprintf("https://signaloneai.com/email-verification?email=%s&verificationCode=%s", loginData.Email, confirmationToken)

	emailObj := e.NewEmail()
	emailObj.From = c.emailClientData.From
	emailObj.To = []string{loginData.Email}
	emailObj.Subject = "Confirm your email address"
	emailObj.HTML = []byte(fmt.Sprintf(`<img alt="Signal0ne" title="Signal0ne Logo" width="196px" height="57px" src="https://signaloneai.com/online-assets/Signal0ne.jpg"
	style="margin-top: 40px;">
	<h1 style="color: black">Hello,</h1> 
	<p style="color: black">Welcome to <span style="font-family: consolas;">Signal0ne!</span> We're excited you're joining us.</p>
	<p style="color: black">Ready to get started? First, verify your email address by clicking the button below: </p>
	<a href="%s" target="_blank"><button style="background-color: #3f51b5; border: none; border-radius: 6px; color: #ffffff; cursor: pointer; padding-bottom: 8px; padding-top: 8px; padding-left: 16px; padding-right: 16px;">Confirm</button></a><br>
	<p>Or you can click the following link: <a href="%s" target="_blank" style="color: #3f51b5">%s</a></p>
	<p style="color: black; margin-bottom: 0; margin-top: 4px;">Best regards,</p>
	<p style="color: black; font-family: consolas; font-size: 15px; font-weight: bold; margin-top: 6px;";>Signal0ne Team</p>`, confirmationLink, confirmationLink, confirmationLink))

	err = emailObj.SendWithStartTLS(c.emailClientData.HostAddress, c.emailClientData.AuthData, c.emailClientData.TlsConfig)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "INVALID_EMAIL"})
	}

	user = models.User{
		UserId:                userId,
		UserName:              loginData.Email,
		PasswordHash:          hashedPassword,
		IsPro:                 false,
		AgentBearerToken:      "",
		Counter:               0,
		Type:                  "signalone",
		EmailConfirmed:        false,
		EmailConfirmationCode: confirmationToken,
	}
	_, err = c.usersCollection.InsertOne(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})

}

func (c *MainController) VerifyEmail(ctx *gin.Context) {
	var user models.User
	var verificationData models.EmailConfirmationPayload

	if err := ctx.ShouldBindJSON(&verificationData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	_, err := VerifyToken(verificationData.ConfirmationToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "INVALID_VERIFICATION_TOKEN"})
		return
	}

	userResult := c.usersCollection.FindOne(ctx, bson.M{"emailConfirmationCode": verificationData.ConfirmationToken})

	err = userResult.Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "INVALID_VERIFICATION_CODE"})
		return
	}

	_, err = c.usersCollection.UpdateOne(ctx,
		bson.M{"userId": user.UserId},
		bson.M{"$set": bson.M{
			"emailConfirmed":        true,
			"emailConfirmationCode": "",
		},
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *MainController) ResendConfirmationEmail(ctx *gin.Context) {
	var user models.User
	var verificationData struct {
		Email string `json:"email"`
	}

	if err := ctx.ShouldBindJSON(&verificationData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	userResult := c.usersCollection.FindOne(ctx, bson.M{"userName": verificationData.Email})

	err := userResult.Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "EMAIL_NOT_FOUND"})
		return
	}

	if user.EmailConfirmed {
		ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "EMAIL_ALREADY_IN_USE"})
		return
	}

	confirmationToken, err := createToken(user.UserId, user.UserName, "refresh")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	_, err = c.usersCollection.UpdateOne(ctx,
		bson.M{"userId": user.UserId},
		bson.M{"$set": bson.M{
			"emailConfirmationCode": confirmationToken,
		},
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"descriptionKey": "ERROR_OCCURED"})
		return
	}

	confirmationLink := fmt.Sprintf("https://signaloneai.com/email-verification?email=%s&verificationCode=%s", verificationData.Email, confirmationToken)

	emailObj := e.NewEmail()
	emailObj.From = c.emailClientData.From
	emailObj.To = []string{verificationData.Email}
	emailObj.Subject = "Confirm your email address"
	emailObj.HTML = []byte(fmt.Sprintf(`<img alt="Signal0ne" title="Signal0ne Logo" width="196px" height="57px" src="https://signaloneai.com/online-assets/Signal0ne.jpg"
	style="margin-top: 40px;">
	<h1 style="color: black">Hello,</h1> 
	<p style="color: black">Welcome to <span style="font-family: consolas;">Signal0ne!</span> We're excited you're joining us.</p>
	<p style="color: black">Ready to get started? First, verify your email address by clicking the button below: </p>
	<a href="%s" target="_blank"><button style="background-color: #3f51b5; border: none; border-radius: 6px; color: #ffffff; cursor: pointer; padding-bottom: 8px; padding-top: 8px; padding-left: 16px; padding-right: 16px;">Confirm</button></a><br>
	<p>Or you can click the following link: <a href="%s" target="_blank" style="color: #3f51b5">%s</a></p>
	<p style="color: black; margin-bottom: 0; margin-top: 4px;">Best regards,</p>
	<p style="color: black; font-family: consolas; font-size: 15px; font-weight: bold; margin-top: 6px;";>Signal0ne Team</p>`, confirmationLink, confirmationLink, confirmationLink))
	err = emailObj.SendWithStartTLS(c.emailClientData.HostAddress, c.emailClientData.AuthData, c.emailClientData.TlsConfig)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"descriptionKey": "INVALID_EMAIL"})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *MainController) RefreshTokenHandler(ctx *gin.Context) {
	var cfg = config.GetInstance()
	var claims = &models.JWTClaimsWithUserData{}
	var data models.RefreshTokenRequest
	var SECRET_KEY = []byte(cfg.SignalOneSecret)

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.ParseWithClaims(data.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	accessTokenString, err := createToken(claims.Id, claims.UserName, "access")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't make authentication token"})
		return
	}

	refreshTokenString, err := createToken(claims.Id, claims.UserName, "refresh")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't make authentication token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Success",
		"accessToken":  accessTokenString,
		"expiresIn":    int64(ACCESS_TOKEN_EXPIRATION_TIME) / int64(time.Second),
		"refreshToken": refreshTokenString,
	})
}

func (c *MainController) AuthenticateAgent(ctx *gin.Context) {
	var user models.User

	userId, err := getUserIdFromToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	result := c.usersCollection.FindOne(ctx, bson.M{"userId": userId})
	err = result.Decode(&user)
	if err != nil {
		fmt.Printf("Error: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.AgentBearerToken != "" {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"token":   user.AgentBearerToken,
		})
		return
	}

	token, err := createToken(userId, user.UserName, "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = c.usersCollection.UpdateOne(ctx,
		bson.M{"userId": userId},
		bson.M{"$set": bson.M{
			"agentBearerToken": token,
		},
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"token":   token,
	})
}

func (c *MainController) CheckAgentAuthorization(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	var token = strings.TrimPrefix(authHeader, "Bearer ")

	err := c.VerifyAgentToken(ctx, token)
	if err != nil {
		fmt.Printf("Error: %s", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.Next()
}

func getUserIdFromToken(ctx *gin.Context) (string, error) {
	bearerToken := ctx.GetHeader("Authorization")

	jwtToken := strings.TrimPrefix(bearerToken, "Bearer ")

	userId, err := VerifyToken(jwtToken)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func createToken(id string, userName string, tokenType string) (string, error) {
	var cfg = config.GetInstance()
	var expTime time.Duration
	var SECRET_KEY = []byte(cfg.SignalOneSecret)

	if tokenType == "refresh" {
		expTime = REFRESH_TOKEN_EXPIRATION_TIME
	} else if tokenType == "access" {
		expTime = ACCESS_TOKEN_EXPIRATION_TIME
	} else {
		expTime = time.Second * 0
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp":      time.Now().Add(expTime).Unix(),
			"id":       id,
			"userName": userName,
		})

	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getGithubData(code string) (models.GithubUserData, error) {
	var cfg = config.GetInstance()
	var githubData = models.GithubUserData{}
	var githubJWTData = models.GithubTokenResponse{}
	var httpClient = &http.Client{}

	ghJWTReqBody := map[string]string{
		"client_id":     cfg.GithubClientId,
		"client_secret": cfg.GithubClientSecret,
		"code":          code,
	}

	jsonData, _ := json.Marshal(ghJWTReqBody)

	ghJWTReq, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(jsonData))
	if err != nil {
		return models.GithubUserData{}, err
	}

	ghJWTReq.Header.Set("Accept", "application/json")
	ghJWTReq.Header.Set("Content-Type", "application/json")

	ghJWTResp, err := httpClient.Do(ghJWTReq)
	if err != nil {
		return models.GithubUserData{}, err
	}

	ghJWTRespBody, err := io.ReadAll(ghJWTResp.Body)
	if err != nil {
		return models.GithubUserData{}, err
	}

	err = json.Unmarshal(ghJWTRespBody, &githubJWTData)
	if err != nil {
		return models.GithubUserData{}, err
	}

	ghUserDataReq, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return models.GithubUserData{}, err
	}

	ghUserDataReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", githubJWTData.AccessToken))

	ghUserDataResp, err := httpClient.Do(ghUserDataReq)
	if err != nil {
		return models.GithubUserData{}, err
	}

	ghUserDataRespBody, err := io.ReadAll(ghUserDataResp.Body)
	if err != nil {
		return models.GithubUserData{}, err
	}

	err = json.Unmarshal(ghUserDataRespBody, &githubData)
	if err != nil {
		return models.GithubUserData{}, err
	}

	return githubData, nil
}

func getGooglePublicKey(keyId string) (string, error) {
	var googleData = map[string]string{}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, &googleData)
	if err != nil {
		return "", err
	}

	key, ok := googleData[keyId]
	if !ok {
		return "", errors.New("key not found")
	}

	return key, nil
}

func validateGoogleJWT(tokenString string) (models.GoogleClaims, error) {
	var cfg = config.GetInstance()
	var claimsStruct = models.GoogleClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}

			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}

			return key, nil
		},
	)

	if err != nil {
		return models.GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*models.GoogleClaims)
	if !ok {
		return models.GoogleClaims{}, errors.New("invalid claims")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return models.GoogleClaims{}, errors.New("iss is invalid")
	}

	audienceToCheck := cfg.GoogleClientId
	found := false

	for _, audience := range claims.Audience {
		if audience == audienceToCheck {
			found = true
			break
		}
	}

	if !found {
		return models.GoogleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Unix() < time.Now().UTC().Unix() {
		return models.GoogleClaims{}, errors.New("jwt is expired")
	}

	return *claims, nil
}

func VerifyToken(tokenString string) (string, error) {
	var cfg = config.GetInstance()
	var claims = &models.JWTClaimsWithUserData{}
	var SECRET_KEY = []byte(cfg.SignalOneSecret)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Id, nil
}

func (c *MainController) VerifyAgentToken(ctx *gin.Context, token string) (err error) {
	var user models.User
	var cfg = config.GetInstance()
	var claims = &models.JWTClaimsWithUserData{}
	var SECRET_KEY = []byte(cfg.SignalOneSecret)

	parser := jwt.NewParser(
		jwt.WithoutClaimsValidation(),
	)

	_, err = parser.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	if err != nil {
		return
	}

	err = c.usersCollection.FindOne(ctx, bson.M{"userId": claims.Id}).Decode(&user)
	if err != nil {
		return
	}

	if user.AgentBearerToken == "" || user.AgentBearerToken != token {
		err = errors.New("unauthorized")
		return
	}

	return
}
