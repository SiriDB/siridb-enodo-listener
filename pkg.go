package main

type pkg struct {
	header []byte
	data   []byte
}

func (p *pkg) GetHeader() []byte {
	return p.header
}

func (p *pkg) GetData() []byte {
	return p.data
}
