package simulation

// Agent is an interface that models any actors that may apear in
// the simulation.
type Agent interface {
	// Act is the method that simulates a tick for that agent.
	// If true is returned the agent has reached its final destination.
	Act(agents []Agent, env Environment) (Agent, bool)
	// GetPosition retrives the agent's current position.
	GetPosition() Vector
	// GetID retrives the agent's ID.
	GetID() int
	// GetCurrentWaypoint retrives the agent's current target waypoint.
	GetCurrentWaypoint() Vector
	// GetSpeed returns the current speed of the agent.
	GetSpeed() float64
	// GetRoute retuns the current route of the agent.
	GetRoute() []Vector
	// GetType returns the type of agent.
	GetType() string
	// GetInfo retuns a json string containing information about
	// the agent.
	GetInfo() string
	// GetFrequency returns the spawn frequency of the agent.
	GetFrequency() int
	// SetID is used to change the agent's id value.
	SetID(id int) Agent
	// SetFrequency changes the agent's frequency.
	SetFrequency(freq int) Agent
}
