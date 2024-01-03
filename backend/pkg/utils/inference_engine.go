package utils

import (
	"fmt"
	"signalone/pkg/components"
	"signalone/pkg/models"
)

type InferenceEngine struct {
	LogSummarizationTaksComponent components.LogSummarizationTaskComponent
	PredictSolutionTaskComponent  components.PredictSolutionTaskComponent
	DetectLogAnomalyTaskComponent components.DetectLogAnomalyTaskComponent
}

func NewInferenceEngine(hfWrapper *HfWrapper, ragwrapper *RagWrapper) *InferenceEngine {
	return &InferenceEngine{
		LogSummarizationTaksComponent: hfWrapper,
		PredictSolutionTaskComponent:  ragwrapper,
		DetectLogAnomalyTaskComponent: nil,
	}
}

func (ie *InferenceEngine) LogSummarization(input string) string {
	return ie.LogSummarizationTaksComponent.Predict(input)
}

func (ie *InferenceEngine) PredictSolutions(input string) components.SolutionPredictionResult {
	var predictedSolution components.SolutionPredictionResult
	var solutionFeaturedAnswers []string
	prompt := `<|user|>
	Based on ther provided summary of the issue: %s and the list of proposed solutions: %v, provide a single paragraph summary of the proposed solution in high technical detail.</s>
	<|assistant|>`
	logSummaryTokenized := ie.PredictSolutionTaskComponent.Tokenize(input)
	solutionItems := ie.PredictSolutionTaskComponent.Predict(logSummaryTokenized)
	for _, solutionItem := range solutionItems {
		predictedSolution.SolutionSources = append(predictedSolution.SolutionSources, models.IssueSolutionPredictionSolutionSource{
			Title: solutionItem.Title,
			Url:   solutionItem.Url,
		})
		solutionFeaturedAnswers = append(solutionFeaturedAnswers, solutionItem.FeaturedAnswer)
	}
	solutionsSummary := ie.LogSummarizationTaksComponent.Predict(fmt.Sprintf(prompt, input, solutionFeaturedAnswers))
	predictedSolution.SolutionSummary = solutionsSummary

	return predictedSolution
}
