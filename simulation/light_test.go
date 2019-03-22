package simulation

import "testing"

// TestNewLight test that the NewLight function correctly
// returns a Light with the paramaters specified.
func TestNewLight(t *testing.T) {
	// Initial test values
	testID := 0
	testPos := Vector{0, 0}
	testStop := true

	l := NewLight(testID, testPos, testStop)

	if testID != l.id ||
		testPos != l.position ||
		testStop != l.stop ||
		l.Logger == nil {
		t.Errorf("NewLight(%v, %v, %v): %v, want %v, %v, %v", testID, testPos, testStop, l, testID, testPos, testStop)
	}

	// Change the ID
	testID = 1

	l = NewLight(testID, testPos, testStop)

	if testID != l.id ||
		testPos != l.position ||
		testStop != l.stop ||
		l.Logger == nil {
		t.Errorf("NewLight(%v, %v, %v): %v, want %v, %v, %v", testID, testPos, testStop, l, testID, testPos, testStop)
	}

	// Change the position
	testPos = Vector{1, 1}

	l = NewLight(testID, testPos, testStop)

	if testID != l.id ||
		testPos != l.position ||
		testStop != l.stop ||
		l.Logger == nil {
		t.Errorf("NewLight(%v, %v, %v): %v, want %v, %v, %v", testID, testPos, testStop, l, testID, testPos, testStop)
	}

	// Change stop value
	testStop = false

	l = NewLight(testID, testPos, testStop)

	if testID != l.id ||
		testPos != l.position ||
		testStop != l.stop ||
		l.Logger == nil {
		t.Errorf("NewLight(%v, %v, %v): %v, want %v, %v, %v", testID, testPos, testStop, l, testID, testPos, testStop)
	}
}

// TestGetID tests that the traffic light's function GetID
// returns the correct ID of the light.
func TestGetID(t *testing.T) {

	// test values is a map with the key the ID and
	// the value if the ID is correct
	testValues := map[int]bool{
		-1: false,
		0:  true,
		1:  false,
	}

	mockLight := Light{id: 0}

	for testID, testPass := range testValues {
		if pass := mockLight.GetID() == testID; pass != testPass {
			t.Errorf("GetID(): %v, Input: %v, GetID() == Input: %v", mockLight.GetID(), testID, pass)
		}
	}
}

// TestGetStop tests that the traffic light's function GetStop
// returns the correct stop value of the light.
func TestGetStop(t *testing.T) {
	mockLight := Light{stop: true}
	// Check that the mockLight.GetStop() is true
	if !mockLight.GetStop() {
		t.Errorf("GetStop() == %v, wan %v", mockLight.GetStop(), true)
	}

	mockLight = Light{stop: false}
	// Check that the mockLight.GetStop() is true
	if mockLight.GetStop() {
		t.Errorf("GetStop() == %v, wan %v", mockLight.GetStop(), false)
	}
}

// TestSetStop check to see if the Light's SetStop function
// correctly changes the value of stop.
func TestSetStop(t *testing.T) {

	// testValues is the list of inputs for the function
	testValues := []bool{
		false, true, false,
	}

	// setup the mock Light struct
	mockLight := NewLight(0, Vector{}, true)

	for _, testValue := range testValues {

		mockLight.SetStop(testValue)

		// check the stop value has been updated
		if mockLight.stop != testValue {
			t.Errorf("SetStop(%v) = %v, want %v", testValue, mockLight.stop, testValue)
		}
	}

}

// TestGetPosition tests that the traffic light's function GetPosition
// returns the correct position of the light.
func TestGetPosition(t *testing.T) {

	// test values is a map with the key the Position and
	// the value if the position is correct.
	testValues := map[Vector]bool{
		Vector{0, 0}: true,
		Vector{0, 1}: false,
		Vector{1, 0}: false,
		Vector{1, 1}: false,
	}

	mockLight := Light{position: Vector{0, 0}}

	for testPos, testPass := range testValues {
		if pass := mockLight.GetPosition().x == testPos.x &&
			mockLight.GetPosition().y == testPos.y; pass != testPass {
			t.Errorf("GetPosition(): %v, Input: %v, GetPosition() == Input: %v", mockLight.GetPosition(), testPos, pass)
		}
	}
}
