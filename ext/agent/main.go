package main

import (
	"os"
	"signal/helpers"
	"signal/jobs"
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
	logger.SetOutput(os.Stdout)

	_ = helpers.GetEnvVariables()
	cli, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
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
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	router.POST("/api/control/state", ControlPower)
	router.GET("/api/control/state", GetState)
	router.POST("/api/control/token", ControlToken)
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
