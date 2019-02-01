package simulation

import "math"

// Vector stores the x and y values.
type Vector struct {
	x float64
	y float64
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
