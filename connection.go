package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"

	qpack "github.com/cesbit/go-qpack"
)

var connectedToHub = false

func sendSeriesUpdate(seriesAndCounts map[string]int) {
	bdata, err := qpack.Pack(seriesAndCounts)
	if err == nil {
		pkg := CreatePackage(1, LISTENER_NEW_SERIES_POINTS, bdata)
		if _, err = hubConn.Write(pkg); err != nil {
			log.Println("Failed to write 'series update' pacakge")
		} else {
			log.Println("Send 'series update' pacakge")
		}
	} else {
		log.Println("Failed to pack 'series update' data")
	}
}

func sendNewSeriesFoundForGroup(seriesName string, groupName string) {
	data := make(map[string]string)
	data["group"] = groupName
	data["series_name"] = seriesName
	bdata, err := qpack.Pack(data)
	if err == nil {
		pkg := CreatePackage(1, LISTENER_ADD_SERIES, bdata)
		if _, err = hubConn.Write(pkg); err != nil {
			log.Println("Failed to write 'found series for group' pacakge")
		} else {
			log.Println("Send 'found series for group' pacakge")
		}
	} else {
		log.Println("Failed to pack 'found series for group' data")
	}
}

func checkUpdates() {
	for {
		timer := time.After(time.Second * 60)
		<-timer
		updateLock.Lock()
		if len(seriesCountUpdate) > 0 {
			sendSeriesUpdate(seriesCountUpdate)
			seriesCountUpdate = make(map[string]int)
		}
		updateLock.Unlock()
	}
}

func heartbeat(wg *sync.WaitGroup) {
	for {
		timer := time.After(time.Second * 25)
		<-timer
		data, err := qpack.Pack(enodoId)
		if err == nil {
			pkg := CreatePackage(1, HEARTBEAT, data)
			if _, err = hubConn.Write(pkg); err != nil {
				log.Println("Failed to write 'heartbeat' package")
				wg.Done()
				return
			} else {
				log.Println("Send heartbeat to hub")
			}
		} else {
			log.Println("Failed to pack 'heartbeat' data")
		}
	}
}

func handshake() {
	data := map[string]interface{}{
		"client_id":   enodoId,
		"client_type": "listener",
		"token":       internal_security_token,
		"version":     Version,
	}

	bdata, err := qpack.Pack(data)

	if err == nil {
		pkg := CreatePackage(1, HANDSHAKE, bdata)
		if _, err = hubConn.Write(pkg); err != nil {
			log.Println("Failed to write 'handshake' pacakge")
		} else {
			log.Println("Send handshake package")
		}
	} else {
		log.Println("Failed to pack 'package' data")
	}
}

func httpReadyWebserver(webserverPort string) {
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if connectedToHub {
			w.WriteHeader(200)
			_, err := w.Write([]byte("ok"))
			if err != nil {
				log.Println("Could not write response to webserver client")
			}
		} else {
			w.WriteHeader(500)
			_, err := w.Write([]byte("Not connected to Hub"))
			if err != nil {
				log.Println("Could not write response to webserver client")
			}
		}
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", webserverPort), nil))
}

func setupHubConn(hubHost string, hubPort string) {
	for {
		hubConn, connErr = net.Dial("tcp", fmt.Sprintf("%s:%s", hubHost, hubPort))
		if connErr == nil {
			log.Println("Connection made to Hub")
			connectedToHub = true
			var wg sync.WaitGroup
			wg.Add(1)
			go watchIncommingData()
			go handshake()
			go heartbeat(&wg)
			go checkUpdates()
			wg.Wait()
		} else {
			log.Println(connErr)
		}
		if hubConn != nil {
			hubConn.Close()
		}
		connectedToHub = false
		timer := time.After(time.Second * 10)
		<-timer
		log.Println("Trying to reconnect")
	}
}

func watchIncommingData() {
	var gds = func(data []byte) (int, error) {
		return int(binary.BigEndian.Uint32(data[0:4])), nil
	}
	pkgCh := make(chan *pkg)
	dataBuf := NewBuffer()
	dataBuf.SetConn(hubConn)
	dataBuf.SetPkgChan(pkgCh)

	go dataBuf.ReadToBuffer(PACKET_HEADER_LEN, gds)
	for {
		data := <-pkgCh
		packageDataBuf := data.GetData()
		_, packageType, _ := ReadHeaderFromBinaryData(data.GetHeader())

		switch packageType {
		case HANDSHAKE_OK:
			log.Println("Hands shaked with Hub")
		case HANDSHAKE_FAIL:
			log.Println("Hub does not want to shake hands")
		case HEARTBEAT:
			log.Println("Heartbeat back from Hub")
		case RESPONSE_OK:
			log.Println("Hub received update correctly")
		case UNKNOWN_CLIENT:
			log.Println("Hub does not recognize us")
		case UPDATE_SERIES:
			log.Println("Received new list of series to watch")

			newSeriesToWatch, err := qpack.Unpack(packageDataBuf, 0)

			if err == nil {
				listSlice, ok := newSeriesToWatch.([]interface{})
				if !ok {
					log.Println("ERROR: var not a slice/list")
					continue
				}
				seriesToWatch = make(map[string]SeriesConfig)
				groupsToWatch = make(map[string]GroupConfig)
				for _, s := range listSlice {
					unboxed, ok := s.(map[interface{}]interface{})
					if !ok {
						log.Println("ERROR: var not a map")
						continue
					}
					name, okName := unboxed["name"].(string)
					isRealtime, okIsRealtime := unboxed["realtime"].(bool)
					isGroup, okIsGroup := unboxed["isGroup"].(bool)
					if okName && okIsRealtime {
						if okIsGroup && isGroup {
							selector, okSelector := unboxed["selector"].(string)
							if okSelector {
								regex, err := regexp.Compile(selector)
								if err == nil {
									gc := GroupConfig{name, regex}
									groupsToWatch[name] = gc
								}
							}
						} else {
							sc := SeriesConfig{name, isRealtime}
							seriesToWatch[name] = sc
						}
					}
				}
				log.Println("SERIES TO WATCH: ", seriesToWatch)
				log.Println("GROUPS TO WATCH: ", groupsToWatch)
			}
		}
	}
}
