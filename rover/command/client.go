package command

import "yonder/rover/common"

//Client defines the client-side API for interacting with the Command node
type Client interface {
	//DirectDrive relays direct motor control information to the drive system
	DirectDrive(common.DirectDrive) error

	//DirectArm relays direct motor control information to the arm system
	DirectArm(common.DirectArm) error

	//JoystickDrive accepts a polar coordinate "joystick" as input for the
	//drive system
	JoystickDrive(common.Joystick) error

	//TankDrive accepts a pair of scalars as "tank drive" input for the
	//drive system
	TankDrive(float32, float32) error
}
