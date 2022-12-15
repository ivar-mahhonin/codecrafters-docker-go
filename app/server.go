package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "6379"
	CONN_TYPE = "tcp"
)

func main() {
	address := fmt.Sprintf("%s:%s", CONN_HOST, CONN_PORT)

	listener, err := net.Listen(CONN_TYPE, address)

	if err != nil {
		fmt.Printf("Failed to bind to port %s \n", CONN_PORT)
		os.Exit(1)
	}

	fmt.Printf("Listening on: %s \n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)

		if _, err := conn.Read(buf); err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Failed reading from connection:", err)
				os.Exit(1)
			}
		}

		message := string(buf)

		fmt.Println("Received read: ", message)

		response := makeResponseMessage(message)
		fmt.Println("Sending response:", response)

		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Failed write:", err)
		}
	}
}

func makeResponseMessage(message string) string {
	response := ""
	switch message {
	default:
		response = stringToReply("PONG")
	}
	return response
}

func stringToReply(message string) string {
	return fmt.Sprintf("+%s\r\n", message)
}
