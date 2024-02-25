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
var jscheduler, _ = gocron.NewScheduler()
var state = false
var token = ""
var userId = ""
var jId = uuid.Nil
var cli *client.Client

type AgentStatePayload struct {
	State bool `json:"state"`
}

type AgentAuthDataPayload struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

func main() {
	var bearerToken = "Bearer " + token
	logger.SetOutput(os.Stdout)

	cfs := helpers.GetEnvVariables()
	cli, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	taskPayload := models.TaskPayload{
		BearerToken: bearerToken,
		BackendUrl:  cfs.BackendApiAddress,
		UserId:      userId,
	}
	j, err := jscheduler.NewJob(
		gocron.DurationJob(time.Second*15),
		gocron.NewTask(jobs.ScanForErrors, cli, logger, taskPayload),
	)
	if err != nil {
		logger.Fatalf("Failed to create job: %v", err)
	}
	jId = j.ID()
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
		jscheduler.Start()
		logger.Infof("Collector started")
	} else {
		state = statePayload.State
		jscheduler.StopJobs()
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
	token = agentAuthDataPayload.Token
	userId = agentAuthDataPayload.UserId
	jscheduler.Update(
		jId,
		gocron.DurationJob(time.Second*15),
		gocron.NewTask(jobs.ScanForErrors, cli, logger, token),
	)
	c.JSON(200, "Success")
	return nil
}
