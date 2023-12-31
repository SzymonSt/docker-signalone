package components

type LogSummarizationTaskComponent interface {
	Predict(input string) string
}

type PredictSolutionTaskComponent interface {
	Predict(input string) []string
}

type DetectLogAnomalyTaskComponent interface {
	Predict(input string) bool
}
