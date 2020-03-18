
build:
	go build -o ./bin/nutactl

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/nutactl

run:
	go run .

format:
	go fmt
