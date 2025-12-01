package main

import (
	"log"
	"net"
)

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
