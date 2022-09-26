package main

import (
	"log"
	"net"
	"strconv"

	qpack "github.com/cesbit/go-qpack"
)

func handlePkg(pkgCh chan *pkg) {
	for {
		data := <-pkgCh
		packageDataBuf := data.GetData()

		unpacked, _ := qpack.Unpack(packageDataBuf, 0)
		unboxed, ok := unpacked.(map[interface{}]interface{})

		if !ok {
			log.Println("ERROR: var not a map")
			continue
		}

		for key, element := range unboxed {
			name, okName := key.(string)

			if okName {
				// converting cstring to string
				nameBytes := []byte(name)
				name = string(nameBytes[:len(nameBytes)-1])
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

	pkgCh := make(chan *pkg)

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

	defer conn.Close()
	log.Printf("Starting listening on UDP port %s\n", udpPort)

	go handlePkg(pkgCh)

	buf := NewUdpBuffer()
	buf.SetConn(conn)
	buf.SetPkgChan(pkgCh)

	go buf.Read()
}
