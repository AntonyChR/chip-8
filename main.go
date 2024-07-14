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

	screenFrecuency := 60
	deltaTime := 1000 / screenFrecuency // ms
	stepsPerCycle := 16
	running := true

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		for range stepsPerCycle {
			cpu.Step()
		}

		// The delay timer is active whenever the delay timer register (DT) is non-zero.
		// This timer does nothing more than subtract 1 from the value of DT at a rate of 60Hz.
		// When DT reaches 0, it deactivates.
		if cpu.Dt > 0 {
			cpu.Dt--
		}

		// The sound timer is active whenever the sound timer register (ST) is non-zero.
		// This timer also decrements at a rate of 60Hz, however, as long as ST's value is greater
		// than zero, the Chip-8 buzzer will sound. When ST reaches zero, the sound timer deactivates.
		if cpu.St > 0 {
			cpu.St--
			sdl.PauseAudio(false)
		} else {
			sdl.PauseAudio(true)
		}

		render(renderer, &cpu)

		sdl.Delay(uint32(deltaTime))
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
