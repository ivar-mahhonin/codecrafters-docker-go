package server

import (
	"fmt"
	"io"
	"net"
	"os"
)

type TcpServer struct {
	host string
	port string
}

func NewTcpServer(host string, port string) *TcpServer {
	return &TcpServer{host, port}
}

func (server *TcpServer) Start() {
	if server.host == "" || server.port == "" {
		fmt.Println("Host and port must be set")
		stop()
	}

	address := fmt.Sprintf("%s:%s", server.host, server.port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Printf("Failed to bind to port %s \n", server.port)
		stop()
	}

	fmt.Printf("Server is listening on: %s \n", address)

	server.Listen(listener)
}

func (server *TcpServer) Stop() {
	stop()
}

func stop() {
	fmt.Println("Stopping server")
	os.Exit(1)
}

func (server *TcpServer) Listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			stop()
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
				println("EOF reached")
				break
			} else {
				fmt.Println("Failed reading from connection:", err)
				stop()
			}
		}

		message := string(buf)

		fmt.Println("Received read: ", message)

		response := makeResponseMessage(message)
		fmt.Println("Sending response:", response)

		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Failed write response:", err)
		} else {
			fmt.Println("Wrote response successfuly")
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
