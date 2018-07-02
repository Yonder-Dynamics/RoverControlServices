package fabric

import "io"

//Router acts as the I/O hub for any Go program which requires IPC
//routers are responsible for abstracting I/O details away from the core program
//to allows communiation details to be changed in the backgroud
type Router interface {
	//In registers an input to the router, which can be linked to an output
	//to create a pipeline
	//takes the input name, and an io.Writer to receive the input
	In(string, io.Writer)

	//Out registers an output to the router, which can be linked to an input
	//to create a pipeline
	//returns an io.Writer to send output on
	Out(string) io.Writer

	//Addr returns the address of this Router
	Addr() string

	//Link connects the given Conduit to the named Router I/O writer
	Link(string, Conduit)

	//Close shuts down the Router
	Close() error
}

//Conduit transfers data from the output of one Router to the input of another
type Conduit interface {
	//In accepts a channel for receiving Transaction objects, which are used to
	//deliver data and return the response
	In(<-chan Transaction)

	//Out accepts a channel for sending Transaction objects, which are used to
	//receive data and send a response
	Out(chan<- Transaction)
}

//Transaction instances represent an input-ouput exchange
type Transaction interface {
	io.ReadWriter

	//End finishes the exchange and allows for a new Transaction to occur
	End()
}
