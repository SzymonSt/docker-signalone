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

func ScanForErrors(cli *client.Client, logger *logrus.Logger, taskPayload models.TaskPayload) {
	containers, err := helpers.ListContainers(cli)
	if err != nil {
		logger.Errorf("Failed to list containers: %v", err)
		return
	}
	wg := sync.WaitGroup{}
	for _, c := range containers {
		wg.Add(1)
		go func(cli *client.Client,
			c types.Container, l *logrus.Logger,
			wg *sync.WaitGroup, taskPayload models.TaskPayload) {
			isErrorState := false
			execTimeOffsetInSeconds := -0.5
			timeTail := time.Now().Add(time.Duration(-15 + execTimeOffsetInSeconds)).Format(time.RFC3339)
			defer wg.Done()
			l.Infof("Authorization: Bearer %s \n", taskPayload.BearerToken)
			container, err := cli.ContainerInspect(context.Background(), c.ID)
			if err != nil {
				l.Errorf("Failed to inspect container %s: %v", c.ID, err)
				return
			}
			logs, err := helpers.CollectLogsForAnalysis(c.ID, cli, timeTail)
			if err != nil {
				l.Errorf("Failed to collect logs for container %s: %v", c.ID, err)
			}
			isErrorState = isContainerInErrorState(container.State)
			if isErrorState {
				helpers.CallLogAnalysis(logs, c.Names[0], taskPayload)
				return
			}
			isErrorState = areLogsIndicatingErrorOrWarning(logs)
			if isErrorState {
				helpers.CallLogAnalysis(logs, c.Names[0], taskPayload)
			}
		}(cli, c, logger, &wg, taskPayload)
	}

	wg.Wait()
}

func isContainerInErrorState(state *types.ContainerState) bool {
	return (state.Error != "" ||
		(!state.Running && state.ExitCode != 0))
}

func areLogsIndicatingErrorOrWarning(logs string) bool {
	regexWarningError := `(?i)(error|warning|exception|err|warn|critical|fatal|stacktrace|
	traceback|issue|crash|hang|freeze|
	timeout|deadlock|corrupt|invalid|illegal|unhandled|uncaught|
	unexpected|unimplemented|unsupported|missing|invalid|illegal|
	unauthorized|denied|forbidden|blocked|rejected|panic|abort)`
	matched, _ := regexp.MatchString(regexWarningError, strings.ToLower(logs))
	return matched
}
