package components

type LogSummarizationTaskComponent interface {
	Predict(input string) string
}

type PredictSolutionTaskComponent interface {
	Predict(input []float32) []string
	Tokenize(input string) []float32
}

type DetectLogAnomalyTaskComponent interface {
	Predict(input string) bool
}
