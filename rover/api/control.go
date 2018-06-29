package api

import (
	"fmt"
	"net"
)

//implements the control API defined for the rover
type ControlHandler struct {
	ln net.Listener
}

//Attach creates a worker thread for handling connections to the supplied Listener
func (f *ControlHandler) Attach(ln net.Listener) {
	go f.listen(ln)
}

func (f *ControlHandler) listen(ln net.Listener) {
	f.ln = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			break
		}
		fmt.Fprintf(conn, "hello there\n")
		conn.Close()
	}
}

//Detach closes the Listener attached to this ControlHandler
func (f *ControlHandler) Detach() error {
	err := f.ln.Close()
	if err != nil {
		return err
	}
	return nil
}
