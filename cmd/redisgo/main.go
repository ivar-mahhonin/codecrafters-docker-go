package main

import (
	"fmt"
	"log"
	"os"

	server "github.com/ivar-mahhonin/redis-go/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	host := envVariable("CONN_HOST")
	port := envVariable("CONN_PORT")

	println(host, port)

	if host == "" || port == "" {
		host = "localhost"
		port = "6379"
		log.Println("Host and port must be set. Using default values")
	}

	server := server.NewTcpServer(host, port)
	server.Start()
}

func envVariable(key string) string {
	dir, _ := os.Getwd()

	err := godotenv.Load(fmt.Sprintf("%s/.env", dir))

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
