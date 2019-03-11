package controller

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"../simulation"
	"github.com/gorilla/mux"
)

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
	c.Logger.Info("Server Shutdown")

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
	c.Logger.Debugf("Env Filepath: %v", simInfo.Environment)
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

// removeSimulation removes a specific simulation from the server.
func (c *Controller) removeSimulation(w http.ResponseWriter, r *http.Request) {
	// response is the information sent back to the client
	// after the request has been executed.
	type response struct {
		// Success is bool that is true if a new sim has been
		// created
		Success bool `json:"success"`
		// Error is a string that is filled if an error occurs
		// while creating a new simulation.
		Error string `json:"error"`
	}

	// Get the id from the url
	params := mux.Vars(r)
	id := params["id"]

	var resp response

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

	// remove the simulation from the server.
	c.simulations.Delete(id)

	// No Simulation found send error
	resp.Success = true

	// Encode response into json
	jsonStr, _ := json.Marshal(resp)

	// Send response
	fmt.Fprint(w, string(jsonStr))
	c.Logger.Infof("Simulation Removed: %v", id)
	return
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

// stopSimulation stops a specified simulation.
func (c *Controller) stopSimulation(w http.ResponseWriter, r *http.Request) {
	type response struct {
		// Success is true if the simulation runs for the amount of steps
		// specified.
		Success bool `json:"success"`
		// Error is a string that is set if something goes wrong.
		Error string `json:"error"`
	}

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

	// Stop the simulation
	i, _ := c.simulations.Load(id)
	sim := i.(simulation.Simulation)
	sim.Stop()
	c.simulations.Store(id, sim)

	resp.Success = true

	// Encode response into json
	jsonStr, _ := json.Marshal(resp)

	// Send response
	fmt.Fprint(w, string(jsonStr))

	c.Logger.Debugf("Sim stopped: %v", id)
}

// addAgent adds agents to a specified simulation.
func (c *Controller) addAgent(w http.ResponseWriter, r *http.Request) {
	type agentInfo struct {
		StartLocation []float64   `json:"startLocation"`
		StartSpeed    float64     `json:"startSpeed"`
		MaxSpeed      float64     `json:"maxSpeed"`
		Acceleration  float64     `json:"acceleration"`
		Deceleration  float64     `json:"deceleration"`
		Route         [][]float64 `json:"route"`
		Type          string      `json:"type"`
		Frequency     int         `json:"frequency"`
	}

	type info struct {
		Agents []agentInfo `json:"agents"`
	}

	type response struct {
		// Success is true if the simulation runs for the amount of steps
		// specified.
		Success bool `json:"success"`
		// Error is a string that is set if something goes wrong.
		Error string `json:"error"`
	}

	var resp response

	// Get the id and type from the url
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

	// Parse the agent data
	var agentsInfo info
	_ = json.NewDecoder(r.Body).Decode(&agentsInfo)

	// Get the simulation from the map
	i, _ := c.simulations.Load(id)
	sim := i.(simulation.Simulation)

	// Create and add new agent for each of the agents information given
	for _, agent := range agentsInfo.Agents {
		switch agent.Type {
		case "vehicle":
			// Convert []float64 to Vector for starting location
			startLoc := simulation.NewVector(agent.StartLocation[0], agent.StartLocation[0])

			// convert [][]float64 to list of Vectors for route
			var route []simulation.Vector
			for i := 0; i < len(agent.Route); i++ {
				newWaypoint := simulation.NewVector(agent.Route[i][0], agent.Route[i][1])
				route = append(route, newWaypoint)
			}

			// Create a vehicle
			newAgent := simulation.NewVehicle(
				-1,
				startLoc,
				agent.StartSpeed,
				agent.MaxSpeed,
				agent.Acceleration,
				agent.Deceleration,
				route,
				agent.Frequency)

			sim.AddAgent(newAgent)

		default:
			// No agent of that type
			resp.Error += "No agent of that type found - " + agent.Type + "\n"

			c.Logger.Warnf("Incorrect agent type given: %v", agent.Type)
		}

	}

	// Store the new simulation with new vehicles
	c.simulations.Store(id, sim)

	resp.Success = true

	// Encode response into json
	jsonStr, _ := json.Marshal(resp)

	// Send response
	fmt.Fprint(w, string(jsonStr))

	c.Logger.Infof("Agents been added to sim: %v", id)
}

// getInfo gets information about a specified simualtion.
func (c *Controller) getInfo(w http.ResponseWriter, r *http.Request) {
	type response struct {
		// Success is true if the simulation runs for the amount of steps
		// specified.
		Success bool `json:"success"`
		// Error is a string that is set if something goes wrong.
		Error string `json:"error"`
		// Info contains the information about the simulation
		Info string `json:"info"`
	}

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

	i, _ := c.simulations.Load(id)
	sim := i.(simulation.Simulation)
	jsonStr := sim.GetInfo()

	fmt.Fprint(w, string(jsonStr))

	c.Logger.Infof("Info returned for sim: %v", id)

}

// getAgentInfo gets information about a specified agent in a simualtion.
func (c *Controller) getAgentInfo(w http.ResponseWriter, r *http.Request) {
	type response struct {
		// Success is true if the simulation runs for the amount of steps
		// specified.
		Success bool `json:"success"`
		// Error is a string that is set if something goes wrong.
		Error string `json:"error"`
		// Info contains the information about the simulation
		Info string `json:"info"`
	}

	var resp response

	// Get the id from the url
	params := mux.Vars(r)
	id := params["id"]

	agentID, err := strconv.Atoi(params["agentId"])
	if err != nil {
		// Incorrect agent Id
		resp.Success = false
		resp.Error = "Agent Id Provided not a number - " + err.Error()

		// Encode response into json
		jsonStr, _ := json.Marshal(resp)

		// Send response
		fmt.Fprint(w, string(jsonStr))

		c.Logger.Warnf("Wrong Agent ID provided: %v", err.Error())
		return
	}

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

	// Load the simulation from the map
	i, _ := c.simulations.Load(id)
	sim := i.(simulation.Simulation)

	// Get the agent from the simulation
	agent := sim.GetAgent(agentID)
	if agent == nil {
		// No Simulation found send error
		resp.Success = false
		resp.Error = "No Agent found with the id: " + params["agentId"]

		// Encode response into json
		jsonStr, _ := json.Marshal(resp)

		// Send response
		fmt.Fprint(w, string(jsonStr))

		c.Logger.Warnf("No Agent found with id: %v", agentID)
		return
	}

	jsonStr := agent.GetInfo()

	// Send the information to the client
	fmt.Fprint(w, string(jsonStr))

	c.Logger.Infof("Info returned for sim: %v", id)

}

// getAgentInfo gets information about a specified agent in a simualtion.
func (c *Controller) getImage(w http.ResponseWriter, r *http.Request) {
	type info struct {
		// Position stores the location that the camera should
		// be set to.
		Position []float64 `json:"cameraPosition"`
		// Direction stores the location that the camera should
		// point towards.
		Direction []float64 `json:"cameraDirection"`
	}
	type response struct {
		// Success is true if the simulation runs for the amount of steps
		// specified.
		Success bool `json:"success"`
		// Error is a string that is set if something goes wrong.
		Error string `json:"error"`
		// Filepath contains the filepath on the server. This could be used
		// if the server is being ran on the same machine as the recipient.
		Filepath string `json:"filepath"`
		// Info contains the information about the simulation.
		Image string `json:"image"`
	}

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

	// Get the simulation
	i, _ := c.simulations.Load(id)
	sim := i.(simulation.Simulation)

	// parse the setup data
	var cameraInfo info
	_ = json.NewDecoder(r.Body).Decode(&cameraInfo)

	positions, goals := sim.GetAgentPositions()

	// Get the filepath for image
	resp.Filepath = c.unityViewer.GetImageFilepath(
		positions,
		sim.GetWaypoints(),
		goals,
		sim.GetTick(),
		cameraInfo.Position,
		cameraInfo.Direction)

	if sendBase64Encoding {
		// Open file and store the base64 encoding of it in the response
		f, err := os.Open(resp.Filepath)
		if err != nil {
			c.Logger.Error(err.Error())
		}
		defer f.Close()
		// Read entire JPG into byte slice.
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)
		// Encode as base64.
		resp.Image = base64.StdEncoding.EncodeToString(content)
		c.Logger.Debugf("Encoded image: %v", resp.Image)
	}

	resp.Success = true

	// Convert to json
	jsonStr, _ := json.Marshal(resp)

	// Send json to client
	fmt.Fprint(w, string(jsonStr))

	c.Logger.Infof("Image sent of sim: %v - %v", id, resp.Filepath)
}
