package main

import (
	"bufio"
	"fmt"
)

type Header struct {
	HeaderFlag          int8
	QueueNameSize       int16
	TailerOrPayloadSize int32
	QueueName           string
}

type MessageFlag byte

const (
	GET_MESSAGE       MessageFlag = 0x1
	WRITE_MESSAGE     MessageFlag = 0x2
	GET_STAT          MessageFlag = 0x3
	CHECK_NEW_MESSAGE MessageFlag = 0x4
)

func handleMessage(header Header, reader *bufio.Reader) {

	flag := header.HeaderFlag

	switch flag {
	case int8(GET_MESSAGE):
		fmt.Println("Get message")
	case int8(WRITE_MESSAGE):
		writeNewQueue(header, reader)
	case int8(GET_STAT):
		fmt.Println("Get stat")
	case int8(CHECK_NEW_MESSAGE):
		fmt.Println("Check new message")
	default:
		fmt.Println("Unknown flag")
	}
}
