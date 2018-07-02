package conduits

import (
	"bytes"
	"net"
	"yonder/fabric"
	"yonder/util"
)

//TCPIn implements the Conduit interface, linking two nodes via TCP
type TCPIn struct {
	ln  net.Listener
	out chan<- fabric.Transaction
}

type tcpTransaction struct {
	buf  bytes.Buffer
	conn net.Conn
	end  util.Signal
}

//NewTCPIn creates a TCPIn instance suitable for recieving TCP input
func NewTCPIn(addr string) (*TCPIn, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	tcp := &TCPIn{ln, make(chan fabric.Transaction)}
	return tcp, nil
}

//In implements the In() method specified by the Conduit interface
//since a TCPIn conduit is meant to carry external inputs, attempting to link
//the input of the conduit will panic
func (tcp *TCPIn) In(<-chan fabric.Transaction) {
	panic("TCPIn Conduits do not accept direct input!")
}

//Out implements the Out() method specified by the Conduit interface
//the TCPIn instance will write all data received on its port to the given
//io.Writer
func (tcp *TCPIn) Out(t chan<- fabric.Transaction) {
	tcp.out = t
	go tcp.serve() //start serving once an output has been registered
}

func (tcp *TCPIn) serve() {
	for {
		conn, err := tcp.ln.Accept()
		if err != nil {
			//listener was closed
			return
		}
		go tcp.connect(conn)
	}
}

func (tcp *TCPIn) connect(conn net.Conn) {
	buf := bytes.Buffer{}
	transaction := &tcpTransaction{buf, conn, make(util.Signal, 1)}
	for {
		_, err := buf.ReadFrom(conn)
		if err != nil { //connection was closed
			return
		}
		tcp.out <- transaction //block until the transaction is received
		transaction.wait()     //block until the transaction has ended
	}
}

//Read from the internal buffer of the tcpTransaction
func (t *tcpTransaction) Read(b []byte) (int, error) {
	return t.buf.Read(b)
}

//Write to the connection represented by the tcpTransaction
func (t *tcpTransaction) Write(b []byte) (int, error) {
	return t.conn.Write(b)
}

//End writes to the internal end channel, allowing a waiting TCP conduit to
//produce further transactions
func (t *tcpTransaction) End() {
	t.end <- util.Empty{}
}

//wait blocks until the tcpTransaction has been ended
func (t *tcpTransaction) wait() {
	<-t.end
}
