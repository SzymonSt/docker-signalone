package utils

import (
	"signalone/pkg/components"
)

type InferenceEngine struct {
	LogSummarizationTaksComponent components.LogSummarizationTaskComponent
	PredictSolutionTaskComponent  components.PredictSolutionTaskComponent
	DetectLogAnomalyTaskComponent components.DetectLogAnomalyTaskComponent
}

func NewInferenceEngine(hfWrapper *HfWrapper) *InferenceEngine {
	return &InferenceEngine{
		LogSummarizationTaksComponent: hfWrapper,
		PredictSolutionTaskComponent:  nil,
		DetectLogAnomalyTaskComponent: nil,
	}
}

func (ie *InferenceEngine) LogSummarization(input string) string {
	return ie.LogSummarizationTaksComponent.Predict(input)
}

func (ie *InferenceEngine) PredictSolutions(input string) string {
	logSummaryTokenized := ie.LogSummarizationTaksComponent.Tokenize(input)
	solutionItems := ie.PredictSolutionTaskComponent.Predict(logSummaryTokenized)
	prompt := "<|user|>...TBD" // add solutionItems to prompt
	solutionsSummary := ie.LogSummarizationTaksComponent.Predict(prompt)
	return solutionsSummary
}
