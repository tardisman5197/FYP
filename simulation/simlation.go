package simulation

import (
	"encoding/json"

	"../view"
	log "github.com/sirupsen/logrus"
)

// Simulation stores all the details of a traffic simulation.
type Simulation struct {
	shouldStop  bool
	agents      []Agent
	environment Environment
	currentTick int

	Logger *log.Entry
}

// NewSimulation creates a new Simulation struct
// and initlises some of the values.
func NewSimulation(env Environment) Simulation {
	var sim Simulation

	// Setup Logger
	sim.Logger = log.WithFields(log.Fields{
		"package": "simulation",
		"section": "sim"})

	// Init shouldStop to false.
	// If set to true the simulation will stop running.
	sim.shouldStop = false

	sim.environment = env

	return sim
}

// Run loops until the simulation's shouldStop variable is set to true.
func (s *Simulation) Run() {
	for {
		s.runOneStep()

		if s.shouldStop {
			break
		}
	}
}

// RunSteps runs the simulation a specified number or until the simulation's
// shouldStop variable is set to true.
func (s *Simulation) RunSteps(noOfSteps int) {
	for i := 0; i < noOfSteps; i++ {
		s.runOneStep()

		// if i%5 == 0 {
		s.getImage()
		// }
		if s.shouldStop {
			break
		}
	}
}

// runOneStep simulates a single second in the simulation.
// At each second each fo the agent's act functions are called.
func (s *Simulation) runOneStep() {
	var toRemove []int

	s.currentTick++
	s.Logger.Infof("Current Tick: %v", s.currentTick)

	// Loop over each agent and execute act function
	for i := 0; i < len(s.agents); i++ {
		removeAgent := false

		s.agents[i], removeAgent = s.agents[i].Act(s.agents, s.environment)

		if removeAgent {
			toRemove = append(toRemove, i)
			continue
		}
	}

	// Remove agents that have reached their destination
	for i := 0; i < len(toRemove); i++ {
		s.removeAgent(toRemove[i])
	}
}

// Stop sets the simulation's shouldStop variable to true.
// If a simulation is currntly running this function should notify the
// simulation to stop at the end of the current tick.
func (s *Simulation) Stop() {
	s.shouldStop = true
}

// AddAgent adds an agent to the simulation.
func (s *Simulation) AddAgent(newAgent Agent) {
	s.Logger.Info("Adding an Agent")
	s.agents = append(s.agents, newAgent)
}

// removeAgent removes the agent at a specified index from the simulation's
// list of agents.
func (s *Simulation) removeAgent(index int) {
	s.Logger.Infof("Removing Agent: %v", s.agents[index].GetID())
	s.agents = append(s.agents[:index], s.agents[index+1:]...)
}

// GetInfo returns a json string containing the current information
// of the simulation.
func (s *Simulation) GetInfo() string {
	type envInfo struct {
		Waypoints [][]float64 `json:"waypoints"`
	}

	type agentInfo struct {
		ID              int         `json:"id"`
		Position        []float64   `json:"position"`
		Speed           float64     `json:"speed"`
		CurrentWaypoint []float64   `json:"currentWaypoint"`
		Route           [][]float64 `json:"route"`
		Type            string      `json:"type"`
	}

	type simInfo struct {
		Agents      []agentInfo `json:"agents"`
		Environment envInfo     `json:"environment"`
		Tick        int         `json:"tick"`
	}

	type response struct {
		Success bool    `json:"success"`
		Error   string  `json:"error"`
		Info    simInfo `json:"info"`
	}

	// Create simInfo to store information about the simulation
	var sim simInfo

	sim.Tick = s.currentTick

	// Set the environment info
	var env envInfo
	// Convert Waypoints from Vectors to []float64
	for _, waypoint := range s.environment.GetWaypoints() {
		env.Waypoints = append(env.Waypoints, waypoint.ConvertToSlice())
	}
	sim.Environment = env

	// Sets the agent information
	for _, agent := range s.agents {
		p := agent.GetPosition()
		cwp := agent.GetCurrentWaypoint()
		// Convert []Vector to [][]float64
		var r [][]float64
		for _, wp := range agent.GetRoute() {
			r = append(r, wp.ConvertToSlice())
		}

		currentAgent := agentInfo{
			ID:              agent.GetID(),
			Position:        p.ConvertToSlice(),
			Speed:           agent.GetSpeed(),
			CurrentWaypoint: cwp.ConvertToSlice(),
			Route:           r,
			Type:            agent.GetType()}

		sim.Agents = append(sim.Agents, currentAgent)
	}

	// Convert the infomation into a json string
	var r response
	r.Success = true
	r.Info = sim
	jsonStr, _ := json.Marshal(r)
	return string(jsonStr)
}

// GetAgent retuns the specified agent by id.
// A nil response means no agent was found.
func (s *Simulation) GetAgent(id int) Agent {
	for _, agent := range s.agents {
		if agent.GetID() == id {
			return agent
		}
	}
	return nil
}

// GetAgentPositions retuns the positions of all the agents in the simulation.
func (s *Simulation) GetAgentPositions() (positions [][]float64) {

	// Loop through all the agents and get their positions in a []float64
	// format and append them to the results
	for _, agent := range s.agents {
		p := agent.GetPosition()
		positions = append(positions, p.ConvertToSlice())
	}

	return
}

// GetWaypoints returns the simulation's waypoints defined in the
// environment in a [][]float64 format.
func (s *Simulation) GetWaypoints() (waypoints [][]float64) {
	for _, w := range s.environment.GetWaypoints() {
		waypoints = append(waypoints, w.ConvertToSlice())
	}
	return
}

// GetTick retuns the current tick the simulation is on.
func (s *Simulation) GetTick() int {
	return s.currentTick
}
func (s *Simulation) getImage() {
	var wp [][]float64
	for _, cwp := range s.environment.waypoints {
		wp = append(wp, cwp.ConvertToSlice())
	}

	var a [][]float64
	for _, ca := range s.agents {
		cp := ca.GetPosition()
		a = append(a, cp.ConvertToSlice())
	}

	view.GenImg(wp, a, s.currentTick)
}
