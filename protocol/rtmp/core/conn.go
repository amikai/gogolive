package core

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type HandShakeSchema int

const (
	SimpleHS = iota
	ComplexHSSchema0
	ComplexHSSchema1
)

type Conn struct {
	net.Conn
	rw *bufio.ReadWriter
}

func NewConn(c net.Conn, bufSize int) *Conn {
	return &Conn{
		Conn: c,
		rw:   bufio.NewReadWriter(bufio.NewReaderSize(c, bufSize), bufio.NewWriterSize(c, bufSize)),
	}
}

func (c *Conn) SrvHandshake() (err error) {
	timeout := 5 * time.Second
	var random [(1 + 1536*2) * 2]byte

	C0C1C2 := random[:1536*2+1]
	C0 := C0C1C2[0:1]
	C1 := C0C1C2[1 : 1536+1]
	C0C1 := C0C1C2[:1536+1]
	C2 := C0C1C2[1536+1:]

	S0S1S2 := random[1536*2+1:]
	S0 := S0S1S2[:1]
	S1 := S0S1S2[1 : 1536+1]
	// S0S1 := S0S1S2[:1536+1]
	S2 := S0S1S2[1536+1:]

	// send C0C1
	c.Conn.SetDeadline(time.Now().Add(timeout))
	if _, err = io.ReadFull(c.rw, C0C1); err != nil {
		return

	}
	fmt.Println()

	// send S0S1S2
	c.Conn.SetDeadline(time.Now().Add(timeout))
	if C0[0] != 3 {
		err = fmt.Errorf("rtmp: handshake version=%d invalid", C0[0])
		return
	}
	S0[0] = 3
	copy(S1, C2)
	copy(S2, C1)

	// c.Conn.SetDeadline(time.Now().Add(timeout))
	if _, err = c.rw.Write(S0S1S2); err != nil {
		return
	}

	c.Conn.SetDeadline(time.Now().Add(timeout))
	if err = c.rw.Flush(); err != nil {
		return
	}

	// send C2
	c.Conn.SetDeadline(time.Now().Add(timeout))
	if _, err = io.ReadFull(c.rw, C2); err != nil {
		return
	}

	return
}
