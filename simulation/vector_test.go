package simulation

import (
	"math"
	"testing"
)

// TestNewVector tests if the correct Vector is created
// witht the specified inputs.
func TestNewVector(t *testing.T) {

	testValues := [][]float64{
		{-math.MaxFloat64, -math.MaxFloat64},
		{-math.MaxFloat64 + 1, -math.MaxFloat64},
		{0, -math.MaxFloat64},
		{math.MaxFloat64 - 1, -math.MaxFloat64},
		{math.MaxFloat64, -math.MaxFloat64},

		{-math.MaxFloat64, -math.MaxFloat64 + 1},
		{-math.MaxFloat64 + 1, -math.MaxFloat64 + 1},
		{0, -math.MaxFloat64 + 1},
		{math.MaxFloat64 - 1, -math.MaxFloat64 + 1},
		{math.MaxFloat64, -math.MaxFloat64 + 1},

		{-math.MaxFloat64, 0},
		{-math.MaxFloat64 + 1, 0},
		{0, 0},
		{math.MaxFloat64 - 1, 0},
		{math.MaxFloat64, 0},

		{-math.MaxFloat64, math.MaxFloat64 - 1},
		{-math.MaxFloat64 + 1, math.MaxFloat64 - 1},
		{0, math.MaxFloat64 - 1},
		{math.MaxFloat64 - 1, math.MaxFloat64 - 1},
		{math.MaxFloat64, math.MaxFloat64 - 1},

		{-math.MaxFloat64, math.MaxFloat64},
		{-math.MaxFloat64 + 1, math.MaxFloat64},
		{0, math.MaxFloat64},
		{math.MaxFloat64 - 1, math.MaxFloat64},
		{math.MaxFloat64, math.MaxFloat64},
	}

	for _, testInputs := range testValues {
		v := NewVector(testInputs[0], testInputs[1])

		if v.x != testInputs[0] || v.y != testInputs[1] {
			t.Errorf("NewVector(%v, %v): %v, Input: %v", testInputs[0], testInputs[1], v, testInputs)
		}
	}
}

// TestMagnitude tests that the Vecor function Magnitude
// correctly calculates the magnitude of the vector.
func TestMagnitude(t *testing.T) {
	// testValues hold inputs and expected results
	// inputs = [i][0:1], exepcted = [i][2]
	testValues := [][]float64{
		{0, 0, 0},
		{1, 0, 1},
		{0, 1, 1},
		{-1, 0, 1},
		{0, -1, 1},
		{1, 2, math.Sqrt(5)},
		{2, 1, math.Sqrt(5)},
		{10, 24, 26},
		{24, 10, 26},
	}

	for _, values := range testValues {
		v := NewVector(values[0], values[1])
		if got := v.Magnitude(); got != values[2] {
			t.Errorf("Magnitude(%v, %v) = %v, want %v", values[0], values[1], got, values[2])
		}
	}
}

// TestInRange checks if the vector correctly deduces if
// the position is within a margin of the target.
func TestInRange(t *testing.T) {
	type testCase struct {
		// pos is the position of the vector
		pos Vector
		// target is the position of the target vector
		target Vector
		// margin is the gap that can be allowed
		// between the vectors
		margin float64
		// expected is the correct result of the function
		expected bool
	}

	testValues := []testCase{
		testCase{
			pos:      NewVector(100, 100),
			target:   NewVector(0, 0),
			margin:   0.0,
			expected: false,
		},
		testCase{
			pos:      NewVector(0, 0),
			target:   NewVector(100, 100),
			margin:   0.0,
			expected: false,
		},
		testCase{
			pos:      NewVector(0, 0),
			target:   NewVector(0, 0),
			margin:   0.0,
			expected: true,
		},
		testCase{
			pos:      NewVector(1, 0),
			target:   NewVector(0, 0),
			margin:   1.0,
			expected: true,
		},
		testCase{
			pos:      NewVector(0, 1),
			target:   NewVector(0, 0),
			margin:   1.0,
			expected: true,
		},
		testCase{
			pos:      NewVector(10, 10),
			target:   NewVector(0, 0),
			margin:   15.0,
			expected: true,
		},
		testCase{
			pos:      NewVector(0, 0),
			target:   NewVector(10, 10),
			margin:   15.0,
			expected: true,
		},
		testCase{
			pos:      NewVector(0, 0),
			target:   NewVector(10, 10),
			margin:   1.0,
			expected: false,
		},
		testCase{
			pos:      NewVector(10, 10),
			target:   NewVector(0, 0),
			margin:   1.0,
			expected: false,
		},
	}

	for _, testCase := range testValues {
		if got := testCase.pos.InRange(testCase.target, testCase.margin); got != testCase.expected {
			t.Errorf("InRange(%v, %v, %v) = %v, expected = %v", testCase.pos, testCase.target, testCase.margin, got, testCase.expected)
		}
	}

}
