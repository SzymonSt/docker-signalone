package components

import "signalone/pkg/models"

type SolutionPredictionResult struct {
	SolutionSummary string
	SolutionSources []models.IssueSolutionPredictionSolutionSource
}

type LogSummarizationTaskComponent interface {
	Predict(input string) string
}

type PredictSolutionTaskComponent interface {
	Predict(input []float32) []models.PredictedSolutionSource
	Tokenize(input string) []float32
}

type DetectLogAnomalyTaskComponent interface {
	Predict(input string) bool
}
