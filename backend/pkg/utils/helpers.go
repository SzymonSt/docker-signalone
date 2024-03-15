package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"signalone/cmd/config"
	"signalone/pkg/models"
	"unicode"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
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
	const LogSimilarityThreshold = 0.6

	isNewIssue = true
	sdm := metrics.NewSorensenDice()
	sdm.CaseSensitive = false
	sdm.NgramSize = 3
	for _, incomingLogTail := range incomingLogTails {
		for _, currentIssueLogTail := range currentIssuesLogTails {
			similarity := strutil.Similarity(incomingLogTail, currentIssueLogTail, sdm)
			fmt.Print("-----------------------\n")
			fmt.Print("Incoming Log Tail: ", incomingLogTail, "\n")
			fmt.Print("Current Issue Log Tail: ", currentIssueLogTail, "\n")
			fmt.Print("Similarity: ", similarity, "\n")
			fmt.Print("Is New Issue: ", (similarity <= LogSimilarityThreshold), "\n")
			fmt.Print("-----------------------\n")
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

func ComparePasswordHashes(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func PasswordValidation(password string) bool {
	if !(len(password) >= 8 && len(password) <= 50) {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}
	return (hasUpper && hasLower && hasDigit)
}
