package protocol

import (
	"io"
	"log"
	"net"
	"sync"
	"yonder"
)

const (
	//BUFSIZE is the size of the tcp read/write buffer
	BUFSIZE = 100
)

//TCPProtocol implements Protocol using TCP connections
//supports concurrent connections, up to a given limit
type TCPProtocol struct {
	port        string
	target      string
	connections map[string]tcpConn
	lock        sync.Mutex
	listener    net.Listener
	receiver    io.Writer
	sem         yonder.Semaphore
}

type tcpConn struct {
	net.Conn
	buffer    []byte
	available chan int
	consumed  yonder.Signal
}

//NewTCPProtocol initializes a TCPProtocol with the proper fields set
func NewTCPProtocol(port string, numSimultaneous int) *TCPProtocol {
	return &TCPProtocol{
		port:        port,
		sem:         yonder.NewSemaphore(numSimultaneous),
		connections: make(map[string]tcpConn),
		lock:        sync.Mutex{},
	}
}

//Receiver sets the receiver object for this protocol
func (p *TCPProtocol) Receiver(w io.Writer) {
	p.receiver = w
}

//Start starts to TCPProtocol by creating a thread to Listen for TCP connections
func (p *TCPProtocol) Start() error {
	log.Printf("Starting TCP Protocol on %s", p.port)
	listener, err := net.Listen("tcp", p.port)
	if err != nil {
		return err
	}
	go p.listen(listener)
	go p.selectConnection()
	return nil
}

func (p *TCPProtocol) listen(listener net.Listener) {
	log.Printf("Started Listener thread")
	p.listener = listener
	for {
		p.sem.P()
		conn, err := p.listener.Accept()
		if err != nil {
			return
		}

		//map the connection and start reading from it
		go p.connect(conn)
	}
}

func (p *TCPProtocol) connect(conn net.Conn) {
	log.Printf("Connecting to: %s", conn.RemoteAddr().String())
	tcpConn := tcpConn{
		conn,
		make([]byte, BUFSIZE),
		make(chan int),
		make(yonder.Signal),
	}
	defer tcpConn.finalize(p)
	p.lock.Lock()
	p.connections[conn.RemoteAddr().String()] = tcpConn
	p.lock.Unlock()
	log.Printf("Mapped %s into connections", conn.RemoteAddr().String())
	read(tcpConn)
}

func (conn *tcpConn) finalize(p *TCPProtocol) {
	log.Printf("Disconnecting from: %s", conn.RemoteAddr().String())
	conn.Close()
	p.lock.Lock()
	delete(p.connections, conn.RemoteAddr().String())
	p.lock.Unlock()
}

//selectConnection waits for data on any of its connections
func (p *TCPProtocol) selectConnection() {
	log.Printf("Started connection multiplexor thread")
	empty := yonder.Empty{}
	for {
		//iterate over connections and check if any have new data
		p.lock.Lock()
		// log.Print("Checking for incoming data")
		for _, conn := range p.connections {
			// log.Printf("Checking: %s", addr)
			select {
			case n := <-conn.available:
				//pipe bytes to this protocol's receiver
				bufsize := len(conn.buffer)
				w, err := p.receiver.Write(conn.buffer[:n])
				if w < n {
					log.Print(err)
				}
				for n == bufsize {
					// log.Printf("Reading data from: %s", addr)
					n, err = conn.Read(conn.buffer)
					// log.Printf("Read %d bytes from: %s", n, addr)
					if err != nil {
						//what do we do if there was an error
						log.Print(err)
					}
					// log.Printf("Writing data from: %s", addr)
					w, err = p.receiver.Write(conn.buffer[:n])
					if w < n {
						log.Print(err)
					}
				}

				//signal that the data has been consumed
				// log.Printf("Signalling data consumption for connection: %s", addr)
				conn.consumed <- empty
			default: //try the next connection
				continue
			}
		}
		// log.Print("Done checking for incoming data")
		p.lock.Unlock()
	}
}

func read(conn tcpConn) {
	for {
		n, err := conn.Read(conn.buffer)
		if err != nil {
			return
		}
		// log.Printf("Got data from %s: %d bytes", conn.RemoteAddr().String(), n)
		conn.available <- n //signal that a read is ready
		<-conn.consumed     //wait for something to consume the read
	}
}

//Finish closes the TCPListener associated with this protocol
func (p *TCPProtocol) Finish() error {
	log.Printf("Closing TCP Protocol on %s", p.listener.Addr().String())
	for _, conn := range p.connections {
		conn.Close()
	}
	return p.listener.Close()
}

//Address changes the address the TCPListener will write to
func (p *TCPProtocol) Address(addr string) error {
	log.Printf("Setting TCP Protocol target to: %s", addr)
	p.target = addr
	return nil
}

//Write writes the given data to the Protocol's target via TCP
//note: having the protocol write to itself will currently allocate a new
//port for sending the data, resulting in two ports being mapped
func (p *TCPProtocol) Write(b []byte) (int, error) {
	if conn, ok := p.connections[p.target]; ok {
		return conn.Write(b)
	}
	conn, err := net.Dial("tcp", p.target)
	if err != nil {
		return 0, err
	}
	go p.connect(conn)
	return conn.Write(b)
}
