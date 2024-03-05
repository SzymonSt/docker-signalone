package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"signalone/cmd/config"
	"signalone/pkg/models"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"go.mongodb.org/mongo-driver/bson"
)

func CalculateNewCounter(currentScore int32, newScore int32, counter int32) int32 {
	return counter + (newScore - currentScore)
}

func GenerateFilter(fields bson.M, operator string) bson.M {
	conditions := make([]bson.M, 0, len(fields))

	for field, value := range fields {
		conditions = append(conditions, bson.M{field: value})
	}

	return bson.M{operator: conditions}
}

func CallPredictionAgentService(jsonData []byte) (analysisResponse models.IssueAnalysis, err error) {
	var cfg = config.GetInstance()

	issueAnalysisReq, err := http.NewRequest("POST", cfg.PredicitonAgentServiceUrl+"/run_analysis", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	issueAnalysisReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(issueAnalysisReq)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rawAnalysisResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(rawAnalysisResponse, &analysisResponse)
	if err != nil {
		return
	}
	return
}

func CompareLogs(incomingLogTails []string, currentIssuesLogTails []string) (isNewIssue bool) {
	const LogSimilarityThreshold = 0.8

	isNewIssue = true
	sdm := metrics.NewSorensenDice()
	sdm.CaseSensitive = false
	sdm.NgramSize = 3
	for _, incomingLogTail := range incomingLogTails {
		for _, currentIssueLogTail := range currentIssuesLogTails {
			similarity := strutil.Similarity(incomingLogTail, currentIssueLogTail, sdm)
			if similarity >= LogSimilarityThreshold {
				isNewIssue = false
				return
			}
		}
	}
	return
}

func FilterForRelevantLogs(logs []string) (relevantLogs []string) {
	//Classes are absractions of different types of logs as different types of issues
	//have different log structures
	// Class 0 = Error or Warning message
	// Class 1 = Exception with stack trace
	issueClassZeroRegex := `(?i)(abort|blocked|corrupt|crash|critical|deadlock|
		denied|deprecated|deprecating|err|error|fatal|forbidden|
		freeze|hang|illegal|invalid|missing|panic|refused|rejected|
		timeout|unauthorized|unsupported|warn|warning)`
	issueClassOneRegex := `(?i)(exception|stacktrace|traceback|uncaught|unhandled)`

	compiledClassZeroRegex := regexp.MustCompile(issueClassZeroRegex)
	compiledClassOneRegex := regexp.MustCompile(issueClassOneRegex)

	for logIndex, log := range logs {
		if matched := compiledClassOneRegex.MatchString(log); matched {
			relevantLogs = append(relevantLogs, logs[logIndex])
			//Add the next and previous log to the relevant logs if stack trace is found
			//To be improved
			if logIndex+1 < len(logs) {
				relevantLogs = append(relevantLogs, logs[logIndex+1])
			}
			if logIndex-1 >= 0 {
				relevantLogs = append(relevantLogs, logs[logIndex-1])
			}
			return
		}

		if matched := compiledClassZeroRegex.MatchString(log); matched {
			relevantLogs = append(relevantLogs, logs[logIndex])
		}
	}

	return
}

func GenerateBearerToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(b)
	return token, nil
}
