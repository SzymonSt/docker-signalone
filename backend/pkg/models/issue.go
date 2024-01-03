package models

type IssueSolutionPredictionSolutionSource struct {
	Title string
	Url   string
}

type Issue struct {
	Id                        string                                  `json:"id"`
	UserId                    string                                  `json:"user_id"`
	Logs                      string                                  `json:"logs"`
	LogSummary                string                                  `json:"log_summary"`
	PredictedSolutionsSummary string                                  `json:"predicted_solutions_summary"`
	PredictedSolutionsSources []IssueSolutionPredictionSolutionSource `json:"issue_predicted_solutions_sources"`
}
