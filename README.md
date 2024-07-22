[Chip-8](https://en.wikipedia.org/wiki/CHIP-8) emulator based on http://devernay.free.fr/hacks/chip8/C8TECH10.HTM with SDL2


See the following instructions to install [SDL2](https://www.libsdl.org/) depending on your operating system: https://github.com/veandco/go-sdl2?tab=readme-ov-file#installation 


Install dependencies
```sh
go mod tidy
```

run emulator
```sh
go run main.go -rom ./roms/INVADERS
```

![screenshoot](./Screenshot_space_invaders.webp)
