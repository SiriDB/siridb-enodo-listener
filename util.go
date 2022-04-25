package main

import (
	"strings"

	"github.com/google/uuid"
)

func generateEnodoId() {
	uuidWithHyphen := uuid.New()
	enodoId = strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}
