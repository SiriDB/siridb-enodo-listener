package main

import (
	"encoding/binary"
	"log"
	"regexp"
	"time"

	qpack "github.com/transceptor-technology/go-qpack"
)

func sendSeriesUpdate(seriesAndCounts map[string]int) {
	bdata, err := qpack.Pack(seriesAndCounts)
	log.Println("SENDING UPDATE")
	if err == nil {
		pkg := CreatePackage(1, LISTENER_NEW_SERIES_POINTS, bdata)
		hubConn.Write(pkg)
	}
}

func sendNewSeriesFoundForGroup(seriesName string, groupName string) {
	data := make(map[string]string)
	data["group"] = groupName
	data["series_name"] = seriesName
	bdata, err := qpack.Pack(data)
	log.Println("SENDING FOUND SERIES FOR GROUP")
	if err == nil {
		pkg := CreatePackage(1, LISTENER_ADD_SERIES, bdata)
		hubConn.Write(pkg)
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

func heartbeat() {
	for {
		timer := time.After(time.Second * 25)
		<-timer
		data, err := qpack.Pack(enodoId)
		if err == nil {
			pkg := CreatePackage(1, HEARTBEAT, data)
			hubConn.Write(pkg)
			log.Println("Send heartbeat to hub")
		}
	}
}

func handshake() {
	data := map[string]interface{}{"client_id": enodoId, "client_type": "listener", "token": internal_security_token, "version": "0.0.1"}
	bdata, err := qpack.Pack(data)

	if err == nil {
		pkg := CreatePackage(1, HANDSHAKE, bdata)
		hubConn.Write(pkg)
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
