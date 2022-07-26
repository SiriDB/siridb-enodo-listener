package main

import (
	"log"
	"net"
)

// buffer is used to read data from a connection.
type buffer struct {
	conn     net.Conn
	data     []byte
	dataSize int
	len      int
	pkg      *pkg
	pkgCh    chan *pkg
}

// newBuffer retur a pointer to a new buffer.
func NewBuffer() *buffer {
	return &buffer{
		conn:     nil,
		data:     make([]byte, 0),
		dataSize: 0,
		len:      0,
		pkg:      nil,
		pkgCh:    make(chan *pkg),
	}
}

func (buf *buffer) SetConn(c net.Conn) {
	buf.conn = c
}

func (buf *buffer) SetPkgChan(pkgCh chan *pkg) {
	buf.pkgCh = pkgCh
}

func (buf buffer) ReadToBuffer(headerSize int, getDataSize func([]byte) (int, error)) {
	defer buf.conn.Close()

	wbuf := make([]byte, 8192)
	for {
		n, err := buf.conn.Read(wbuf)
		if err != nil {
			log.Println("Closing connection ", buf.conn.RemoteAddr())
			connectedToHub = false
			return
		}
		buf.len += n
		buf.data = append(buf.data, wbuf[:n]...)
		for buf.len >= headerSize {
			if buf.pkg == nil {
				buf.dataSize, err = getDataSize(buf.data)
				if err != nil {
					log.Println("Failed reading data size from ", buf.conn.RemoteAddr())
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
