package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Header struct {
	HeaderFlag          int8
	QueueNameSize       int16
	TailerOrPayloadSize int32
	QueueName           string
	TailerName          string
}

type MessageFlag byte

const (
	GET_MESSAGE       MessageFlag = 0x1
	WRITE_MESSAGE     MessageFlag = 0x2
	GET_STAT          MessageFlag = 0x3
	CHECK_NEW_MESSAGE MessageFlag = 0x4
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		headerBytes := make([]byte, 7)
		_, err := reader.Read(headerBytes)
		header := serializeHeader([7]byte(headerBytes))
		queueNameBytes := make([]byte, header.QueueNameSize)
		reader.Read(queueNameBytes)

		fmt.Println(header.QueueNameSize)

		if header.HeaderFlag != int8(WRITE_MESSAGE) {
			tailerNameBytes := make([]byte, header.TailerOrPayloadSize)
			reader.Read(tailerNameBytes)
			header.TailerName = string(tailerNameBytes)
		}

		queueName := string(queueNameBytes)
		header.QueueName = queueName

		fmt.Println("QueueName " + header.QueueName)
		handleMessage(header, reader, writer)
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}
		fmt.Printf("Received")

	}
}

func handleMessage(header Header, reader *bufio.Reader, writer *bufio.Writer) {

	flag := header.HeaderFlag

	switch flag {
	case int8(GET_MESSAGE):
		fmt.Println("Get message")
	case int8(WRITE_MESSAGE):
		writeMessage(header, reader)
	case int8(GET_STAT):
		getMessage(header, writer)
	case int8(CHECK_NEW_MESSAGE):
		fmt.Println("Check new message")
	default:
		fmt.Println("Unknown flag")
	}
}
