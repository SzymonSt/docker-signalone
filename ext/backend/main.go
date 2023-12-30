package main

import (
	"context"
	"os"

	"github.com/docker/docker/client"
	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	logger.SetOutput(os.Stdout)

	logger.Infof("Starting collector")
	jscheduler, err := gocron.NewScheduler()
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatalf("Failed to create docker client: %v", err)
	}

}
