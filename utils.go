package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
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

func Serialize(input interface{}) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(input)
	RaiseError(err)

	return result.Bytes()
}
