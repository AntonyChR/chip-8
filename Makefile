run:
	go run main.go -rom ./roms/INVADERS
debug:
	go build -gcflags="-N -L" -o chip8
build:
	go build -ldflags="-s -w" -o target/chip8

test:
	go test ./chip8/ -v
