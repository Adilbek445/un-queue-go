package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

func IsFileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

func isDirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
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

func collectErrorBuf(errorMessage string) []byte {
	header := byte(0x02)
	errorMessageSize := int16(len(errorMessage))

	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(errorMessageSize))

	errorMessageByte := []byte(errorMessage)
	payload := []byte{header}
	payload = append(payload, b...)
	payload = append(payload, errorMessageByte...)
	return payload
}

func IsAlphanumericASCII(s string) bool {
	for _, r := range s {
		isLetterOrDigit := (r >= '0' && r <= '9') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z')

		if !isLetterOrDigit {
			return false
		}
	}

	return len(s) > 0
}
