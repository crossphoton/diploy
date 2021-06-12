install:
	go install

build:
	go build -o diploy

server: build
	DIPLOY_LOG_PATH=./ ./diploy

setup: build
	DIPLOY_LOG_PATH=./ ./diploy server setup