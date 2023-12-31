package utils

import (
	"signalone/pkg/components"
)

type InferenceEngine struct {
	logSummarizationTaksComponent components.LogSummarizationTaskComponent
	predcitSolutionTaskComponent  components.PredictSolutionTaskComponent
	detectLogAnomalyTaskComponent components.DetectLogAnomalyTaskComponent
}

func NewInferenceEngine(hfWrapper *HfWrapper) *InferenceEngine {
	return &InferenceEngine{
		logSummarizationTaksComponent: hfWrapper,
	}
}
