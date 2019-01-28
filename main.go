package main

import (
	"os"

	"./controller"
	log "github.com/sirupsen/logrus"
)

// serverAddr is the port at which the server can be accessed
// e.g. "127.0.0.1:8080" (aka "localhost:8080")
const serverAddr = ":8080"

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

	// Create a controler and start listening
	c := controller.Controller{}
	c.Setup(serverAddr)
	c.Listen()

	// Once the http server is no longer listening the server stops
	logger.Warn("Server Stopping")
}
