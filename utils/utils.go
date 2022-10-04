package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
)

func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

func Contains(num int, array []int) bool {
	for _, temp := range array {
		if num == temp {
			return true
		}
	}

	return false
}

func TimeFormat(timeStamp uint64) string {
	unixTime := time.Unix(int64(timeStamp), 0)

	return unixTime.Format("2006-01-02 15:04:05")
}
