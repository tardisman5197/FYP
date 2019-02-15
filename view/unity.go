package view

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// UnityServer handels the communication between the go application and
// the unity application.
type UnityServer struct {
	// port is the port that the tcp server opens e.g. ":6666"
	port string
	// conn is the object that stores the information about the client
	// connected to the server
	conn      net.Conn
	connected bool
	incoming  chan string
	outgoing  chan string
	stop      chan bool

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

// StartServer ...
func (u *UnityServer) StartServer() {
	u.Logger.Info("Starting Unity server")
	u.Logger.Infof("Port: '%v'", u.port)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", u.port)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	for {
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
}

// handelClient ...
func (u *UnityServer) handleClient() {
	defer u.disconnect()

handlerLoop:
	for {
		select {
		case stop := <-u.stop:
			u.Logger.Infof("Stop Flag Unity server %v", stop)
			if stop {
				u.Logger.Info("Stopping Unity server")
				break handlerLoop
			}
		case msg := <-u.outgoing:
			u.Logger.Debugf("Sending: %v", msg)
			_, err := u.conn.Write([]byte(msg))
			if err != nil {
				panic(err)
			}
			u.Logger.Debugf("Messsage sent: %v", msg)
		}
	}
}

// SendMessage ...
func (u *UnityServer) SendMessage(msg string) {
	u.Logger.Debugf("Adding to outgoing: %v", msg)
	u.outgoing <- msg
}

// StopServer ...
func (u *UnityServer) StopServer() {
	u.Logger.Debugf("Setting stop flag to: %v", true)
	u.stop <- true
}

// Connected ...
func (u *UnityServer) Connected() bool {
	return u.connected
}

// disconnect ...
func (u *UnityServer) disconnect() {
	u.Logger.Info("Disconnecting unity")
	u.conn.Close()
	u.connected = false
	u.Logger.Info("Connection Closed")
}
