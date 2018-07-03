package status

import "yonder/rover"

//Client defines the client-side API for interacting with the Status node
type Client interface {
	//Update pulls all rover status information and writes it to a Rover object
	Update() (*rover.Rover, error)

	Odom() (*rover.Odom, error)

	Drive() (*rover.Drive, error)

	Arm() (*rover.Arm, error)

	Power() (*rover.Power, error)

	Network() (*rover.Network, error)

	Core() (*rover.Core, error)
}
