package simulation

import "math"

// Vector stores the x and y values.
type Vector struct {
	x float64
	y float64
}

// NewVector creates an instance of a Vector.
func NewVector(x float64, y float64) Vector {
	v := Vector{}
	v.x = x
	v.y = y
	return v
}

// Magnitude calclates the overall vector.
func (v *Vector) Magnitude() float64 {
	return math.Sqrt((v.x * v.x) + (v.y * v.y))
}

// InRange checks to see if the vectoe is within a certain margin of
// another vector.
func (v *Vector) InRange(target Vector, margin float64) bool {
	target.x -= v.x
	target.y -= v.y
	return target.Magnitude() <= margin
}

// ConvertToSlice changes the vector into slice of float64s.
func (v *Vector) ConvertToSlice() []float64 {
	var s []float64
	s = append(s, v.x)
	s = append(s, v.y)
	return s
}

// Equals checks if the given vector has the same values.
func (v *Vector) Equals(other Vector) bool {
	return v.x == other.x && v.y == other.y
}

// DistanceTo calculates the distance to a target.
func (v *Vector) DistanceTo(target Vector) float64 {
	dx := v.x - target.x
	dy := v.y - target.y
	return math.Sqrt((dx * dx) + (dy * dy))
}
