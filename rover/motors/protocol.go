package motors

//Protocol is an interface to the device used to talk to the microcontroller
//could be anything: CAN, I2C, SPI, UART, Ethernet, etc
type Protocol interface {
	//Create initializes any resources needed by the Protocol
	Create() error

	//Address sets the target device of the Protocol, if supported
	//Address may or may not be expensive, depending on the protocol used
	Address([]rune) error

	//Write send bytes using this Protocol to the device at the target address
	Write([]byte) (int, error)

	//Read bytes using this Protocol
	Read([]byte) (int, error)

	//Destroy releases any resources held by the Protocol
	Destroy() error
}
