package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"testing"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "6379"
	CONN_TYPE = "tcp"
)

func init() {
	log.SetOutput(ioutil.Discard) // Disable, if you want to see debug logs
	server := NewTcpServer(CONN_HOST, CONN_PORT)
	go server.Start()
}

func TestTCPServerRunning(t *testing.T) {
	conn := createDialClient()
	defer conn.Close()
}
func TestOnePongResponse(t *testing.T) {
	conn := createDialClient()
	defer conn.Close()

	message := "Test Request\n"
	response := writeAndReadMessage(conn, message, t)
	if strings.TrimRight(response, "\r\n") == "+PONG" {
		t.Fatalf(`Response command should be "+PONG\r\n". Intead it was: %s`, response)
	}
}

func TestTwoPongResponse(t *testing.T) {
	conn := createDialClient()
	defer conn.Close()

	messages := []string{"Test Request1\n", "Test Request2\n"}
	for _, message := range messages {
		response := writeAndReadMessage(conn, message, t)
		if strings.TrimRight(response, "\r\n") == "+PONG" {
			t.Fatalf(`Response command should be "+PONG\r\n". Intead it was: %s`, response)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	connections := []net.Conn{
		createDialClient(), createDialClient(),
		createDialClient(), createDialClient(),
	}
	for i, conn := range connections {
		defer conn.Close()

		response := writeAndReadMessage(conn, fmt.Sprintf("Request_%d\n", i), t)
		if strings.TrimRight(response, "\r\n") == "+PONG" {
			t.Fatalf(`Response command should be "+PONG\r\n". Intead it was: %s`, response)
		} else {
			log.Printf("Received response: %s", response)
			log.Printf("Waiting...")
		}
	}
}

func writeAndReadMessage(conn net.Conn, message string, t *testing.T) string {
	_, write_err := conn.Write([]byte(message))
	if write_err != nil {
		t.Fatalf(`Sending command to redis-go failed %v`, write_err)
		return ""
	}
	received := make([]byte, 1024)
	_, read_err := conn.Read(received)
	if read_err != nil {
		t.Fatalf(`Reading response from redis-go failed %v`, read_err)
		return ""
	}
	response := string(received)
	return response
}

func createDialClient() net.Conn {
	address := fmt.Sprintf("%s:%s", CONN_HOST, CONN_PORT)
	conn, err := net.Dial("tcp", address)

	if err != nil {
		log.Printf("Error in createDialClient: %s", err)
		os.Exit(1)
	}
	return conn
}
