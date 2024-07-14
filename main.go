package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	chip8 "github.com/AntonyChR/chip-8/chip8"
	sdl "github.com/veandco/go-sdl2/sdl"
)

const (
	WIDTH  = 64
	HEIGHT = 32

	WINDOW_WIDTH  = 640
	WINDOW_HEIGHT = 320

	PIXEL_SIZE = 10

	FPS = 30
)

func main() {
	var ROM = flag.String("rom", "", "-rom path/to/rom")

	flag.Parse()

	cpu := chip8.CPU{}
	chip8.InitializeCPU(&cpu)

	err := chip8.LoadRom(&cpu, *ROM)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Initilize sdl
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("Error initializing SDL: %v", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("CHIP-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Error creating window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Error creating renderer: %v", err)
	}
	defer renderer.Destroy()

	audioSpec := chip8.CreateAudioSpec()
	if err := sdl.OpenAudio(audioSpec, nil); err != nil {
		log.Println(err)
		return
	}

	sdl.PauseAudio(true)
	defer sdl.CloseAudio()

	var cycles uint64 = 0
	var lastTick uint64 = 0

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		if sdl.GetTicks64()-cycles > 1 {
			if cpu.WaitKey == -1 {
				chip8.Step(&cpu)
			} else {
				var key uint8
				for key = 0; key <= 0xF; key++ {
					if chip8.IsKeyPressed(key) {
						cpu.V[cpu.WaitKey] = key
						cpu.WaitKey = -1
						break
					}
				}
			}
			cycles = sdl.GetTicks64()
		}

		if sdl.GetTicks64()-lastTick > (1000 / FPS) {
			if cpu.Dt != 0 {
				cpu.Dt--
			}
			if cpu.St != 0 {
				cpu.St--
				if cpu.St == 0{
					sdl.PauseAudio(true)
				} else {
					sdl.PauseAudio(false)
				}
			}

			render(renderer, &cpu)
			lastTick = sdl.GetTicks64()
		}
	}
}

func render(renderer *sdl.Renderer, cpu *chip8.CPU) {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			if cpu.Gfx[y*WIDTH+x] == 1 {
				renderer.SetDrawColor(255, 255, 255, 255)
				rect := sdl.Rect{X: int32(x * PIXEL_SIZE), Y: int32(y * PIXEL_SIZE), H: PIXEL_SIZE, W: PIXEL_SIZE}
				renderer.FillRect(&rect)
			}
		}
	}
	renderer.Present()
}
