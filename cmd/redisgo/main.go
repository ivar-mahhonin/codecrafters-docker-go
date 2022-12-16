package main

import server "github.com/ivar-mahhonin/redis-go/internal/server"

const (
	CONN_HOST = "localhost"
	CONN_PORT = "6379"
	CONN_TYPE = "tcp"
)

func main() {
	server := server.NewTcpServer(CONN_HOST, CONN_PORT)
	server.Start()
}
