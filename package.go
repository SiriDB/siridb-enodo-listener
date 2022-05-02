package main

import (
	"encoding/binary"
)

const (
	HANDSHAKE                  = 1
	HANDSHAKE_OK               = 2
	HANDSHAKE_FAIL             = 3
	UNKNOWN_CLIENT             = 4
	HEARTBEAT                  = 5
	_                          = 6 // SHUTDOWN
	_                          = 7 // CLIENT_SHUTDOWN
	_                          = 8 // ADD_SERIES
	_                          = 9 // REMOVE_SERIES
	LISTENER_ADD_SERIES        = 10
	_                          = 11 // LISTENER_REMOVE_SERIES
	LISTENER_NEW_SERIES_POINTS = 12
	UPDATE_SERIES              = 13
	RESPONSE_OK                = 14
	_                          = 15 // WORKER_JOB
	_                          = 16 // WORKER_JOB_RESULT
	_                          = 21 // WORKER_JOB_CANCEL
	_                          = 22 // WORKER_JOB_CANCELLED
	_                          = 23 // WORKER_UPDATE_BUSY
	_                          = 24 // WORKER_REFUSED
)

const PACKET_HEADER_LEN = 6 // bytes

func CreatePackage(packageID int, packageType int, data []byte) []byte {
	pkg := make([]byte, PACKET_HEADER_LEN+len(data))

	binary.BigEndian.PutUint32(pkg[0:], uint32(len(data)))
	pkg[4] = uint8(packageType)
	pkg[5] = uint8(packageID)
	copy(pkg[6:], data)

	return pkg
}

func ReadHeaderFromBinaryData(data []byte) (int, int, int) {
	return int(binary.BigEndian.Uint32(data[0:4])), int(uint8(data[4])), int(uint8(data[5]))
}
