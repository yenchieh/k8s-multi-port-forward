.PHONY: run build

run:
	go run main.go

build:
	go build -o ./build/main main.go
