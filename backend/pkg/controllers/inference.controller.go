package controllers

type InferenceController struct {
	hfWrapper *HfWrapper
}

func NewInferenceController(hfWrapper *HfWrapper) *InferenceController {
	return &InferenceController{
		hfWrapper: hfWrapper,
	}
}

func (ic *InferenceController) LogSummarizationTask(logs string) string {
	return ic.hfWrapper.Predict(logs)
}
