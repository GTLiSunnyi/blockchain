package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

// 序列化
func Serialize(something interface{}) []byte {
	data, err := json.Marshal(something)
	if err != nil {
		log.Println("序列化失败")
		log.Panic(err)
	}
	return data
}
