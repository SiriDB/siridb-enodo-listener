package main

import (
	"log"
	"net"
	"strconv"

	qpack "github.com/cesbit/go-qpack"
)

func handleData(dataCh chan []byte) {
	for {
		data := <-dataCh

		unpacked, _ := qpack.Unpack(data, 0)
		unboxed, ok := unpacked.(map[interface{}]interface{})

		if !ok {
			log.Println("ERROR: var not a map")
			continue
		}

		for key, element := range unboxed {
			name, okName := key.(string)

			if okName {
				name = name[:len(name)-1]
				if series, ok := seriesToWatch[name]; ok {
					pointsList, okPointsList := element.([]interface{})
					if okPointsList {
						if series.IsRealtime {
							singleUpdate := make(map[string]int)
							singleUpdate[name] = len(pointsList)
							sendSeriesUpdate(singleUpdate)
						} else {
							updateLock.Lock()
							seriesCountUpdate[name] += len(pointsList)
							updateLock.Unlock()
						}
					}
				} else {
					for _, groupConfig := range groupsToWatch {
						if found := groupConfig.Regex.MatchString(name); found {
							sendNewSeriesFoundForGroup(name, groupConfig.Name)
						}
					}
				}
			}
		}
	}
}

func readFromUdp() {
	const headerSize = 8

	dataCh := make(chan []byte)

	port, err := strconv.Atoi(udpPort)

	if err != nil {
		log.Fatal("Incorrect port")
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	log.Printf("Starting listening on UDP port %s\n", udpPort)

	go handleData(dataCh)

	buf := NewUdpBuffer()
	buf.SetConn(conn)
	buf.SetDataChan(dataCh)

	go buf.Read()
}
