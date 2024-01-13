package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	//Inference Engine API
	InferenceApiUrl      string `mapstructure:"INFERENCE_API_URL"`
	HuggingFaceApiKey    string `mapstructure:"HUGGINGFACE_API_KEY"`
	HuggingFaceBaseModel string `mapstructure:"BASE_MODEL_NAME"`

	//Tokenized Solution for Prediction Database Details
	SolutionDbHost         string `mapstructure:"SOLUTION_DB_HOST"`
	SolutionDbName         string `mapstructure:"SOLUTION_DB_NAME"`
	SolutionCollectionName string `mapstructure:"SOLUTION_COLLECTION_NAME"`

	//Application Database Details
	ApplicationDbUrl          string `mapstructure:"APPLICATION_DB_URL"`
	ApplicationDbName         string `mapstructure:"APPLICATION_DB_NAME"`
	ApplicationCollectionName string `mapstructure:"APPLICATION_COLLECTION_NAME"`

	//Saved Analysis Database Details
	SavedAnalysisDbUrl          string `mapstructure:"SAVED_ANALYSIS_DB_URL"`
	SavedAnalysisDbName         string `mapstructure:"SAVED_ANALYSIS_DB_NAME"`
	SavedAnalysisCollectionName string `mapstructure:"SAVED_ANALYSIS_COLLECTION_NAME"`
}

func New() (config *Config) {
	viper.SetConfigName(".default")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil
	}
	return config
}
