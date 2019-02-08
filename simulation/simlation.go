package simulation

import (
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
