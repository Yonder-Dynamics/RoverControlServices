package cluster

//Node is an object which implements a small sub-service within a helper thread
type Node interface {
	//start creates a worker thread to process requests to the node and
	//allocates any resources needed by the Node
	//returns the I/O channels used by the Node
	start() (in <-chan interface{}, out chan<- interface{})

	//stop signals the worker thread to terminate and releases any resources
	//held by the node
	stop()
}
