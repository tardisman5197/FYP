package simulation

import (
	"math"

	log "github.com/sirupsen/logrus"
)

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
}

// Vehicle implements the agent interface.
// The agent represents a simple vehicle.
type Vehicle struct {
	id int
	// position stores the current postion of the vehicle
	position Vector
	// speed stores the current speed of the vehicle
	speed float64
	// maxSpeed stores the maximum speed the vehicle can achieve.
	maxSpeed float64
	// route stores a list of waypoints the vehicle must visit to
	// reach its final destination.
	route []Vector
	// currentWaypoint stores the position of the vehicles current destination.
	currentWaypoint Vector
	// acceleration stores the rate at which the
	// vehicle can increase its speed.
	acceleration float64
	// decceleration stroes the rate at which the
	// vehicle can decrease its speed.
	decceleration float64

	// Logger is used to give a context based log to the stdout
	Logger *log.Entry
}

// NewVehicle creates a new vehicle and intilises its values
// using the paramaters provided.
func NewVehicle(
	id int,
	startPostiion Vector,
	startSpeed float64,
	maxSpeed float64,
	acceleration float64,
	route []Vector) Vehicle {

	// Init values
	v := Vehicle{}

	v.id = id
	// Setup the logger
	v.Logger = log.WithFields(log.Fields{
		"package": "simulation",
		"section": "vehicle",
		"id":      v.id})

	v.position = startPostiion
	v.speed = startSpeed
	v.maxSpeed = maxSpeed
	v.acceleration = acceleration
	v.route = route
	// Get the first waypoint
	v.getNextWaypoint()

	return v
}

// Act simulates a vehicles behaviour for one tick in the simulation.
// If true is returned the vehicle has reached its final destination.
func (v Vehicle) Act(agents []Agent, env Environment) (Agent, bool) {
	// Check if the vehicle has reached its current waypoint
	if v.position.InRange(v.currentWaypoint, margin) {
		// Check if the vehicle has reached its final destination
		// and update the vehicle's next destination.
		if !v.getNextWaypoint() {
			v.Logger.Infof("Reached Destination: %v", v.position)
			return v, true
		}
	}

	// Update speed based on surroundings
	v.updateSpeed()

	// Update the position of the vehicle
	v.updatePosition()
	return v, false
}

// GetPosition retrives the vehicles current position.
func (v Vehicle) GetPosition() Vector {
	return v.position
}

// GetID retrives the vehicles id.
func (v Vehicle) GetID() int {
	return v.id
}

// updateVelocity calculates the vehicles next velocity based upon
// the vehicle's surroundings.
func (v *Vehicle) updateSpeed() {
	// Acceleration of free vehicles:
	// 	Each vehicle of speed v < vmax with gap ≥ v+1 accelerates to v+1.
	// 	if v < vmax & gap ≥ v + a then v = v + a
	// if v.speed < v.maxSpeed {

	// }

	// Slowing down due to other cars:
	//	Each vehicle (speed v) with gap ≤ v−1 reduces its speed to gap: v → gap.
	//	if gap ≤ v-d then v = gap

	// Randomization:
	//	Each vehicle reduces its speed by one with probability
	//	1/2: v → max[ v − 1, 0 ]

}

// updatePosition uses the vehicle's current speed to calculate
// the vehicles new velocity and position.
func (v *Vehicle) updatePosition() {
	// Convert the vehicle's current speed into its x and y velocitys

	// Go East
	if v.currentWaypoint.x > v.position.x {
		// Go North
		if v.currentWaypoint.y < v.position.y {
			// NE
			// a = Atan(py-ty/tx-px)
			// y' = -s(Sin(a))
			// x' = s(Cos(a))
			dy := v.position.y - v.currentWaypoint.y
			dx := v.currentWaypoint.x - v.position.x
			angle := math.Atan(dy / dx)

			v.position.y -= v.speed * math.Sin(angle)
			v.position.x += v.speed * math.Cos(angle)
		} else {

			// SE
			// a = Atan(ty-py/tx-px)
			// y' = s(Sin(a))
			// x' = s(Cos(a))
			dy := v.currentWaypoint.y - v.position.y
			dx := v.currentWaypoint.x - v.position.x
			angle := math.Atan(dy / dx)

			v.position.y += v.speed * math.Sin(angle)
			v.position.x += v.speed * math.Cos(angle)
		}
	} else {

		// Go North
		if v.currentWaypoint.y < v.position.y {
			// NW
			// a = Atan(py-ty/px-tx)
			// y' = -s(Sin(a))
			// x' = -s(Cos(a))
			dy := v.position.y - v.currentWaypoint.y
			dx := v.position.x - v.currentWaypoint.x
			angle := math.Atan(dy / dx)

			v.position.y -= v.speed * math.Sin(angle)
			v.position.x -= v.speed * math.Cos(angle)
		} else {
			// SW
			// a = Atan(ty-py/px-tx)
			// y' = s(Sin(a))
			// x' = -s(Cos(a))
			dy := v.currentWaypoint.y - v.position.y
			dx := v.position.x - v.currentWaypoint.x
			angle := math.Atan(dy / dx)

			v.position.y += v.speed * math.Sin(angle)
			v.position.x -= v.speed * math.Cos(angle)
		}
	}
}

// getNextWaypoint gets the vehicles next waypoint from the route.
// If the vehicle does not have another waypoint to visit false is returned.
func (v *Vehicle) getNextWaypoint() bool {
	if len(v.route) == 0 {
		// No more waypoints for the vehicle.
		// Vehicle must be at its final destination.
		return false
	}

	// Update the currentWaypoint to the next waypoint on route
	v.currentWaypoint = v.route[0]
	// Remove the currentWaypoint from the route.
	v.route = v.route[1:]
	v.Logger.Debugf("CW: %v, Route: %v", v.currentWaypoint, v.route)

	return true
}
