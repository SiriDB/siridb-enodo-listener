package main

type collect struct {
	data    [][]byte
	seq_sz  int
	data_sz int
}

func (p *pkg) NewCollect(seq_sz) *collect {
	return &collect{
		data:      nil,
		data:      make([]byte, 0),
		dataSize:  0,
		len:       0,
		collector: make(map[int][][]byte),
		pkgCh:     make(chan *pkg),
	}
}

func (p *pkg) GetData() []byte {
	return p.data
}
