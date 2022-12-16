run_server:
	go run ./cmd/redisgo/main.go
run_tests:
	go test -v ./internal/server -count=1