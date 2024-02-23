package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"signalone/cmd/config"
	"signalone/pkg/models"

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
