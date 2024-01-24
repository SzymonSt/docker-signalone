package main

import (
	"net"
	"os"
	"signal/helpers"
	"signal/jobs"
	"time"

	"github.com/docker/docker/client"
	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var jscheduler, _ = gocron.NewScheduler()
var state AgentState = AgentState{State: false}
var token AgentToken = AgentToken{Token: ""}

type AgentState struct {
	State bool `json:"state"`
}

type AgentToken struct {
	Token string `json:"token"`
}

func main() {
	var bearerToken = "Bearer " + token.Token
	var socketPath = "/run/guest-services/backend.sock"
	logger.SetOutput(os.Stdout)

	logger.Infof("Starting collector")
	_ = helpers.GetEnvVariables()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatalf("Failed to create docker client: %v", err)
	}
	_, err = jscheduler.NewJob(
		gocron.DurationJob(time.Second*15),
		gocron.NewTask(jobs.ScanForErrors, cli, logger, &bearerToken),
	)
	if err != nil {
		logger.Fatalf("Failed to create job: %v", err)
	}
	router := echo.New()
	router.HideBanner = true
	startUrl := ""

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		logger.Fatalf("Failed to create socket: %v", err)
	}
	router.Listener = l
	router.POST("/api/control/power", ControlPower)
	router.POST("/api/control/token", ControlToken)
	logger.Fatal(router.Start(startUrl))

}

func ControlPower(c echo.Context) error {
	if err := c.Bind(&state); err != nil {
		return c.JSON(400, "Invalid value for on")
	}
	if state.State {
		jscheduler.Start()
	} else {
		jscheduler.StopJobs()
	}
	return nil
}

func ControlToken(c echo.Context) error {
	if err := c.Bind(&token); err != nil {
		return c.JSON(400, "Invalid value for token")
	}
	return nil
}
