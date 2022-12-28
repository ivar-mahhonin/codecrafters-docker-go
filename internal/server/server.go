package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	redisStore "github.com/ivar-mahhonin/redis-go/internal/store"
)

const (
	CMND_ECHO = "echo"
	CMND_GET  = "get"
	CMND_SET  = "set"
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
		command = []string{message}
	}

	if command != nil && response == "" {
		response = executeCommand(command)
	}
	return response
}

func isReplArray(command string) bool {
	return strings.HasPrefix(command, "*")
}

func extractArgumentsFromReplArray(command string) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(command))
	scanner.Scan()
	numArgs, err := stringToInt(scanner.Text()[1:])

	if err != nil {
		log.Println("Failed to convert numArgs to int")
		return nil, errors.New("failed to convert numArgs to int")
	}

	arguments := make([]string, numArgs)

	for i := 0; i < numArgs; i++ {
		scanner.Scan()
		argLength, err := stringToInt(scanner.Text()[1:])

		if err != nil {
			log.Println("Failed to convert numArgs to int")
			return nil, errors.New("failed to convert numArgs to int")
		}

		scanner.Scan()
		arguments[i] = scanner.Text()[:argLength]
	}

	return arguments, nil
}

func makeErrorResponse(message string) string {
	return fmt.Sprintf("-%s\r\n", message)
}

func executeCommand(command []string) string {
	response := ""
	switch strings.ToLower(command[0]) {
	default:
		response = stringToReply("PONG")
	case CMND_ECHO:
		response = stringToReply(command[1])
	case CMND_SET:
		key := command[1]
		value := command[2]
		response = setValue(key, value)
	case CMND_GET:
		key := command[1]
		response = getVlaue(key)
	}
	return response
}

func setValue(key string, value string) string {
	store.Set(key, value)
	return stringToReply("OK")
}

func getVlaue(key string) string {
	return stringToReply(store.Get(key))
}

func stringToReply(message string) string {
	return fmt.Sprintf("+%s\r\n", message)
}

func stringToInt(str string) (int, error) {
	i, err := strconv.Atoi(str)
	return i, err
}
