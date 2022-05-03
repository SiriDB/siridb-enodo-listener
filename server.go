package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

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

func readFromTcp() {
	const headerSize = 8

	var gds = func(data []byte) (int, error) {
		dataSize := int(binary.LittleEndian.Uint32(data[0:4]))
		return dataSize, nil
	}

	log.Printf("Starting listening on TCP port %s\n", tcpPort)
	pkgCh := make(chan *pkg)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", tcpPort))
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	go handlePkg(pkgCh)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept Error", err)
			continue
		}

		buf := NewBuffer()
		buf.SetConn(conn)
		buf.SetPkgChan(pkgCh)

		go buf.ReadToBuffer(headerSize, gds)

		log.Println("Accepted ", conn.RemoteAddr())
	}
}
