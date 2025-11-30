package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

func getTimeNow() int64 {
	now := time.Now()
	return now.UnixMilli()
}

func serializeHeader(buf [7]byte) Header {
	header := Header{}

	header.HeaderFlag = int8(buf[0])
	header.QueueNameSize = int16(binary.BigEndian.Uint16(buf[1:3]))
	fmt.Println("QueueNameSize %u", header.QueueNameSize)
	header.TailerOrPayloadSize = int32(binary.BigEndian.Uint32(buf[3:7]))
	return header
}
