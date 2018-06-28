package cluster

//Router is an object which maps API calls into calls to its cluster's nodes
type Router interface {
	start()
	stop()
}
