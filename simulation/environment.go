package simulation

import (
	"strings"

	shp "github.com/jonas-p/go-shp"
	log "github.com/sirupsen/logrus"
)

// Environment models the road network for the traffic simulation.
type Environment struct {
	// waypoints store the locations of the road network
	waypoints []Vector

	// Logger is used to give a context based log to the stdout
	Logger *log.Entry
}

// NewEnvironment creates a new Environment and intilises its variables.
func NewEnvironment() Environment {
	var env Environment

	// Setup the logger
	env.Logger = log.WithFields(log.Fields{
		"package": "simulation",
		"section": "Env"})

	return env
}

// ReadShapefile takes a shape file and sets up the environment.
func (e *Environment) ReadShapefile(fileName string) {
	shape, err := shp.Open(fileName)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer shape.Close()

	// Get the fields from the attribute table (DBF)
	fields := shape.Fields()

	// Loop through all the features in the shapefile
	for shape.Next() {
		n, p := shape.Shape()

		for k := range fields {
			val := shape.ReadAttribute(n, k)

			// Add the waypoints to the environment
			if strings.Contains(val, "waypoint") {
				e.waypoints = append(e.waypoints, Vector{x: p.BBox().MaxX, y: p.BBox().MaxY})
			}
		}
	}
	e.Logger.Debug(e.waypoints)
}

// WriteShapeFile writes a test shape file.
func (e *Environment) WriteShapeFile(fileName string) {
	// points to write
	points := []shp.Point{
		shp.Point{0.0, 0.0},
		shp.Point{30.0, 30.0},
	}

	// fields to write
	fields := []shp.Field{
		// String attribute field with length 25
		shp.StringField("name", 25),
	}

	// create and open a shapefile for writing points
	shape, err := shp.Create(fileName, shp.POINT)
	if err != nil {
		log.Fatal(err)
	}
	defer shape.Close()

	// setup fields for attributes
	shape.SetFields(fields)

	// write points and attributes
	for n, point := range points {
		shape.Write(&point)

		// write attribute for object n for field 0 (NAME)
		shape.WriteAttribute(n, 0, "waypoint")
	}
}
