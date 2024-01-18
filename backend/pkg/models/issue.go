package models

type IssueSolutionPredictionSolutionSource struct {
	Title string
	Url   string
}

type Issue struct {
	Id                        string                                  `json:"id"`
	UserId                    string                                  `json:"userId"`
	Logs                      string                                  `json:"logs"`
	LogSummary                string                                  `json:"logSummary"`
	PredictedSolutionsSummary string                                  `json:"predictedSolutionsSummary"`
	PredictedSolutionsSources []IssueSolutionPredictionSolutionSource `json:"issuePredictedSolutionsSources"`
}
