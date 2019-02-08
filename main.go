package main

import (
	"os"

	"./simulation"
	log "github.com/sirupsen/logrus"
)

// init is called before main and intiilises the loggers paramaters.
func init() {
	// Setup the logger
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: false,
		FullTimestamp:    true,
	})
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	// Only display Debug or higher
	log.SetLevel(log.DebugLevel)
}

// main is ran when the application is executed. Main setsup a controller and
// starts it listening.
func main() {
	logger := log.WithFields(log.Fields{"package": "main"})
	logger.Info("Server Starting")

	demo()
	// // Create a controler and start listening
	// c := controller.NewController(serverAddr)
	// c.Listen()

	// Once the http server is no longer listening the server stops
	logger.Warn("Server Stopping")
}

func demo() {
	// Create the environment
	env := simulation.NewEnvironment()
	env.WriteShapeFile("resources/test.shp")
	env.ReadShapefile("resources/test.shp")

	waypoints := env.GetWaypoints()

	// Test Sim
	sim := simulation.NewSimulation(env)

	// TEMP Setup test vehicle

	startLoc := simulation.NewVector(0, 0)
	route := waypoints[1:]
	startSpeed := 0.0

	id := 0
	maxSpeed := 8.0
	acc := 3.0
	dec := 3.0
	sim.AddAgent(simulation.NewVehicle(id, startLoc, startSpeed, maxSpeed, acc, dec, route))

	sim.RunSteps(5)

	id = 1
	maxSpeed = 10.0
	acc = 3.0
	dec = 3.0
	sim.AddAgent(simulation.NewVehicle(id, startLoc, startSpeed, maxSpeed, acc, dec, route))

	sim.RunSteps(5)

	id = 2
	maxSpeed = 10.0
	acc = 4.0
	dec = 4.0
	sim.AddAgent(simulation.NewVehicle(id, startLoc, startSpeed, maxSpeed, acc, dec, route))

	sim.RunSteps(5)

	id = 3
	maxSpeed = 80.0
	acc = 5.0
	dec = 5.0
	sim.AddAgent(simulation.NewVehicle(id, startLoc, startSpeed, maxSpeed, acc, dec, route))

	sim.RunSteps(100)
}
