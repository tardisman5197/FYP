package controller

import (
	"net/http"
	"sync"
	"time"

	"../view"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Controller handles the network requests and manages simualtions.
type Controller struct {
	// server is a http server that accepts network requests
	server *http.Server
	// simulations stores a slice of simulations that are created
	// on request by a user
	simulations sync.Map
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

	// simulation endpoints
	router.HandleFunc("/simulation/new", c.newSimulation).Methods("POST")
	router.HandleFunc("/simulation/remove/{id}", c.removeSimulation).Methods("GET")
	router.HandleFunc("/simulation/run/{id}", c.runSimulation).Methods("POST")
	router.HandleFunc("/simulation/stop/{id}", c.stopSimulation).Methods("GET")
	router.HandleFunc("/simulation/add/{id}", c.addAgent).Methods("POST")
	router.HandleFunc("/simulation/info/agent/{id}/{agentId}", c.getAgentInfo).Methods("GET")
	router.HandleFunc("/simulation/info/{id}", c.getInfo).Methods("GET")
	router.HandleFunc("/simulation/view/{id}", c.getImage).Methods("POST")

	// server endpoints
	router.HandleFunc("/shutdown", c.Shutdown).Methods("GET")

	// Setup the http server
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	c.server = &http.Server{
		Addr:           port,
		Handler:        handlers.CORS(headersOk, originsOk, methodsOk)(router),
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
