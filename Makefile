install:
	go install

build:
	go build -o diploy

server: build
	./diploy
