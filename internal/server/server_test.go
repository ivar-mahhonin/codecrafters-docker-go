package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "6379"
)

func init() {
	//server := NewTcpServer(CONN_HOST, CONN_PORT)
	//go server.Start()
}

func TestTCPServerRunning(t *testing.T) {
	conn := createDialClient()
	defer conn.Close()
}
func TestOnePongResponse(t *testing.T) {
	conn := createDialClient()
	defer conn.Close()

	message := "PING\n"
	response := writeAndReadMessage(conn, message, t)
	if response != "+PONG" {
		t.Fatalf(`Response command should be "+PONG\r\n". Intead it was: %s`, response)
	}
}

func TestTwoPongResponse(t *testing.T) {
	conn := createDialClient()
	defer conn.Close()

	messages := []string{"Test Request1\n", "Test Request2\n"}
	for _, message := range messages {
		response := writeAndReadMessage(conn, message, t)
		if response != "+PONG" {
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
		if response != "+PONG" {
			t.Fatalf(`Response command should be "+PONG\r\n". Intead it was: %s`, response)
		} else {
			log.Printf("Received response: %s", response)
			log.Printf("Waiting...")
		}
	}
}

func TestEchoResponse(t *testing.T) {
	conn := createDialClient()
	defer conn.Close()

	message := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	response := writeAndReadMessage(conn, message, t)

	if response != "+hey" {
		t.Fatalf(`Response command should be "+hey". Intead it was: %s`, response)
	}
}

func TestSetCommand(t *testing.T) {
	conn := createDialClient()
	defer conn.Close()

	message := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	response := writeAndReadMessage(conn, message, t)

	if response != "+OK" {
		t.Fatalf(`Response command should be "+OK". Intead it was: %s`, response)
	}
}

func TestGetCommand(t *testing.T) {
	conn := createDialClient()
	defer conn.Close()

	message := "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"
	response := writeAndReadMessage(conn, message, t)

	if response != "+value" {
		t.Fatalf(`Response command should be "+OK". Intead it was: %s`, response)
	}
}

func TestSetCommandWithGettingExpiredValue(t *testing.T) {
	conn := createDialClient()
	defer conn.Close()

	message := fmt.Sprintf("*5\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n$2\r\nPX\r\n$3\r\n%s\r\n", "100")
	response := writeAndReadMessage(conn, message, t)

	if response != "+OK" {
		t.Fatalf(`Response command should be "+OK". Intead it was: %s`, response)
	}

	time.Sleep(101 * time.Millisecond)

	message = "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"
	response = writeAndReadMessage(conn, message, t)

	if response != "$-1" {
		t.Fatalf(`Response command should be "$-1". Intead it was: %s`, response)
	}
}

func TestSetCommandWithGettingNonExpiredValue(t *testing.T) {
	conn := createDialClient()
	defer conn.Close()

	message := fmt.Sprintf("*5\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n$2\r\nPX\r\n$5\r\n%s\r\n", "10000")
	response := writeAndReadMessage(conn, message, t)

	if response != "+OK" {
		t.Fatalf(`Response command should be "+OK". Intead it was: %s`, response)
	}

	time.Sleep(99 * time.Millisecond)

	message = "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"
	response = writeAndReadMessage(conn, message, t)

	if response != "+value" {
		t.Fatalf(`Response command should be "+value". Intead it was: %s`, response)
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
	response := string(bytes.Trim(received, "\x00"))
	return strings.TrimRight(response, "\r\n")
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
