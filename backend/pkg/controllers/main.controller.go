package controllers

import (
	"fmt"
	"signalone/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LogPayload struct {
	userId      string
	containerId string
	logs        string
}

type MainController struct {
	iEngine *utils.InferenceEngine
}

func NewMainController(iEngine *utils.InferenceEngine) *MainController {
	return &MainController{
		iEngine: iEngine,
	}
}

func (c *MainController) LogAnalysisTask(ctx *gin.Context) {
	summarizationTaskPromptTemplate := "<|user|>Summarize these logs and generate a single paragraph summary of what is happening in these logs in high technical detail: \n%s\n</s>\n<|assistant|>"
	issueId := uuid.New().String()
	var logPayload LogPayload
	if err := ctx.ShouldBindJSON(&logPayload); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	generatedSummary := c.iEngine.LogSummarization(fmt.Sprintf(summarizationTaskPromptTemplate, logPayload.logs))
	// IF PRO: use inference engine to generate solutions with sources
	proposed_solutions := make([]string, 0)
	//put new issue into database (user_id, issue_id, container_id, logs, summary, proposed_solutions)
	ctx.JSON(200, gin.H{
		"message": "Success",
		"issueId": issueId,
	})
}
