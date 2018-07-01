package protocol

import "io"

//Protocol is an interface to the device used to talk to another program
//could be anything: CAN, I2C, SPI, UART, Ethernet, etc
type Protocol interface {
	//Create initializes any resources needed by the Protocol
	Start() error

	//Address sets the target device of the Protocol, if supported
	//Address may or may not be expensive, depending on the protocol used
	//ex: opening and closing TCP ports vs changing an I2C address
	Address(string) error

	//Read and Write bytes using this Protocol to communicate with the device
	//at the target address
	Write([]byte) (int, error)

	//Receiver sets the protocol's receiver, which is an io.Writer that accepts
	//data sent to the host using the protocol
	Receiver(io.Writer)

	//Destroy releases any resources held by the Protocol
	Finish() error
}
