package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"sync"
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
	router.HandleFunc("/simulation/run/{id}", c.runSimulation).Methods("POST")

	// server endpoints
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

// newSimulation creates and adds a simulation to the controller.
func (c *Controller) newSimulation(w http.ResponseWriter, r *http.Request) {
	// info is a struct containing the information that is sent
	// as part of the newSimulation endpoint.
	type info struct {
		// Environment stores the filepath to the environment
		// shape file
		Environment string `json:"environment"`
	}

	// response is the information sent back to the client
	// after the request has been executed.
	type response struct {
		// Key stores the unique key given to the simulation
		// created
		Key string `json:"key"`
		// Success is bool that is true if a new sim has been
		// created
		Success bool `json:"success"`
		// Error is a string that is filled if an error occurs
		// while creating a new simulation.
		Error string `json:"error"`
	}

	// parse the setup data
	var simInfo info
	_ = json.NewDecoder(r.Body).Decode(&simInfo)

	// create response type to fill
	var resp response

	// Generate a unique key for the simulation
	var key string
	for {
		// Generate a random string key
		key = ""
		for i := 0; i < keyLength; i++ {
			// 97{a} - 122{z}
			key += string(rand.Intn(26) + 97)
		}

		// Check of key is already in use
		if _, ok := c.simulations.Load(key); !ok {
			break
		}
	}
	resp.Key = key

	// Generate the simulation environment
	env := simulation.NewEnvironment()
	// env.WriteShapeFile("resources/test.shp")
	env.ReadShapefile(simInfo.Environment)

	// Create the simulation
	sim := simulation.NewSimulation(env)

	// Add the simulation to the map
	c.simulations.Store(key, sim)

	resp.Success = true
	// Encode response into json
	jsonStr, _ := json.Marshal(resp)

	// Send response
	fmt.Fprint(w, string(jsonStr))

	c.Logger.Infof("New Simulation Created: %v", key)
}

// runSimulation runs a specified simulation.
func (c *Controller) runSimulation(w http.ResponseWriter, r *http.Request) {
	type info struct {
		// Steps is the number of steps that the simulation needs to run.
		// If the number is negative the simulation should run until
		// it is told to stop.
		Steps int `json:"steps"`
	}

	type response struct {
		// Success is true if the simulation runs for the amount of steps
		// specified.
		Success bool `json:"success"`
		// Error is a string that is set if something goes wrong.
		Error string `json:"error"`
	}

	// parse the setup data
	var cmdInfo info
	_ = json.NewDecoder(r.Body).Decode(&cmdInfo)

	// Create response for client
	var resp response

	// Get the id from the url
	params := mux.Vars(r)
	id := params["id"]
	// Check if the id exists
	if _, ok := c.simulations.Load(id); !ok {
		// No Simulation found send error
		resp.Success = false
		resp.Error = "No Simulation found with the id - " + id

		// Encode response into json
		jsonStr, _ := json.Marshal(resp)

		// Send response
		fmt.Fprint(w, string(jsonStr))

		c.Logger.Warnf("No Simulation found with id: %v", id)
		return
	}

	// Get the simulation and run for specifed number of steps
	i, _ := c.simulations.Load(id)
	sim := i.(simulation.Simulation)

	// Check if the number of steps to run is positive
	if cmdInfo.Steps > -1 {
		sim.RunSteps(cmdInfo.Steps)

		c.simulations.Store(id, sim)
	} else {
		// No Simulation found send error
		resp.Success = false
		resp.Error = "Can not have negative Step value (yet)"

		// Encode response into json
		jsonStr, _ := json.Marshal(resp)

		// Send response
		fmt.Fprint(w, string(jsonStr))

		c.Logger.Warnf("Negative Step value: %v", cmdInfo.Steps)
		return
	}

	resp.Success = true

	// Encode response into json
	jsonStr, _ := json.Marshal(resp)

	// Send response
	fmt.Fprint(w, string(jsonStr))

	c.Logger.Debugf("Sim: %v Run for %v steps", id, cmdInfo.Steps)
}
