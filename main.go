package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var hubConn net.Conn
var connErr error

var enodoId string
var seriesToWatch map[string]SeriesConfig
var groupsToWatch map[string]GroupConfig
var seriesCountUpdate map[string]int = make(map[string]int)
var updateLock sync.RWMutex = sync.RWMutex{}

var hubHost = os.Getenv("ENODO_HUB_HOSTNAME")
var hubPort = os.Getenv("ENODO_HUB_PORT")
var tcpPort = os.Getenv("ENODO_TCP_PORT")
var webserverPort = os.Getenv("ENODO_READY_PORT")
var internal_security_token = os.Getenv("ENODO_INTERNAL_SECURITY_TOKEN")

func main() {
	generateEnodoId()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		os.Exit(0)
	}(sigc)

	var wg sync.WaitGroup
	wg.Add(1)

	go httpReadyWebserver(webserverPort)
	go setupHubConn(hubHost, hubPort)
	go readFromTcp()

	wg.Wait()
}
