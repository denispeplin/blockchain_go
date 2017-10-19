package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func RaiseError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	RaiseError(err)

	return buff.Bytes()
}
