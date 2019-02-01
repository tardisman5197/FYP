package main

import (
	"os"

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

	// // Test Sim
	// sim := simulation.NewSimulation()
	// sim.RunSteps(15)

	// // Create a controler and start listening
	// c := controller.NewController(serverAddr)
	// c.Listen()

	// Once the http server is no longer listening the server stops
	logger.Warn("Server Stopping")
}
