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
	SHUTDOWN                   = 6
	CLIENT_SHUTDOWN            = 7
	ADD_SERIES                 = 8
	REMOVE_SERIES              = 9
	LISTENER_ADD_SERIES        = 10
	LISTENER_REMOVE_SERIES     = 11
	LISTENER_NEW_SERIES_POINTS = 12
	UPDATE_SERIES              = 13
	RESPONSE_OK                = 14
	WORKER_JOB                 = 15
	WORKER_JOB_RESULT          = 16
	WORKER_JOB_CANCEL          = 21
	WORKER_JOB_CANCELLED       = 22
	WORKER_UPDATE_BUSY         = 23
	WORKER_REFUSED             = 24
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
