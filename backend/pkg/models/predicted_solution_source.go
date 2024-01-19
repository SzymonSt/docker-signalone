package models

type PredictedSolutionSource struct {
	Url            string `json:"url"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	FeaturedAnswer string `json:"featuredAnswer"`
}
