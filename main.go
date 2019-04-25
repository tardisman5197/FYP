package main

import (
	"os"

	"./controller"
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

	demoServer()

	// Once the http server is no longer listening the server stops
	logger.Warn("Server Stopping")
}

// demoServer tests the API
func demoServer() {
	// Create a controller and start listening
	c, err := controller.NewController(serverAddr, unityAddr)
	if err != nil {
		panic(err)
	}

	c.Listen()
}
