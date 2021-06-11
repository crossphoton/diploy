install:
	go install

build: diploy
	go build -o diploy

server: build
	./diploy
