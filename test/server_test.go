package test

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "6379"
	CONN_TYPE = "tcp"
)

func TestTCPServerRunning(t *testing.T) {
	conn, err := createDialClient()
	if err != nil {
		t.Fatalf(`Failed to connect to server %v`, err)
	}
	defer conn.Close()
}
func TestOnePongResponse(t *testing.T) {
	conn, _ := createDialClient()
	defer conn.Close()

	message := "Test Request\n"
	_, write_err := conn.Write([]byte(message))
	if write_err != nil {
		t.Fatalf(`Sending command to redis-go failed %v`, write_err)
	}
	received := make([]byte, 1024)
	_, read_err := conn.Read(received)
	if read_err != nil {
		t.Fatalf(`Reading response from redis-go failed %v`, read_err)
	}

	response := string(received)
	if strings.TrimRight(response, "\r\n") == "+PONG" {
		t.Fatalf(`Response command should be "+PONG\r\n". Intead it was: %s`, response)
	}
}

func TestTwoPongResponse(t *testing.T) {
	conn, _ := createDialClient()
	defer conn.Close()

	messages := make([]string, 2)
	messages[0] = "Test Request1\n"
	messages[1] = "Test Request2\n"
	for _, message := range messages {
		_, write_err := conn.Write([]byte(message))
		if write_err != nil {
			t.Fatalf(`Sending command to redis-go failed %v`, write_err)
		}
		received := make([]byte, 1024)
		_, read_err := conn.Read(received)
		if read_err != nil {
			t.Fatalf(`Reading response from redis-go failed %v`, read_err)
		}

		response := string(received)
		if strings.TrimRight(response, "\r\n") == "+PONG" {
			t.Fatalf(`Response command should be "+PONG\r\n". Intead it was: %s`, response)
		}
	}
}

func createDialClient() (net.Conn, error) {
	address := fmt.Sprintf("%s:%s", CONN_HOST, CONN_PORT)
	conn, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Printf("Error in createDialClient: %s", err)
	}
	return conn, err
}
