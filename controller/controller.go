package controller

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"time"

	"../simulation"
	"../view"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Controller handles the network requests and manages simualtions.
type Controller struct {
	// server is a http server that accepts network requests
	server *http.Server
	// simulations stores a slice of simulations that are created
	// on request by a user
	simulations []simulation.Simulation
	// unityViewer is a server that handles the connection to unity
	// and is responsable for visualising the simulation
	unityViewer view.UnityServer

	// Logger is used to output information about the servers condition
	Logger *log.Entry
}

// NewController ...
func NewController(apiPort, unityPort string) (Controller, error) {
	var c Controller
	c.setup(apiPort)
	c.unityViewer = view.NewUnityServer(unityPort)
	err := c.unityViewer.StartServer()
	if err != nil {
		c.Logger.Errorf("Error: starting unityServer - %v", err)
		return c, err
	}
	return c, nil
}

// Setup intilises the logger and server
func (c *Controller) setup(port string) {
	c.Logger = log.WithFields(log.Fields{"package": "controller"})

	c.Logger.Debug("Setting up the server")
	c.Logger.Debug("Port: " + port)
	c.Logger.Debug("ReadTimeout: 10s")
	c.Logger.Debug("WriteTimeout: 10s")

	// Setup the router
	// The router maps the subdomains to functions
	router := mux.NewRouter()

	// Assign endpoints
	router.HandleFunc("/test", c.test).Methods("GET")
	router.HandleFunc("/shutdown", c.Shutdown).Methods("GET")

	// Setup the http server
	c.server = &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

// Listen tells the http server to start listening for requests
func (c *Controller) Listen() {
	if c.server == nil {
		c.Logger.Fatal("Server has not been setup")
	}

	c.Logger.Info("Server is Listening ...")
	c.server.ListenAndServe()
}

// test responds to the request with a simple message
func (c *Controller) test(w http.ResponseWriter, r *http.Request) {
	c.Logger.Debug("Received: " + html.EscapeString(r.URL.Path))

	fmt.Fprintf(w, "Everything is working, %q", html.EscapeString(r.URL.Path))
}

// Shutdown gracefully stops the http and unity server.
func (c *Controller) Shutdown(w http.ResponseWriter, r *http.Request) {
	c.Logger.Debug("Received: " + html.EscapeString(r.URL.Path))

	// Check server has been setup
	if c.server == nil {
		c.Logger.Fatal("Server has not been setup")
		return
	}

	// Close the unity server connection
	c.unityViewer.StopServer()

	// Stop the server from listening
	c.server.Shutdown(context.Background())

}
