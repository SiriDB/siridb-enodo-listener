package main

import (
	"log"
	"net"
)

// buffer is used to read data from a connection.
type udpBuffer struct {
	conn     net.UDPConn
	data     []byte
	dataSize int
	len      int
	pkg      *pkg
	pkgCh    chan *pkg
}

// newBuffer retur a pointer to a new buffer.
func NewUdpBuffer(c net.UDPConn) *udpBuffer {
	return &udpBuffer{
		conn:     c,
		data:     make([]byte, 0),
		dataSize: 0,
		len:      0,
		pkg:      nil,
		pkgCh:    make(chan *pkg),
	}
}

func (buf *udpBuffer) SetConn(c net.UDPConn) {
	buf.conn = c
}

func (buf *udpBuffer) SetPkgChan(pkgCh chan *pkg) {
	buf.pkgCh = pkgCh
}

func (buf udpBuffer) ReadToBuffer(headerSize int, getDataSize func([]byte) (int, error)) {
	defer buf.conn.Close()

	wbuf := make([]byte, 8192)
	for {
		n, remote, err := buf.conn.ReadFromUDP(wbuf)
		if err != nil {
			log.Println("Closing connection ", remote)
			return
		}
		buf.len += n
		buf.data = append(buf.data, wbuf[:n]...)
		for buf.len >= headerSize {
			if buf.pkg == nil {
				buf.dataSize, err = getDataSize(buf.data)
				if err != nil {
					log.Println("Failed reading data size from ", remote)
					return
				}
				buf.pkg = &pkg{make([]byte, headerSize), make([]byte, buf.dataSize)}
			}
			total := buf.dataSize + headerSize

			if buf.len < total {
				break
			}

			buf.pkg.header = buf.data[0:headerSize]
			buf.pkg.data = buf.data[headerSize:total]

			buf.pkgCh <- buf.pkg

			buf.data = buf.data[total:]
			buf.len -= total
			buf.pkg = nil
		}
	}
}
