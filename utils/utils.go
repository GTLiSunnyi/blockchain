package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

// uint64转[]byte
func UintToByte(num uint64) []byte {
	buffer := bytes.NewBuffer([]byte{})
	err := binary.Write(buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
