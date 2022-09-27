package main

import (
	"bytes"
)

type collect struct {
	data   [][]byte
	pkg_id int
	count  int
}

func NewCollect(seq_sz int, pkg_id int) *collect {
	return &collect{
		data:   make([][]byte, seq_sz),
		pkg_id: pkg_id,
		count:  0,
	}
}

func (c *collect) Size() int {
	return len(c.data)
}

func (c *collect) SetPart(seq_id int, part []byte) {
	dst := make([]byte, len(part))
	copy(dst, part)
	c.data[seq_id] = dst
	c.count += 1
}

func (c *collect) IsComplete() bool {
	return c.count == len(c.data)
}

func (c *collect) GetData() []byte {
	sep := make([]byte, 0)
	raw := bytes.Join(c.data, sep)
	return raw
}
