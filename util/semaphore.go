package util

//Semaphore implementation
type Semaphore chan empty

//NewSemaphore creates a Semaphore instance with a given capacity
func NewSemaphore(capacity int) Semaphore {
	sem := make(Semaphore, capacity)
	return sem
}

//P waits for and consumes one resource from the Semaphore
func (sem Semaphore) P() {
	sem <- empty{}
}

//V releases one resource from the Semaphore
func (sem Semaphore) V() {
	<-sem
}
