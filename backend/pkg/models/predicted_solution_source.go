package models

type PredictedSolutionSource struct {
	Url            string `json:"url:string_value"`
	Title          string `json:"title:string_value"`
	Description    string `json:"description:string_value"`
	FeaturedAnswer string `json:"featuredAnswer:string_value"`
}
