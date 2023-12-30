package jobs

import (
	"context"
	"regexp"
	"signal/helpers"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func ScanForErrors(cli *client.Client, logger *logrus.Logger) {
	containers, err := helpers.ListContainers(cli)
	if err != nil {
		logger.Errorf("Failed to list containers: %v", err)
		return
	}
	wg := sync.WaitGroup{}
	for _, c := range containers {
		wg.Add(1)
		go func(cli *client.Client, c types.Container, l *logrus.Logger, wg *sync.WaitGroup) {
			isErrorState := false
			timeTail := time.Now().Add(time.Second * -15).Format(time.RFC3339)
			defer wg.Done()
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
				// TODO: send logs to analysis
				return
			}
			isErrorState = areLogsIndicatingErrorOrWarning(logs)
			if isErrorState {
				// TODO: send logs to analysis
			}
		}(cli, c, logger, &wg)
	}
}

func isContainerInErrorState(state *types.ContainerState) bool {
	return (state.Error != "" ||
		(!state.Running && state.ExitCode != 0))
}

func areLogsIndicatingErrorOrWarning(logs string) bool {
	regexWarningError := `(?i)(error|warning|exception|err|warn)`
	matched, _ := regexp.MatchString(regexWarningError, strings.ToLower(logs))
	return matched
}
