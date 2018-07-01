package protocol

import (
	"io"
	"net"
	"yonder"
)

//TCPProtocol implements Protocol using TCP connections
//TCPProtocol supports concurrent connections, up to a given limit
type TCPProtocol struct {
	port        string
	target      string
	connections map[string]tcpConn
	listener    net.Listener
	receiver    io.Writer
	sem         yonder.Semaphore
}

type tcpConn struct {
	net.Conn
	available yonder.Signal
}

//NewTCPProtocol initializes a TCPProtocol with the proper fields set
func NewTCPProtocol(port string, numSimultaneous int) TCPProtocol {
	return TCPProtocol{port: port, sem: yonder.NewSemaphore(numSimultaneous)}
}

//Start starts to TCPProtocol by creating a thread to Listen for TCP connections
func (p *TCPProtocol) Start() error {
	listener, err := net.Listen("tcp", p.port)
	if err == nil {
		return err
	}
	go p.listen(listener)
	go p.selectConnection()
	return nil
}

func (p *TCPProtocol) listen(listener net.Listener) {
	p.listener = listener
	for {
		p.sem.P()
		conn, err := p.listener.Accept()
		if err != nil {
			return
		}

		//map the connection and start reading from it
		go p.connect(tcpConn{conn, make(yonder.Signal)})
	}
}

func (p *TCPProtocol) connect(conn tcpConn) {
	defer conn.finalize(p)
	p.connections[conn.RemoteAddr().String()] = conn
	read(conn)
}

func (conn *tcpConn) finalize(p *TCPProtocol) {
	conn.Close()
	delete(p.connections, conn.RemoteAddr().String())
}

//selectConnection waits for data on any of its connections
func (p *TCPProtocol) selectConnection() {
	empty := yonder.Empty{}
	bufsize := 100
	buffer := make([]byte, bufsize)
	var err error
	var n int
	var w int
	for {
		//iterate over connections and check if any have new data
		for c := range p.connections {
			select {
			case <-p.connections[c].available:
				//pipe bytes to this protocol's receiver
				n = bufsize
				for n == bufsize {
					n, err = p.connections[c].Read(buffer)
					if err != nil {
						//what do we do if there was an error
					}
					w, err = p.receiver.Write(buffer[:n])
					if w < n {
						//error writing to the receiver
					}
				}

				//signal that the data has been consumed
				p.connections[c].available <- empty
			default: //try the next connection
				continue
			}
		}
	}
}

func read(conn tcpConn) {
	buffer := make([]byte, 0)
	empty := yonder.Empty{}
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			return
		}
		conn.available <- empty //signal that a read is ready
		<-conn.available        //wait for something to consume the read
	}
}

//Finish closes the TCPListener associated with this protocol
func (p *TCPProtocol) Finish() error {
	//TODO: close any connections
	return p.listener.Close()
}

//Address changes the address the TCPListener will write to
func (p *TCPProtocol) Address(addr string) {
	p.target = addr
}
