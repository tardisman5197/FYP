package simulation

import "testing"

// TestUpdateSpeed checks to see if the vehicle correctly changes
// their speed depending on the senario.
func TestUpdateSpeed(t *testing.T) {
	// Branch / Decision Coverage - two test cases for each decision.
	// Branches:
	//	1. Slow down for lights
	//	2. Slow down to touch waypoint
	// 	3. Slow down due to other agents
	//	4. Accelerate if space

	// 1. Slow down for lights
	// Prerequisite:
	// * The vehicle has to be close to a light
	// * The vehicle's speed must result in it passing the light
	// * The traffic light must be in a stop state
	// Result:
	// * The vehicle will not pass the light
	// * The vehicle's speed will decrease.

	// Setup the simulation
	// Create environment to test the agent in
	var env Environment
	env = NewEnvironment()
	env.AddLight(NewVector(10, 10), true)

	// Create a list of agents to populate the simulation with
	var agents []Agent

	// Test case 1

	// Setup test vehicle
	testVehicle := NewVehicle(
		0,
		NewVector(0, 0),
		20,
		30,
		5,
		10,
		[]Vector{
			NewVector(10, 10),
			NewVector(20, 20),
		},
		0,
	)

	// Perform function
	testVehicle.updateSpeed(agents, env)

	// Check if the vehicle has not decreased its speed
	// and has passed the light
	if testVehicle.GetSpeed() >= 20 &&
		testVehicle.GetPosition().x > 10 &&
		testVehicle.GetPosition().y > 10 {

		t.Errorf("UpdateSpeed() - Not slowing down for light - Speed: %v !< %v || Position: %v !< %v",
			testVehicle.GetSpeed(),
			20,
			testVehicle.GetPosition(),
			[]float64{10, 10},
		)
	}

	// Test case 2
	// Setup test vehicle
	testVehicle = NewVehicle(
		0,
		NewVector(5, 5),
		10,
		30,
		5,
		10,
		[]Vector{
			NewVector(10, 10),
			NewVector(20, 20),
		},
		0,
	)

	// Perform function
	testVehicle.updateSpeed(agents, env)

	// Check if the vehicle has not decreased its speed
	// and has passed the light
	if testVehicle.GetSpeed() >= 10 &&
		testVehicle.GetPosition().x > 10 &&
		testVehicle.GetPosition().y > 10 {

		t.Errorf("UpdateSpeed() - Not slowing down for light - Speed: %v !< %v || Position: %v !< %v",
			testVehicle.GetSpeed(),
			10,
			testVehicle.GetPosition(),
			[]float64{10, 10},
		)
	}

	// 2. Slow down to touch waypoint
	// Prerequisite:
	//  * The vehicle must be approaching a waypoint
	//  * The vehicle's speed must result in it passing the waypoint
	//	* There is no traffic light at the waypoint
	// Result:
	//  * The vehicle's speed will decrease

	// Setup the simulation
	// Create environment to test the agent in
	env = NewEnvironment()

	// Test case 1

	// Setup test vehicle
	testVehicle = NewVehicle(
		0,
		NewVector(0, 0),
		20,
		30,
		5,
		10,
		[]Vector{
			NewVector(10, 10),
			NewVector(20, 20),
		},
		0,
	)

	// Perform function
	testVehicle.updateSpeed(agents, env)

	// Check if the vehicle has not decreased its speed
	// and has passed the light
	if testVehicle.GetSpeed() >= 20 &&
		testVehicle.GetPosition().x > 10 &&
		testVehicle.GetPosition().y > 10 {

		t.Errorf("UpdateSpeed() - Not slowing down for waypoint - Speed: %v !< %v || Position: %v !< %v",
			testVehicle.GetSpeed(),
			20,
			testVehicle.GetPosition(),
			[]float64{10, 10},
		)
	}

	// Test case 2

	// Setup test vehicle
	testVehicle = NewVehicle(
		0,
		NewVector(7, 7),
		10,
		30,
		5,
		10,
		[]Vector{
			NewVector(10, 10),
			NewVector(20, 20),
		},
		0,
	)

	// Perform function
	testVehicle.updateSpeed(agents, env)

	// Check if the vehicle has not decreased its speed
	// and has passed the light
	if testVehicle.GetSpeed() >= 10 &&
		testVehicle.GetPosition().x > 10 &&
		testVehicle.GetPosition().y > 10 {

		t.Errorf("UpdateSpeed() - Not slowing down for waypoint - Speed: %v !< %v || Position: %v !< %v",
			testVehicle.GetSpeed(),
			10,
			testVehicle.GetPosition(),
			[]float64{10, 10},
		)
	}

	// 3. Slow down due to other agents
	// Prerequisite:
	//  * The vehicle is close to another agent
	//  * The vehicle's speed would result in it crashing into
	//    the other agent
	//  * There are no waypoints or traffic lights inbetween the
	//    vehicle and agent
	// Result:
	//  * The vehicle's speed will decrease

	// Setup the simulation
	// Create environment to test the agent in
	env = NewEnvironment()

	// Test case 1

	// Add other agent
	agents = []Agent{
		NewVehicle(
			0,
			NewVector(4, 4),
			5,
			5,
			3,
			3,
			[]Vector{
				NewVector(10, 10),
				NewVector(20, 20),
			},
			0,
		),
	}

	// Setup test vehicle
	testVehicle = NewVehicle(
		1,
		NewVector(0, 0),
		10,
		15,
		3,
		5,
		[]Vector{
			NewVector(10, 10),
			NewVector(20, 20),
		},
		0,
	)

	// Perform function
	testVehicle.updateSpeed(agents, env)

	// Check if the vehicle has not decreased its speed
	// and has passed the light
	if testVehicle.GetSpeed() >= 10 &&
		testVehicle.GetPosition().x > 4 &&
		testVehicle.GetPosition().y > 4 {

		t.Errorf("UpdateSpeed() - Not slowing down for other agent - Speed: %v !< %v || Position: %v !< %v",
			testVehicle.GetSpeed(),
			10,
			testVehicle.GetPosition(),
			[]float64{4, 4},
		)
	}

	// Test case 2

	// Add other agent
	agents = []Agent{
		NewVehicle(
			0,
			NewVector(7, 7),
			5,
			5,
			3,
			3,
			[]Vector{
				NewVector(10, 10),
				NewVector(20, 20),
			},
			0,
		),
	}

	// Setup test vehicle
	testVehicle = NewVehicle(
		1,
		NewVector(1, 1),
		20,
		30,
		5,
		8,
		[]Vector{
			NewVector(10, 10),
			NewVector(20, 20),
		},
		0,
	)

	// Perform function
	testVehicle.updateSpeed(agents, env)

	// Check if the vehicle has not decreased its speed
	// and has passed the light
	if testVehicle.GetSpeed() >= 10 &&
		testVehicle.GetPosition().x > 7 &&
		testVehicle.GetPosition().y > 7 {

		t.Errorf("UpdateSpeed() - Not slowing down for other agent - Speed: %v !< %v || Position: %v !< %v",
			testVehicle.GetSpeed(),
			20,
			testVehicle.GetPosition(),
			[]float64{7, 7},
		)
	}

	// 4. Accelerate if space
	// Prerequisite:
	//  * There are no waypoints or agents close to the agent
	// Result:
	//  * The vehicle's speed will increase

	// Setup the simulation
	// Create environment to test the agent in
	env = NewEnvironment()
	agents = []Agent{}

	// Test case 1

	// Setup test vehicle
	testVehicle = NewVehicle(
		1,
		NewVector(0, 0),
		5,
		15,
		5,
		5,
		[]Vector{
			NewVector(20, 20),
		},
		0,
	)

	// Perform function
	testVehicle.updateSpeed(agents, env)

	// Check if the vehicle has not decreased its speed
	// and has passed the light
	if testVehicle.GetSpeed() < 5 &&
		testVehicle.GetPosition().x <= 0 &&
		testVehicle.GetPosition().y <= 0 {

		t.Errorf("UpdateSpeed() - Not Speeding up when there is space - Speed: %v !> %v || Position: %v !> %v",
			testVehicle.GetSpeed(),
			5,
			testVehicle.GetPosition(),
			[]float64{0, 0},
		)
	}

	// Test case 2

	// Setup test vehicle
	testVehicle = NewVehicle(
		1,
		NewVector(5, 5),
		10,
		15,
		5,
		5,
		[]Vector{
			NewVector(20, 20),
		},
		0,
	)

	// Perform function
	testVehicle.updateSpeed(agents, env)

	// Check if the vehicle has not decreased its speed
	// and has passed the light
	if testVehicle.GetSpeed() < 10 &&
		testVehicle.GetPosition().x <= 5 &&
		testVehicle.GetPosition().y <= 5 {

		t.Errorf("UpdateSpeed() - Not Speeding up when there is space - Speed: %v !> %v || Position: %v !> %v",
			testVehicle.GetSpeed(),
			10,
			testVehicle.GetPosition(),
			[]float64{5, 5},
		)
	}
}
