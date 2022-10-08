package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

func IsContains(num int, array []int) bool {
	for _, temp := range array {
		if num == temp {
			return true
		}
	}

	return false
}
