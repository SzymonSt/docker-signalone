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

func ScanForErrors(dockerClient *client.Client, logger *logrus.Logger, taskPayload models.TaskPayload) {
	containers, err := helpers.ListContainers(dockerClient)
	if err != nil {
		logger.Errorf("Failed to list containers: %v", err)
		return
	}
	wg := sync.WaitGroup{}
	for _, c := range containers {
		wg.Add(1)
		go func(dockerClient *client.Client,
			c types.Container, l *logrus.Logger,
			wg *sync.WaitGroup, taskPayload models.TaskPayload) {
			isErrorState := false
			execTimeOffsetInSeconds := -5
			timeTail := time.Now().Add(time.Duration(-15 + execTimeOffsetInSeconds)).Format(time.RFC3339)
			defer wg.Done()
			l.Infof("Authorization: Bearer %s \n", taskPayload.BearerToken)
			container, err := dockerClient.ContainerInspect(context.Background(), c.ID)
			if err != nil {
				l.Errorf("Failed to inspect container %s: %v", c.ID, err)
				return
			}
			logs, err := helpers.CollectLogsForAnalysis(c.ID, dockerClient, timeTail)
			if err != nil {
				l.Errorf("Failed to collect logs for container %s: %v", c.ID, err)
			}
			isErrorState = isContainerInErrorState(container.State)
			if isErrorState && logs != "" {
				err := helpers.CallLogAnalysis(logs, c.Names[0], taskPayload)
				if err != nil {
					l.Errorf("Failed to call log analysis for container %s: %v", c.Names[0], err)
				}
				return
			}
			isErrorState = areLogsIndicatingErrorOrWarning(logs)
			if isErrorState {
				err := helpers.CallLogAnalysis(logs, c.Names[0], taskPayload)
				if err != nil {
					l.Errorf("Failed to call log analysis for container %s: %v", c.Names[0], err)
				}
			}
		}(dockerClient, c, logger, &wg, taskPayload)
	}

	wg.Wait()
}

func isContainerInErrorState(state *types.ContainerState) bool {
	return (state.Error != "" ||
		(!state.Running && state.ExitCode != 0))
}

func areLogsIndicatingErrorOrWarning(logs string) bool {
	regexWarningError := `(?i)(abort|blocked|corrupt|crash|critical|deadlock|denied|
		err|error|exception|fatal|forbidden|freeze|hang|illegal|invalid|issue|missing|
		panic|rejected|refused|stacktrace|timeout|traceback|unauthorized|uncaught|unexpected|unhandled|
		unimplemented|unsupported|warn|warning)`
	matched, _ := regexp.MatchString(regexWarningError, strings.ToLower(logs))
	return matched
}
