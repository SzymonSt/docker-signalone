package models

type IssueAnalysis struct {
	Title              string   `json:"title"`
	LogSummary         string   `json:"logsummary"`
	PredictedSolutions string   `json:"predictedSolutions"`
	Sources            []string `json:"sources"`
}
