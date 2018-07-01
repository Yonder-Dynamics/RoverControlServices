package main

import (
	"net"
	"yonder/rover/api"
)

//onboard rover control server executable
//listens on TCP ports for API calls and passes them to the appropriate handler
//APIs Mapped:
//	control/	-> api/control.go
//	status/		-> api/status.go

//Handler is an object whose API is exposed on a TCP port
type Handler interface {
	//attach Handler to a port via a TCPListener
	Attach(net.Listener)

	//detach Handler from its port
	Detach() error
}

func main() {
	controller := api.CommandHandler{}
	controllerListener, err := net.Listen("tcp", ":8081")
	if err != nil {
		return
	}
	controller.Attach(controllerListener)
	for {
	}
}
