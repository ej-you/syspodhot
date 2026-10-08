go_exec="./main.go"

dev:
	go run $(go_exec) serve

manager:
	go run $(go_exec) manager

lint:
	golangci-lint run
