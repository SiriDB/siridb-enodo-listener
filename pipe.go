package main

import (
	"encoding/binary"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	enodolib "github.com/SiriDB/siridb-enodo-go-lib"
	qpack "github.com/transceptor-technology/go-qpack"
)

func readFromPipe() {
	os.Remove(pipe_path)
	log.Println("Starting echo server")
	ln, err := net.Listen("unix", pipe_path)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		ln.Close()
		os.Exit(0)
	}(ln, sigc)

	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		var gds = func(data []byte) (uint32, error) {
			dataSize := uint32(binary.LittleEndian.Uint32(data[0:4]))
			return dataSize, nil
		}

		dataBuf := enodolib.NewBuffer()
		dataBuf.SetConn(fd)
		pkgCh := dataBuf.GetPkgChan()
		go dataBuf.ReadToBuffer(8, gds)

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
								log.Println("COUNT: ", len(pointsList))
								updateLock.Lock()
								if _, ok := seriesCountUpdate[name]; ok {
									seriesCountUpdate[name] += len(pointsList)
								} else {
									seriesCountUpdate[name] = len(pointsList)
								}
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
}
