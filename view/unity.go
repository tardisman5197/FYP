package view

import (
	"encoding/json"
	"net"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// message is struct containg the information that needs to be sent
// to the unity application.
type message struct {
	// Agents stores the x, y coordinates of each agent
	Agents []vector2 `json:"agents"`
	// Waypoints stores the x,y coordinates of the waypoints
	Waypoints []vector2 `json:"waypoints"`
	// Tick stores the tick of the simulation which the
	// information represents
	Tick int `json:"tick"`
}

// receipt is a struct that stores the information that the unity
// application sends.
type receipt struct {
	// Filepath stores the path to the image that represents a tick
	// in the simulation
	Filepath string `json:"filepath"`
}

// vector2 struct is used to convert x,y valeus into a json map format.
type vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// UnityServer handels the communication between the go application and
// the unity application.
type UnityServer struct {
	// port is the port that the tcp server opens e.g. ":6666"
	port string
	// conn is the object that stores the information about the client
	// connected to the server
	conn net.Conn
	// connected shows if the server is connected to the unity application
	connected bool
	// incoming stores any messages received from the unity app
	incoming chan string
	// outgoing stores messages that need to be sent to the unity app
	outgoing chan string
	// stop can be set to true of the TCP server needs to be stopped
	stop chan bool

	// Logger is used to give a context based log to the stdout
	Logger *log.Entry
}

// NewUnityServer creates a new TCP server for the unity application to
// communicate to.
func NewUnityServer(port string) UnityServer {
	u := UnityServer{}
	u.Logger = log.WithFields(log.Fields{
		"package": "view",
		"section": "unityServer"})
	u.port = port
	u.incoming = make(chan string)
	u.outgoing = make(chan string)
	u.stop = make(chan bool)
	return u
}

// startUnityApp runs the unity application executable.
func (u *UnityServer) startUnityApp() {
	if pathToUnity != "" {
		cmd := exec.Command(pathToUnity)
		cmd.Run()
	} else {
		u.Logger.Warn("Path to unity not set")
	}
}

// StartServer creats a TCP server which listens for new connections
func (u *UnityServer) StartServer() error {
	u.Logger.Info("Starting Unity server")
	u.Logger.Infof("Port: '%v'", u.port)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", u.port)
	if err != nil {
		u.Logger.Error("Error: Unable to resolve TCP addr")
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		u.Logger.Error("Error: Unable to listen on tcp addr")
		return err
	}

	startedUnity := false

	for {
		if !startedUnity {
			u.Logger.Info("Starting Unity App")
			go u.startUnityApp()
			startedUnity = true
		}

		u.conn, err = listener.Accept()
		if err != nil {
			continue
		}
		if err == nil {
			u.connected = true
			u.Logger.Infof("Connection Established: %v", u.conn.RemoteAddr())
			go u.handleClient()
			break
		}
	}
	return nil
}

// handelClient checks for messages to be sent or received
func (u *UnityServer) handleClient() {
	// Disconnect from the server after the loop has been broken out of
	defer u.disconnect()

	// Constantly check for messages from unity
	go u.readMessage()

handlerLoop:
	for {
		select {
		// Check if the server should stop
		case stop := <-u.stop:
			if stop {
				break handlerLoop
			}
		// Check if their are any messages to send
		case msg := <-u.outgoing:
			u.writeToSocket(msg)
		// Check if their are any incoming messages
		case msg := <-u.incoming:
			u.parseMessage(msg)
		}
	}
}

// SendMessage adds a message to be sent to the unbity application.
func (u *UnityServer) SendMessage(msg string) {
	u.Logger.Debugf("Adding to outgoing: %v", msg)
	u.outgoing <- msg
}

// SendSimulation creates a json string and sends it to the unity application.
func (u *UnityServer) SendSimulation(agents, waypoints [][]float64, tick int) {
	// Convert agents [][]float64 into []vector2
	var agentVec []vector2
	for i := 0; i < len(agents); i++ {
		agentVec = append(agentVec, vector2{X: agents[i][0], Y: agents[i][1]})
	}

	// Convert waypoints [][]float64 into []vector2
	var waypointVec []vector2
	for i := 0; i < len(waypoints); i++ {
		waypointVec = append(waypointVec, vector2{X: waypoints[i][0], Y: waypoints[i][1]})
	}

	// convert message to json string
	jsonStr, err := json.Marshal(message{
		Agents:    agentVec,
		Waypoints: waypointVec,
		Tick:      tick,
	})
	if err != nil {
		u.Logger.Error("Error: Converting to JSON")
		return
	}

	// Send the json string to unity
	u.SendMessage(string(jsonStr))
}

// writeToSocket writes a messaage to the socket file for the unity
// application to read.
func (u *UnityServer) writeToSocket(msg string) {
	// Write tries to write the message to the file
	// if an EOF error is found the connection has been closed
	_, err := u.conn.Write([]byte(msg))
	if err != nil {
		u.Logger.Error("Error: sending message")
		return
	}
	u.Logger.Debugf("Messsage sent: %v", msg)
}

// readMessage checks for messages from the unity application.
func (u *UnityServer) readMessage() {
	for {
		data := make([]byte, 1024)
		n, err := u.conn.Read(data)
		if err != nil {
			u.Logger.Errorf("Error: Reading socket - %v", err)
			u.stop <- true
			break
		}
		u.incoming <- string(data[:n])
	}
}

// parseMessage reads in messages sent from the unity application.
func (u *UnityServer) parseMessage(msg string) {
	u.Logger.Infof("Recieved: %v", msg)

	// Convert json to struct
	var r receipt
	err := json.Unmarshal([]byte(msg), &r)
	if err != nil {
		u.Logger.Errorf("Error: Trying to parse message - %v", err)
		return
	}

	// Check if the image exists
	if _, err := os.Stat(r.Filepath); err == nil {
		u.Logger.Debug("File does exist")

	}
}

// StopServer closes the communication between the server and unity
// application.
func (u *UnityServer) StopServer() {
	u.stop <- true
}

// disconnect cloeses the communication between the server and unity
// application and updates the connected value.
func (u *UnityServer) disconnect() {
	u.Logger.Info("Disconnecting unity")
	u.conn.Close()
	u.connected = false
}

// Connected returns the state of the server. If connected is true
// the server is connected to the unity application.
func (u *UnityServer) Connected() bool {
	return u.connected
}
