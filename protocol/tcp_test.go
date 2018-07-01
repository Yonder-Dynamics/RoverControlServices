package protocol

import (
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
)

type syncWriter struct {
	w    io.Writer
	done chan bool
}

func TestTCPProtocol(t *testing.T) {
	t.Log("Starting TestTCPProtocol")
	builder := strings.Builder{}
	writer := syncWriter{&builder, make(chan bool)}
	var tcp Protocol = NewTCPProtocol("127.0.0.1:8080", 5)
	tcp.Receiver(&writer)
	tcp.Start()

	t.Log("Testing Protocol Receiver()")
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Fatal(err)
	}
	msg := "hello tcp"
	tryMsg(t, conn, &msg, &builder, &writer)

	t.Log("Testing multiple Protocol Connections")
	conn2, err2 := net.Dial("tcp", "127.0.0.1:8080")
	if err2 != nil {
		t.Fatal(err)
	}
	msg2 := "hello tcp 2"
	tryMsg(t, conn2, &msg2, &builder, &writer)

	t.Log("Testing protocol Write()")
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		t.Fatal(err)
	}
	//launch a goroutine to accept the Protocol's outgoing connection
	//and stream the data to the syncWriter used by tryMsg
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Fatal(err)
		}
		buf := make([]byte, 100)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		writer.Write(buf[:n])
	}()
	err = tcp.Address("127.0.0.1:8081")
	if err != nil {
		t.Fatal(err)
	}
	msg3 := "hello tcp 3"
	tryMsg(t, tcp, &msg3, &builder, &writer)

	t.Log("Testing connection reuse")
	msg4 := "hello tcp 4"
	tryMsg(t, conn2, &msg4, &builder, &writer)

	tcp.Finish()
}

func tryMsg(t *testing.T, conn io.Writer, msg *string, b *strings.Builder, s *syncWriter) {
	_, err := fmt.Fprint(conn, *msg)
	if err != nil {
		t.Fatal(err)
	}

	s.Wait()
	str := b.String()
	if dif := strings.Compare(str, *msg); dif != 0 {
		t.Errorf("%s != %s", str, *msg)
	}
	b.Reset()
}

func (s *syncWriter) Write(b []byte) (int, error) {
	n, err := s.w.Write(b)
	s.done <- true
	return n, err
}

func (s *syncWriter) Wait() {
	<-s.done
}
