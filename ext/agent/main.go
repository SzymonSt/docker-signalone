package main

import (
	"net"
	"os"
	"signal/helpers"
	"signal/jobs"
	"time"

	"github.com/docker/docker/client"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var jscheduler, _ = gocron.NewScheduler()
var state = false
var token = ""
var jId = uuid.Nil
var cli *client.Client

type AgentStatePayload struct {
	State bool `json:"state"`
}

type AgentTokenPayload struct {
	Token string `json:"token"`
}

func main() {
	var bearerToken = "Bearer " + token
	var socketPath = "/run/guest-services/backend.sock"
	logger.SetOutput(os.Stdout)

	logger.Infof("Starting collector")
	_ = helpers.GetEnvVariables()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatalf("Failed to create docker client: %v", err)
	}
	j, err := jscheduler.NewJob(
		gocron.DurationJob(time.Second*15),
		gocron.NewTask(jobs.ScanForErrors, cli, logger, bearerToken),
	)
	if err != nil {
		logger.Fatalf("Failed to create job: %v", err)
	}
	jId = j.ID()
	router := echo.New()
	router.HideBanner = true
	startUrl := ""

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		logger.Fatalf("Failed to create socket: %v", err)
	}
	router.Listener = l
	router.POST("/api/control/power", ControlPower)
	router.GET("/api/control/state", GetState)
	router.POST("/api/control/token", ControlToken)
	logger.Fatal(router.Start(startUrl))

}

func GetState(c echo.Context) error {
	var statePayload AgentStatePayload
	statePayload.State = state
	c.JSON(200, statePayload)
	return nil
}

func ControlPower(c echo.Context) error {
	var statePayload AgentStatePayload
	if err := c.Bind(&statePayload); err != nil {
		c.JSON(400, "Invalid value for on")
		return nil
	}
	if state == statePayload.State {
		c.JSON(200, "Success")
		return nil
	}
	if state {
		jscheduler.Start()
	} else {
		jscheduler.StopJobs()
	}
	c.JSON(200, "Success")
	return nil
}

func ControlToken(c echo.Context) error {
	var tokenPayload AgentTokenPayload
	if err := c.Bind(&tokenPayload); err != nil {
		c.JSON(400, "Invalid value for token")
		return nil
	}
	token = tokenPayload.Token
	jscheduler.Update(
		jId,
		gocron.DurationJob(time.Second*15),
		gocron.NewTask(jobs.ScanForErrors, cli, logger, token),
	)
	c.JSON(200, "Success")
	return nil
}
