package main

import (
	"os"
	"signal/helpers"
	"signal/jobs"
	"signal/models"
	"time"

	"github.com/docker/docker/client"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var jobScheduler, _ = gocron.NewScheduler()
var state = false
var taskPayload = models.TaskPayload{
	BearerToken: "",
	BackendUrl:  "",
	UserId:      "",
}
var jobId = uuid.Nil
var dockerClient *client.Client
var containersState = make(map[string]*time.Time)

type AgentStatePayload struct {
	State bool `json:"state"`
}

type AgentAuthDataPayload struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

func main() {
	logger.SetOutput(os.Stdout)

	cfs := helpers.GetEnvVariables()
	dockerClient, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	taskPayload.BackendUrl = cfs.BackendApiAddress
	job, err := jobScheduler.NewJob(
		gocron.DurationJob(time.Second*10),
		gocron.NewTask(jobs.ScanForErrors, dockerClient, logger, taskPayload, containersState),
	)
	if err != nil {
		logger.Fatalf("Failed to create job: %v", err)
	}
	jobId = job.ID()
	router := echo.New()
	router.HideBanner = true
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	router.POST("/api/control/state", ControlPower)
	router.GET("/api/control/state", GetState)
	router.POST("/api/control/auth_data", ControlAuthData)
	logger.Fatal(router.Start(":37002"))

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
	logger.Infof("State: %v", statePayload.State)
	if state == statePayload.State {
		c.JSON(200, "Success")
		return nil
	}
	if statePayload.State {
		state = statePayload.State
		logger.Infof("Starting collector")
		jobScheduler.Start()
		logger.Infof("Collector started")
	} else {
		state = statePayload.State
		jobScheduler.StopJobs()
	}
	c.JSON(200, "Success")
	return nil
}

func ControlAuthData(c echo.Context) error {
	var agentAuthDataPayload AgentAuthDataPayload
	if err := c.Bind(&agentAuthDataPayload); err != nil {
		c.JSON(400, "Invalid value for token")
		return nil
	}
	taskPayload.BearerToken = agentAuthDataPayload.Token
	taskPayload.UserId = agentAuthDataPayload.UserId
	jobScheduler.Update(
		jobId,
		gocron.DurationJob(time.Second*10),
		gocron.NewTask(jobs.ScanForErrors, dockerClient, logger, taskPayload, containersState),
	)
	c.JSON(200, "Success")
	return nil
}
