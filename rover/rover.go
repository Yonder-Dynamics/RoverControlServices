package rover

//Rover encapsulates the full status of the rover being polled
type Rover struct {
	odom    Odom
	network Network
	core    Core
}

//Odom encapsulates the odometry information of the rover
type Odom struct {
	orientation Orientation
	position    Position
}

//Orientation describes the rover's angular position
type Orientation struct {
	//some representation of orientation (quaternion?)
}

//Position describes the rover's planar position
type Position struct {
	//estimated lat-long position
}

//Network describes the networks the rover is linked to
type Network struct {
	//status of connected computers/nodes
	antenna Antenna
}

//Antenna describes the antenna connection
type Antenna struct {
	//connection strength
	//connection latency
}

//Core describes the status of the onboard computer
type Core struct {
	//CPU usage and temp
	//RAM usage
	//program logs endpoint
}

//Drive describes the rover's drive system
type Drive struct {
	//wheel power use and efficiency
}

//Power describes the rover's power systems
type Power struct {
	//battery charge, if we have that info
	//approximate power usage
}

//Arm describes the rover's arm configuration
type Arm struct {
	//joint positions and lengths
}
