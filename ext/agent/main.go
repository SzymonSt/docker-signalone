package main

import (
	"os"
	"signal/helpers"
	"signal/jobs"
	"time"

	"github.com/docker/docker/client"
	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	logger.SetOutput(os.Stdout)

	logger.Infof("Starting collector")
	_ = helpers.GetEnvVariables()
	jscheduler, err := gocron.NewScheduler()
	if err != nil {
		logger.Fatalf("Failed to create scheduler: %v", err)
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatalf("Failed to create docker client: %v", err)
	}
	_, err = jscheduler.NewJob(
		gocron.DurationJob(time.Second*15),
		gocron.NewTask(jobs.ScanForErrors, cli, logger),
	)
	if err != nil {
		logger.Fatalf("Failed to create job: %v", err)
	}
	jscheduler.Start()

}
