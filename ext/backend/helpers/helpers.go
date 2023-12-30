package helpers

import (
	"context"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

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
			Since: logTimeTail,
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
