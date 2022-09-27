package main

import (
	"encoding/binary"
	"net"
)

const MAX_PER_SERVER = 100

// buffer is used to read data from a connection.
type udpBuffer struct {
	conn      *net.UDPConn
	collector map[int]*collect
	pkgCh     chan *pkg
}

// newBuffer retur a pointer to a new buffer.
func NewUdpBuffer() *udpBuffer {
	return &udpBuffer{
		conn:      nil,
		collector: make(map[int]*collect),
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
		n, _, err := buf.conn.ReadFromUDP(rbuf)
		if err != nil {
			continue
		}

		id := int(binary.LittleEndian.Uint16(rbuf[0:2]))
		pkg_id := int(binary.LittleEndian.Uint16(rbuf[2:4]))
		seq_sz := int(binary.LittleEndian.Uint16(rbuf[4:6]))
		seq_id := int(binary.LittleEndian.Uint16(rbuf[6:8]))
		key := (id << 16) + (pkg_id % MAX_PER_SERVER)

		collector, exists := buf.collector[key]
		if !exists || collector.pkg_id != pkg_id || collector.Size() != seq_sz {
			collector = NewCollect(seq_sz, pkg_id)
			buf.collector[key] = collector
		}

		collector.SetPart(seq_id, rbuf[8:n])
		if collector.IsComplete() {
			data := collector.GetData()
			pkg := &pkg{data[0:8], data[8:]}
			buf.pkgCh <- pkg
			delete(buf.collector, key)
		}
	}
}
