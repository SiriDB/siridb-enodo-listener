package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var hubConn net.Conn
var connErr error

var enodoId string
var seriesToWatch map[string]SeriesConfig
var groupsToWatch map[string]GroupConfig
var seriesCountUpdate map[string]int = make(map[string]int)
var updateLock sync.RWMutex = sync.RWMutex{}

var hostname = os.Getenv("ENODO_HUB_HOSTNAME")
var port = os.Getenv("ENODO_HUB_PORT")
var pipe_path = os.Getenv("ENODO_PIPE_PATH")
var internal_security_token = os.Getenv("ENODO_INTERNAL_SECURITY_TOKEN")

func main() {
	getEnodoId()
	for {
		hubConn, connErr = net.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port))
		if connErr == nil {
			log.Println("Connection made to Hub")
			var wg sync.WaitGroup
			wg.Add(4)
			go watchIncommingData()
			go handshake()
			go heartbeat()
			go readFromPipe()
			go checkUpdates()
			wg.Wait()
		} else {
			log.Println(connErr)
		}
		timer := time.After(time.Second * 10)
		<-timer
	}
}
