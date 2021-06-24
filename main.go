package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

var hubConn net.Conn

var enodoId string
var seriesToWatch map[string]SeriesConfig
var seriesCountUpdate map[string]int = make(map[string]int)
var updateLock sync.RWMutex = sync.RWMutex{}

var hostname = os.Getenv("ENODO_HUB_HOSTNAME")
var port = os.Getenv("ENODO_HUB_PORT")
var pipe_path = os.Getenv("ENODO_PIPE_PATH")
var internal_security_token = os.Getenv("ENODO_INTERNAL_SECURITY_TOKEN")

func main() {
	getEnodoId()
	hubConn, _ = net.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port))
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)
	wg.Add(1)
	wg.Add(1)
	go watchIncommingData()
	go handshake()
	go heartbeat()
	go readFromPipe()
	go checkUpdates()
	wg.Wait()
}
