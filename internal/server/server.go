package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	redisStore "github.com/ivar-mahhonin/redis-go/internal/store"
)

const (
	CMD_PING = "ping"
	CMD_ECHO = "echo"
	CMD_GET  = "get"
	CMD_SET  = "set"
	CMD_PX   = "px"
)

const (
	BUFFER_SIZE = 1024
)

var store = redisStore.GetStore()

type TcpServer struct {
	host string
	port string
}

func NewTcpServer(host string, port string) *TcpServer {
	return &TcpServer{host, port}
}

func (server *TcpServer) Start() {
	if server.host == "" || server.port == "" {
		log.Println("Host and port must be set")
		stop()
	}
	address := fmt.Sprintf("%s:%s", server.host, server.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("Failed to bind to port %s \n", server.port)
		stop()
	}
	log.Printf("Server is listening on: %s \n", address)
	server.Listen(listener)
}

func (server *TcpServer) Stop() {
	stop()
}

func stop() {
	log.Println("Stopping server")
	os.Exit(1)
}

func (server *TcpServer) Listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			stop()
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, BUFFER_SIZE)

		if _, err := conn.Read(buf); err != nil {
			if err == io.EOF {
				log.Println("EOF reached")
				break
			} else {
				fmt.Println("Failed reading from connection:", err)
				stop()
			}
		}

		message := string(buf)

		log.Println("Received read:\n", message)

		response := parseAndExecuteCommand(message)

		log.Println("Sending response:", response)

		_, err := conn.Write([]byte(response))
		if err != nil {
			log.Println("Failed write response:", err)
		} else {
			log.Println("Wrote response successfuly")
		}
	}
}

func parseAndExecuteCommand(message string) string {
	var command []string
	response := ""
	if isReplArray(message) {
		cmd, err := extractArgumentsFromReplArray(message)
		if err != nil {
			log.Println("Failed to extract arguments from repl array")
			response = makeErrorResponse("Failed to extract arguments from repl array")
		} else {
			command = cmd
		}
	} else {
		command = []string{strings.Replace(message, "\n", "", -1)}
	}

	if command != nil && response == "" {
		response = executeCommand(command)
	}
	return response
}

func executeCommand(command []string) string {
	response := ""
	switch strings.ToLower(command[0]) {
	default:
		response = stringToReply("PONG")
	case CMD_ECHO:
		response = stringToReply(command[1])
	case CMD_SET:
		key := command[1]
		value := command[2]
		if isCommandWithExpiration(command) {
			exp, err := makeExpirationFromString(command[4])
			if err != nil {
				response = makeErrorResponse(err.Error())
				break
			}
			response = setStoreValueWithExpiration(key, value, exp)
		} else {
			response = setStoreValueWithoutExpiration(key, value)
		}
	case CMD_GET:
		key := command[1]
		value, found := getStoreValue(key)
		if !found {
			response = makeNullValueResponse()
		} else {
			response = stringToReply(value)
		}
	}
	return response
}

func isCommandWithExpiration(command []string) bool {
	return len(command) >= 5 && strings.ToLower(command[3]) == CMD_PX && command[4] != ""
}

func makeExpirationFromString(expiration string) (time.Time, error) {
	milliseconds, err := strconv.ParseInt(expiration, 10, 64)

	if err != nil {
		return time.Time{}, errors.New(fmt.Sprintf("Could not parse expiration time. PX is not integer: %s", expiration))
	}

	return time.Now().Add(time.Duration(milliseconds) * time.Millisecond), nil
}

func setStoreValueWithoutExpiration(key string, value string) string {
	var expiration time.Time
	store.Set(key, value, expiration)
	return stringToReply("OK")
}

func setStoreValueWithExpiration(key string, value string, expiration time.Time) string {
	store.Set(key, value, expiration)
	return stringToReply("OK")
}

func getStoreValue(key string) (string, bool) {
	return store.Get(key)
}
