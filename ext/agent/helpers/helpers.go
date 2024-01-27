package helpers

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

type ConfigServer struct {
	BackendApiKey     string `mapstructure:"BACKEND_API_KEY"`
	BackendApiAddress string `mapstructure:"BACKEND_API_ADDRESS"`
}

func ListContainers(cli *client.Client) ([]types.Container, error) {

	containers, err := cli.ContainerList(context.Background(),
		types.ContainerListOptions{
			All: true,
		},
	)
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func CollectLogsForAnalysis(containerID string, cli *client.Client, logTimeTail string) (string, error) {
	logs, err := cli.ContainerLogs(context.Background(),
		containerID,
		types.ContainerLogsOptions{
			Since:      logTimeTail,
			ShowStdout: true,
			ShowStderr: true,
		})
	if err != nil {
		return "", err
	}
	defer logs.Close()
	logBytes, err := ioutil.ReadAll(logs)
	if err != nil {
		return "", err
	}
	return string(logBytes), nil
}

func GetEnvVariables() (cfs ConfigServer) {
	viper.SetConfigName(".default")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	err = viper.Unmarshal(&cfs)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	return
}
