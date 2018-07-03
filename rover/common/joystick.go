package common

import "yonder/util"

//Joystick defines a virtual joystick state for controlling rover systems
type Joystick interface {
	//Direction returns the angle of the stick
	//angle: radian clockwise angle from "forward"
	Direction() (angle float32)

	//Magnitude returns the distance of the stick from its resting position
	//val: value in [0,1] range, maps from rest position to max displacement
	Magnitude() (val float32)

	util.Serializable
}
