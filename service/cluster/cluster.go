package cluster

//Cluster is a set of nodes and a router which act as a single service
type Cluster interface {
	start()
	stop()
}
