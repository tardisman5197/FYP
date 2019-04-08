package controller

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"../simulation"
	log "github.com/sirupsen/logrus"
)

type testInfo struct {
	Position []float64 `json:"position"`
	Stop     bool      `json:"stop"`
}

type testResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// TestAddLight checks to see if the endpoint to add a traffic
// light to a simualtion correctly responds to the request.
func TestAddLight(t *testing.T) {
	// Robust BVA
	// For each inputs:
	//  2x invalid inputs
	//  2x edge cases
	//  2x nominal

	// Input: Position
	// 2x invalid
	testRequest := testInfo{
		Position: []float64{},
		Stop:     true,
	}
	response := sendTestAddLightRequest(testRequest)
	if response.Success != false {
		t.Errorf(
			"addLight() - success: %v != %v",
			response.Success,
			false,
		)
	}

	testRequest = testInfo{
		Position: []float64{1, 1, 1},
		Stop:     true,
	}
	response = sendTestAddLightRequest(testRequest)
	if response.Success != false {
		t.Errorf(
			"addLight() - success: %v != %v",
			response.Success,
			false,
		)
	}

	// 2x edge cases
	testRequest = testInfo{
		Position: []float64{-math.MaxFloat64, -math.MaxFloat64},
		Stop:     true,
	}
	response = sendTestAddLightRequest(testRequest)
	if response.Success != true {
		t.Errorf(
			"addLight() - success: %v != %v",
			response.Success,
			true,
		)
	}

	testRequest = testInfo{
		Position: []float64{math.MaxFloat64, math.MaxFloat64},
		Stop:     true,
	}
	response = sendTestAddLightRequest(testRequest)
	if response.Success != true {
		t.Errorf(
			"addLight() - success: %v != %v",
			response.Success,
			true,
		)
	}

	// 2x nominal cases
	testRequest = testInfo{
		Position: []float64{10, 10},
		Stop:     true,
	}
	response = sendTestAddLightRequest(testRequest)
	if response.Success != true {
		t.Errorf(
			"addLight() - success: %v != %v",
			response.Success,
			true,
		)
	}

	testRequest = testInfo{
		Position: []float64{20.5, 20.5},
		Stop:     true,
	}
	response = sendTestAddLightRequest(testRequest)
	if response.Success != true {
		t.Errorf(
			"addLight() - success: %v != %v",
			response.Success,
			true,
		)
	}

	// Input: Stop
	testRequest = testInfo{
		Position: []float64{20.5, 20.5},
		Stop:     true,
	}
	response = sendTestAddLightRequest(testRequest)
	if response.Success != true {
		t.Errorf(
			"addLight() - success: %v != %v",
			response.Success,
			true,
		)
	}

	testRequest = testInfo{
		Position: []float64{20.5, 20.5},
		Stop:     false,
	}
	response = sendTestAddLightRequest(testRequest)
	if response.Success != true {
		t.Errorf(
			"addLight() - success: %v != %v",
			response.Success,
			true,
		)
	}

}

func testRouter() *mux.Router {
	var c Controller
	c.Logger = log.WithFields(log.Fields{"package": "controller"})
	c.simulations.Store("test", simulation.NewSimulation(simulation.NewEnvironment()))

	r := mux.NewRouter()
	r.HandleFunc("/simulation/light/add/{id}", c.addLight).Methods("POST")
	return r
}

func sendTestAddLightRequest(testRequest testInfo) testResponse {
	testJSONStr, _ := json.Marshal(testRequest)

	// Create the test request
	req := httptest.NewRequest(
		"POST",
		"/simulation/light/add/test",
		bytes.NewReader(testJSONStr),
	)
	req.Header.Set("Content-Type", "application/json; param=value")
	w := httptest.NewRecorder()

	// Send request
	testRouter().ServeHTTP(w, req)

	// Decode response
	var r testResponse
	_ = json.NewDecoder(w.Body).Decode(&r)
	return r
}
