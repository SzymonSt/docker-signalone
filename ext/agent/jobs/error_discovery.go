package jobs

import (
	"context"
	"regexp"
	"signal/helpers"
	"signal/models"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func ScanForErrors(dockerClient *client.Client, logger *logrus.Logger, taskPayload models.TaskPayload, containersState map[string]*time.Time) {
	var currentsIDs = make([]string, 0)
	var mutex = &sync.RWMutex{}

	containers, err := helpers.ListContainers(dockerClient)
	if err != nil {
		logger.Errorf("Failed to list containers: %v", err)
		return
	}
	wg := sync.WaitGroup{}
	for _, container := range containers {
		currentsIDs = append(currentsIDs, container.ID)
		mutex.RLock()
		_, exists := containersState[container.ID]
		mutex.RUnlock()
		if !exists {
			containerCreationTime := time.Unix(container.Created, 0)
			mutex.RLock()
			containersState[container.ID] = &containerCreationTime
			mutex.RUnlock()
		}

		wg.Add(1)
		go func(dockerClient *client.Client,
			containerDefinition types.Container, logger *logrus.Logger,
			wg *sync.WaitGroup, taskPayload models.TaskPayload) {
			var timestampCheckpoint time.Time
			isErrorState := false
			logString := ""
			severity := "INFO"
			defer wg.Done()
			container, err := dockerClient.ContainerInspect(context.Background(), containerDefinition.ID)
			if err != nil {
				logger.Errorf("Failed to inspect container %s: %v", containerDefinition.ID, err)
				return
			}

			logs, err := helpers.CollectLogsForAnalysis(containerDefinition.ID, dockerClient)
			if err != nil {
				logger.Errorf("Failed to collect logs for container %s: %v", containerDefinition.ID, err)
			}

			mutex.RLock()
			for _, log := range logs {
				if log.Timestamp.Add(-1 * time.Second).After(*containersState[containerDefinition.ID]) {
					logString += (log.Log + "\n")
					timestampCheckpoint = log.Timestamp
				}
			}
			if logString != "" {
				containersState[containerDefinition.ID] = &timestampCheckpoint
			}
			mutex.RUnlock()

			isErrorState = checkContainerErrorState(container.State)
			if isErrorState && logString != "" {
				severity = "CRITICAL"
				err := helpers.CallLogAnalysis(logString, containerDefinition.Names[0], containerDefinition.ID, severity, taskPayload)
				if err != nil {
					logger.Errorf("Failed to call log analysis for container %s: %v", containerDefinition.Names[0], err)
				}
				return
			}

			isErrorState, severity = checkLogsForIssue(logString)
			if isErrorState {
				err := helpers.CallLogAnalysis(logString, containerDefinition.Names[0], containerDefinition.ID, severity, taskPayload)
				if err != nil {
					logger.Errorf("Failed to call log analysis for container %s: %v", containerDefinition.Names[0], err)
				}
				return
			}
		}(dockerClient, container, logger, &wg, taskPayload)
	}

	wg.Wait()
	mutex.RLock()
	for key, _ := range containersState {
		if verifyIfContainerDeleted(key, currentsIDs) {
			delete(containersState, key)
			helpers.DeleteContainerIssues(key, taskPayload)
		}
	}
	mutex.RUnlock()
}

func checkContainerErrorState(state *types.ContainerState) bool {
	return (state.Error != "" ||
		(!state.Running && state.ExitCode != 0))
}

func checkLogsForIssue(logs string) (matched bool, severity string) {
	regexWarning := `(?i)(deprecated|deprecating|unsupported|warn|warning)`
	matched, _ = regexp.MatchString(regexWarning, strings.ToLower(logs))
	if matched {
		severity = "WARNING"
	}

	regexError := `(?i)(abort|blocked|corrupt|crash|critical|deadlock|denied|
		err|error|exception|fatal|forbidden|freeze|hang|illegal|invalid|missing|
		panic|rejected|refused|stacktrace|timeout|traceback|unauthorized|uncaught|unexpected|unhandled|
		unimplemented)`
	matched, _ = regexp.MatchString(regexError, strings.ToLower(logs))
	if matched {
		severity = "CRITICAL"
	}

	return
}

func verifyIfContainerDeleted(key string, currentsIDs []string) bool {
	for _, containerId := range currentsIDs {
		if containerId == key {
			return false
		}
	}
	return true
}
