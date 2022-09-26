package main

import (
	"encoding/binary"
	"log"
	"net"
)

// buffer is used to read data from a connection.
type udpBuffer struct {
	conn      *net.UDPConn
	collector map[int]collect
	pkgCh     chan *pkg
}

// newBuffer retur a pointer to a new buffer.
func NewUdpBuffer() *udpBuffer {
	return &udpBuffer{
		conn:      nil,
		collector: make(map[int][][]byte),
		pkgCh:     make(chan *pkg),
	}
}

func (buf *udpBuffer) SetConn(c *net.UDPConn) {
	buf.conn = c
}

func (buf *udpBuffer) SetPkgChan(pkgCh chan *pkg) {
	buf.pkgCh = pkgCh
}

func (buf udpBuffer) Read() {
	defer buf.conn.Close()

	rbuf := make([]byte, 508)
	for {
		n, remote, err := buf.conn.ReadFromUDP(rbuf)
		if err != nil {
			continue
		}

		id := int(binary.LittleEndian.Uint16(rbuf[0:2]))
		pkg_id := int(binary.LittleEndian.Uint16(rbuf[2:4]))
		seq_sz := int(binary.LittleEndian.Uint16(rbuf[4:6]))
		seq_id := int(binary.LittleEndian.Uint16(rbuf[6:8]))
		key := (id << 16) + (pkg_id % 100)

		collector, exists := buf.collector[key]
		if (exists) {
			collector.
		}

		buf.len += n
		buf.data = append(buf.data, wbuf[:n]...)
		for buf.len >= 8 {

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
