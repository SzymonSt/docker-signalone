package helper

import (
	"github.com/spf13/viper"
)

type Config struct {
	HuggingFaceApiKey    string `mapstructure:"HUGGINGFACE_API_KEY"`
	HuggingFaceBaseModel string `mapstructure:"BASE_MODEL_NAME"`
	DbHost               string `mapstructure:"DB_HOST"`
	DbPassword           string `mapstructure:"DB_PASSWORD"`
	DbUser               string `mapstructure:"DB_USER"`
	DbName               string `mapstructure:"DB_NAME"`
}

func New() (config *Config) {
	viper.SetConfigName(".default")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return config
}
