package simulation

const margin = 1.0

// Agent is an interface that models any actors that may apear in
// the simulation
type Agent interface {
	// Act is the method that simulates a tick for that agent.
	// If true is returned the agent has reached its final destination.
	Act(agents []Agent, env Environment) (Agent, bool)
	// GetPosition retrives the agent's current position
	GetPosition() Vector
	// GetID retrives the agent's ID
	GetID() int
	// GetCurrentWaypoint retrives the agent's current target waypoint
	GetCurrentWaypoint() Vector
}
