package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		headerBytes := make([]byte, 7)
		_, err := reader.Read(headerBytes)
		header := serializeHeader([7]byte(headerBytes))
		queueNameBytes := make([]byte, header.QueueNameSize)
		reader.Read(queueNameBytes)

		fmt.Println(header.QueueNameSize)

		queueName := string(queueNameBytes)
		header.QueueName = queueName

		fmt.Println("QueueName " + header.QueueName)
		handleMessage(header, reader)
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}
		fmt.Printf("Received")
		// Можно отправить ответ клиенту
		conn.Write([]byte("ACK\n"))
	}
}

func main() {
	listenPort := ":8070"
	currentDirectory = getCurrentDirectory()
	listener, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", listenPort, err)
	}
	defer listener.Close()

	log.Printf("Server is listening on %s", listenPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go handleConnection(conn) // goroutine — аналог асинхронного ожидания epoll
	}
}
