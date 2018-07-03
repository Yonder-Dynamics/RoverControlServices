package common

import "yonder/util"

//DirectDrive represents a message for direct drive motor control (dir/pwm)
type DirectDrive interface {
	util.Serializable
}

//DirectArm represents a message for direct arm motor control (dir/pwm)
type DirectArm interface {
	util.Serializable
}
