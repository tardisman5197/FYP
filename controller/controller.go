package controller

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Controller handles the network requests and manages simualtions.
type Controller struct {
	// server is a http server that accepts network requests and
	server *http.Server

	// Logger is used to output information about the servers condition
	Logger *log.Entry
}

// NewController ...
func NewController(addr string) Controller {
	var c Controller
	c.setup(addr)
	return c
}

// Setup intilises the logger and server
func (c *Controller) setup(addr string) {
	c.Logger = log.WithFields(log.Fields{"package": "controller"})

	c.Logger.Debug("Setting up the server")
	c.Logger.Debug("Address: " + addr)
	c.Logger.Debug("ReadTimeout: 10s")
	c.Logger.Debug("WriteTimeout: 10s")

	// Setup the router
	// The router maps the subdomains to functions
	router := mux.NewRouter()

	router.HandleFunc("/test", c.test).Methods("GET")
	router.HandleFunc("/shutdown", c.Shutdown).Methods("GET")

	// Setup the http server
	c.server = &http.Server{
		Addr:           addr,
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

// Shutdown gracefully stops the http server.
func (c *Controller) Shutdown(w http.ResponseWriter, r *http.Request) {
	c.Logger.Debug("Received: " + html.EscapeString(r.URL.Path))

	if c.server == nil {
		c.Logger.Fatal("Server has not been setup")
	}

	c.server.Shutdown(context.Background())
}
