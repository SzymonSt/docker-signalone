package models

import "time"

type IssueSolutionPredictionSolutionSource struct {
	Title string `json:"title" bson:"title"`
	Url   string `json:"url" bson:"url"`
}

type IssueSearchResult struct {
	Id            string `json:"id" bson:"_id"`
	ContainerName string `json:"containerName" bson:"containerName"`
	Title         string `json:"title" bson:"title"`
	Severity      string `json:"severity" bson:"severity"`
}

type Issue struct {
	Id                        string                                  `json:"id" bson:"_id"`
	UserId                    string                                  `json:"userId" bson:"userId"`
	ContainerName             string                                  `json:"containerName" bson:"containerName"`
	Severtiy                  string                                  `json:"severity" bson:"severity"`
	Logs                      string                                  `json:"logs" bson:"logs"`
	Title                     string                                  `json:"title" bson:"title"`
	TimeStamp                 time.Time                               `json:"timestamp" bson:"timestamp"`
	LogSummary                string                                  `json:"logSummary" bson:"logSummary"`
	PredictedSolutionsSummary string                                  `json:"predictedSolutionsSummary" bson:"predictedSolutionsSummary"`
	PredictedSolutionsSources []IssueSolutionPredictionSolutionSource `json:"issuePredictedSolutionsSources" bson:"issuePredictedSolutionsSources"`
}
