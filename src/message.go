package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
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

var (
	queueLocks   = make(map[string]*sync.Mutex)
	queueLocksMu sync.RWMutex
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	headerBytes := make([]byte, 7)
	_, err := reader.Read(headerBytes)
	header := serializeHeader([7]byte(headerBytes))
	queueNameBytes := make([]byte, header.QueueNameSize)
	reader.Read(queueNameBytes)

	fmt.Println(header.QueueNameSize)

	if header.HeaderFlag != int8(WRITE_MESSAGE) {
		tailerNameBytes := make([]byte, header.TailerOrPayloadSize)
		if _, err := io.ReadFull(reader, tailerNameBytes); err != nil {
			log.Println("Connection closed (tailer read):", err)
			return
		}
		header.TailerName = string(tailerNameBytes)
	}

	queueName := string(queueNameBytes)
	header.QueueName = queueName

	fmt.Println("QueueName " + header.QueueName)
	fmt.Println("tailer " + header.TailerName)
	handleMessage(header, reader, writer)
	if err != nil {
		log.Println("Connection closed:", err)
		return
	}
	fmt.Printf("Received")

}

func handleMessage(header Header, reader *bufio.Reader, writer *bufio.Writer) {
	queue := header.QueueName

	mu := getQueueLock(queue)

	fmt.Println("TRY LOCK:", queue, "flag:", header.HeaderFlag)
	mu.Lock()
	fmt.Println("LOCKED:", queue)

	defer func() {
		mu.Unlock()
		fmt.Println("UNLOCK:", queue)
	}()

	switch header.HeaderFlag {
	case int8(GET_MESSAGE):
		getMessage(header, writer)
	case int8(WRITE_MESSAGE):
		writeMessage(header, reader)
	case int8(GET_STAT):
		getStat(header, writer)
	case int8(CHECK_NEW_MESSAGE):
		checkNewMessage(header, writer)
	default:
		fmt.Println("Unknown flag")
	}
}

func getQueueLock(queue string) *sync.Mutex {
	queueLocksMu.RLock()
	m, ok := queueLocks[queue]
	queueLocksMu.RUnlock()

	if ok {
		fmt.Println("LOCK: existing mutex for", queue)
		return m
	}

	queueLocksMu.Lock()
	defer queueLocksMu.Unlock()

	if m2, exists := queueLocks[queue]; exists {
		fmt.Println("LOCK: mutex created by someone else earlier for", queue)
		return m2
	}

	fmt.Println("LOCK: creating NEW mutex for", queue)
	m = &sync.Mutex{}
	queueLocks[queue] = m
	return m
}

func validationRequest() {}
