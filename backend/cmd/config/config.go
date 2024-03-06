package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort   string `mapstructure:"SERVER_PORT"`
	Mode         string `mapstructure:"MODE"`
	CertFilePath string `mapstructure:"CERT_FILE_PATH"`
	KeyFilePath  string `mapstructure:"KEY_FILE_PATH"`

	//GitHub Data
	GithubClientId     string `mapstructure:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `mapstructure:"GITHUB_CLIENT_SECRET"`

	//Google Data
	GoogleClientId     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`

	//SignalOne Data
	SignalOneSecret string `mapstructure:"SIGNAL_ONE_SECRET"`

	//Inference Engine API
	PredicitonAgentServiceUrl string `mapstructure:"PREDICTION_AGENT_SERVICE_URL"`
	InferenceApiUrl           string `mapstructure:"INFERENCE_API_URL"`
	InferenceApiKey           string `mapstructure:"INFERENCE_API_KEY"`
	InferenceBaseModel        string `mapstructure:"BASE_MODEL_NAME"`

	//Tokenized Solution for Prediction Database Details
	SolutionDbHost         string `mapstructure:"SOLUTION_DB_HOST"`
	SolutionCollectionName string `mapstructure:"SOLUTION_COLLECTION_NAME"`

	//Application Database Details
	ApplicationDbUrl                string `mapstructure:"APPLICATION_DB_URL"`
	ApplicationDbName               string `mapstructure:"APPLICATION_DB_NAME"`
	ApplicationIssuesCollectionName string `mapstructure:"APPLICATION_ISSUES_COLLECTION_NAME"`
	ApplicationUsersCollectionName  string `mapstructure:"APPLICATION_USERS_COLLECTION_NAME"`

	//Saved Analysis Database Details
	SavedAnalysisDbUrl          string `mapstructure:"SAVED_ANALYSIS_DB_URL"`
	SavedAnalysisDbName         string `mapstructure:"SAVED_ANALYSIS_DB_NAME"`
	SavedAnalysisCollectionName string `mapstructure:"SAVED_ANALYSIS_COLLECTION_NAME"`
}

var (
	once   sync.Once
	config *Config
)

func GetInstance() *Config {
	once.Do(func() {
		viper.SetConfigName(".default")
		viper.AddConfigPath(".")
		viper.SetConfigType("env")
		viper.AutomaticEnv()

		err := viper.ReadInConfig()

		if err != nil {
			return
		}

		config = &Config{}
		err = viper.Unmarshal(config)

		if err != nil {
			return
		}
	})

	return config
}
