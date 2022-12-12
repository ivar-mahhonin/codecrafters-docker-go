package main

import (
	"fmt"
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
		fmt.Println(fmt.Sprintf("Failed to bind to port %s", CONN_PORT))
		os.Exit(1)
	}

	fmt.Printf("Listening on: %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println("Failed reading from connection:", err)
			return
		}

		fmt.Println("Received read: ", string(buf))

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Failed write:", err)
			return
		}

		conn.Close()
	}
}
