package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"signal/models"
	"strconv"
	"strings"
	"time"

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

	for _, container := range containers {
		if _, exists := container.Labels["com.docker.desktop.extension"]; !exists {
			filteredContainers = append(filteredContainers, container)
		}
	}

	return filteredContainers, nil
}

func CollectLogsForAnalysis(containerID string, dockerClient *client.Client) ([]models.LogEntry, error) {
	const MaxLogTail = 8
	const LogStringBuffer = 8
	const LogTimestampParsingTemplate = "2006-01-02T15:04:05.000000000Z"

	var logEntries []models.LogEntry
	logs, err := dockerClient.ContainerLogs(context.Background(),
		containerID,
		types.ContainerLogsOptions{
			Timestamps: true,
			Tail:       strconv.Itoa(MaxLogTail),
			ShowStdout: true,
			ShowStderr: true,
		})
	if err != nil {
		return nil, err
	}
	defer logs.Close()

	logBytes, err := ioutil.ReadAll(logs)
	if err != nil {
		return nil, err
	}

	logSlice := bytes.Split(logBytes, []byte("\n"))
	for _, log := range logSlice {
		if len(log) < MaxLogTail {
			continue
		}
		logStringSlice := string(log[LogStringBuffer:])
		logTimestamp, err := time.Parse(LogTimestampParsingTemplate, strings.Fields(logStringSlice)[0])
		if err != nil {
			return nil, err
		}
		logString := strings.Fields(logStringSlice)[1:]
		logStringSlice = strings.Join(logString, " ")
		entry := models.LogEntry{
			Timestamp: logTimestamp,
			Log:       logStringSlice,
		}
		logEntries = append(logEntries, entry)
	}
	return logEntries, nil
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

func CallLogAnalysis(logs string, containerName string, containerId string, severity string, taskPayload models.TaskPayload) (err error) {
	data := map[string]string{
		"logs":          logs,
		"severity":      severity,
		"containerName": containerName,
		"containerId":   containerId,
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

func DeleteContainerIssues(containerId string, taskPayload models.TaskPayload) (err error) {
	issueDeletionReq, err := http.NewRequest("DELETE", taskPayload.BackendUrl+"/api/agent/issues/"+containerId, nil)
	if err != nil {
		return
	}
	issueDeletionReq.Header.Set("Content-Type", "application/json")
	issueDeletionReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", taskPayload.BearerToken))
	client := &http.Client{}
	resp, err := client.Do(issueDeletionReq)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to delete container issues: %v", resp.Status)
	}
	return
}
