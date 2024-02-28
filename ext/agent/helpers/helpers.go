package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"signal/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

type ConfigServer struct {
	BackendApiKey     string `mapstructure:"BACKEND_API_KEY"`
	BackendApiAddress string `mapstructure:"BACKEND_API_ADDRESS"`
}

func ListContainers(cli *client.Client) ([]types.Container, error) {
	filteredContainers := make([]types.Container, 0)
	containers, err := cli.ContainerList(context.Background(),
		types.ContainerListOptions{
			All: true,
		},
	)
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		if _, exists := c.Labels["com.docker.desktop.extension"]; !exists {
			filteredContainers = append(filteredContainers, c)
		}
	}

	return filteredContainers, nil
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

func CallLogAnalysis(logs string, containerName string, taskPayload models.TaskPayload) (err error) {
	data := map[string]string{
		"logs":          logs,
		"containerName": containerName,
		"userId":        taskPayload.UserId,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	issueAnalysisReq, err := http.NewRequest("PUT", taskPayload.BackendUrl+"/api/agent/issues/analysis", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	issueAnalysisReq.Header.Set("Content-Type", "application/json")
	issueAnalysisReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", taskPayload.BearerToken))
	client := &http.Client{}
	resp, err := client.Do(issueAnalysisReq)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to call log analysis: %v", resp.Status)
	}
	return
}
