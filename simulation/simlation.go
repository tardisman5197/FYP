package simulation

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// Simulation stores all the details of a traffic simulation.
type Simulation struct {
	// shouldStop is true if the simualtion should
	// stop and no longer run the simulation
	shouldStop bool
	// agents is a list of all the agents in the
	// simulation
	agents []Agent
	// environment stores the information about
	// the road network
	environment Environment
	// currentTick is the time that the simulation
	// is currently at
	currentTick int
	// agentsToSpawn is a list of Agents (key) that need to
	// be spawned thorughout the simualtion.
	agentsToSpawn []Agent
	// currentAgentID stores the current ID for an agent.
	// This is used to assing new agents IDs.
	currentAgentID int

	// Logger is used to print messages to the stdout
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
		// s.getImage()
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

	// Spawn agents that have a frequency
	for _, agent := range s.agentsToSpawn {
		if s.currentTick%agent.GetFrequency() == 0 {
			s.Logger.Debugf("Spawning Agent, f: %v, tick: %v", agent.GetFrequency(), s.currentTick)

			// Add the new vehicle, however change the frequency
			// so the new agent doesn't get added to the agentsToSpawn.
			s.AddAgent(agent.SetFrequency(0))
		}
	}

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

	// if the frequency is more than 0 the agent
	// needs to be spawned more than once
	if newAgent.GetFrequency() > 0 {
		s.agentsToSpawn = append(s.agentsToSpawn, newAgent)
	}

	// Give the agent a unique id
	if newAgent.GetID() < 0 {
		s.currentAgentID++
		s.Logger.Debugf("Assinging ID: %v", s.currentAgentID)
		newAgent = newAgent.SetID(s.currentAgentID)
	}

	s.Logger.Infof("Adding an Agent: %v", newAgent.GetID())
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
	type lightInfo struct {
		Stop     bool      `json:"stop"`
		Position []float64 `json:"position"`
		ID       int       `json:"id"`
	}
	type envInfo struct {
		Waypoints [][]float64 `json:"waypoints"`
		Lights    []lightInfo `json:"lights"`
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

	// Convert lights to []lightInfo
	lights := s.environment.GetLights()
	var lightsInfo []lightInfo
	// Loop through the lights
	for _, light := range lights {
		var cli lightInfo

		cli.ID = light.GetID()
		cli.Stop = light.GetStop()

		// Convert the position from vector to []float64
		lightPos := light.GetPosition()
		cli.Position = lightPos.ConvertToSlice()

		lightsInfo = append(lightsInfo, cli)
	}

	env.Lights = lightsInfo

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
// It also retuns the current waypoint for each agent.
func (s *Simulation) GetAgentPositions() (positions [][]float64, goals [][]float64) {

	// Loop through all the agents and get their positions in a []float64
	// format and append them to the results
	for _, agent := range s.agents {
		p := agent.GetPosition()
		positions = append(positions, p.ConvertToSlice())
		wp := agent.GetCurrentWaypoint()
		goals = append(goals, wp.ConvertToSlice())
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

// GetAgents retuns the list of agents in the simulation.
func (s *Simulation) GetAgents() []Agent {
	return s.agents
}

// AddLight adds a traffic light at the given position with
// the stop state specified.
func (s *Simulation) AddLight(pos Vector, stop bool) {
	s.environment.AddLight(pos, stop)
}

// UpdateLight updates the state of a given light in the
// simulation.
func (s *Simulation) UpdateLight(id int, stop bool) {
	s.environment.UpdateLight(id, stop)
}

// GetLights returns the positions and current states of all the lights in
// the environment in the form of [][]flaot64 and []bool.
func (s *Simulation) GetLights() (positions [][]float64, states []bool) {
	for _, light := range s.environment.GetLights() {
		currentPos := light.GetPosition()
		positions = append(positions, currentPos.ConvertToSlice())
		states = append(states, light.GetStop())
	}
	return
}
