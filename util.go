package main

import enodolib "github.com/SiriDB/siridb-enodo-go-lib"

func getEnodoId() {
	enodoId = enodolib.GenerateEnodoId()
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}
