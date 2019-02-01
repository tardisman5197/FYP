package simulation

import log "github.com/sirupsen/logrus"

// Simulation stores all the details of a traffic simulation.
type Simulation struct {
	shouldStop  bool
	agents      []Agent
	environment Environment

	Logger *log.Entry
}

// NewSimulation creates a new Simulation struct
// and initlises some of the values.
func NewSimulation() Simulation {
	var sim Simulation

	// Setup Logger
	sim.Logger = log.WithFields(log.Fields{
		"package": "simulation",
		"section": "sim"})

	// Init shouldStop to false.
	// If set to true the simulation will stop running.
	sim.shouldStop = false

	// Get a new Environment and set
	sim.environment = NewEnvironment()
	sim.environment.WriteShapeFile("resources/test.shp")
	sim.environment.ReadShapefile("resources/test.shp")

	// TEMP Setup test vehicle
	startLoc := Vector{x: 0, y: 0}
	sim.agents = append(sim.agents,
		NewVehicle(0, startLoc, 2, 2, 0, sim.environment.waypoints))

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

		if s.shouldStop {
			break
		}
	}
}

// runOneStep simulates a single second in the simulation.
// At each second each fo the agent's act functions are called.
func (s *Simulation) runOneStep() {
	var toRemove []int

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

// removeAgent removes the agent at a specified index from the simulation's
// list of agents.
func (s *Simulation) removeAgent(index int) {
	s.Logger.Infof("Removing Agent: %v", s.agents[index].GetID())
	s.agents = append(s.agents[:index], s.agents[index+1:]...)
}
