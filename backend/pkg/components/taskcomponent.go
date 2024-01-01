package components

type LogSummarizationTaskComponent interface {
	Predict(input string) string
	Tokenize(input string) []float32
}

type PredictSolutionTaskComponent interface {
	Predict(input []float32) []string
}

type DetectLogAnomalyTaskComponent interface {
	Predict(input string) bool
}
