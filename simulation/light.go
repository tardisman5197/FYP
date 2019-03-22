package simulation

import log "github.com/sirupsen/logrus"

// Light represents a traffic light on a road
type Light struct {
	// id is a unique integer used to identify the
	// traffic light.
	id int
	// stop is true if the traffic light is
	// red and false if green.
	stop bool
	// position is a vecotor storing the location of the traffic light.
	// This location should be the same location as a target waypoint
	// for the vehicle you want to stop.
	position Vector

	// Logger is used to give a context based log to the stdout
	Logger *log.Entry
}

// NewLight returns a Light with the specified paramaters.
func NewLight(id int, pos Vector, stop bool) Light {
	var l Light
	l.id = id
	l.position = pos
	l.stop = stop

	// Setup the logger
	l.Logger = log.WithFields(log.Fields{
		"package": "simulation",
		"section": "Light"})

	return l
}

// GetID returns the id of the traffic light.
func (l *Light) GetID() int {
	return l.id
}

// GetStop retuns the current value of the stop bool.
func (l *Light) GetStop() bool {
	return l.stop
}

// SetStop updates the stop bool to the value specified.
func (l *Light) SetStop(stop bool) {
	l.Logger.Debugf("%v Light set to: %v", l.id, l.stop)
	l.stop = stop
}

// GetPosition retuns the position of the traffic light.
func (l *Light) GetPosition() Vector {
	return l.position
}
