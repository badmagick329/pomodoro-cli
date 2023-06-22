.PHONY: build clean

BINARY_NAME=pomodoro
BUILD_DIR=./bin
CONFIG_FILE=go_pomodoro_config.json
SOUND_FILE=bell.mp3

build:
	mkdir -p ${BUILD_DIR}
	GOARCH=amd64 GOOS=linux go build -o ${BUILD_DIR}/${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=windows go build -o ${BUILD_DIR}/${BINARY_NAME}-windows.exe main.go
	cp ${CONFIG_FILE} ${BUILD_DIR}
	cp ${SOUND_FILE} ${BUILD_DIR}

run: build
	./${BINARY_NAME}

clean:
	go clean
	find ${BUILD_DIR} -name "${BINARY_NAME}-*" -type f -delete
	find . -name "cover.*" -type f -delete

test:
	go test ./... -v -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
