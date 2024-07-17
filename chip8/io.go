package chip8

import (
	sdl "github.com/veandco/go-sdl2/sdl"
)

/**
	The computers which originally used the Chip-8 Language had a 16-key
	hexadecimal keypad with the following layout:

	original        remaped

	1 2 3 C         1 2 3 4
	4 5 6 D         q w e r
	7 8 9 E         a s d f
	A 0 B F         a s d f
**/

var EMULATOR_KEYS = []int{
	sdl.SCANCODE_X, // 0
	sdl.SCANCODE_1, // 1
	sdl.SCANCODE_2, // 2
	sdl.SCANCODE_3, // 3
	sdl.SCANCODE_Q, // 4
	sdl.SCANCODE_W, // 5
	sdl.SCANCODE_E, // 6
	sdl.SCANCODE_A, // 7
	sdl.SCANCODE_S, // 8
	sdl.SCANCODE_D, // 9
	sdl.SCANCODE_Z, // A
	sdl.SCANCODE_C, // B
	sdl.SCANCODE_4, // C
	sdl.SCANCODE_R, // D
	sdl.SCANCODE_F, // E
	sdl.SCANCODE_V, // F
}

func IsKeyPressed(key uint8) bool {
	keyboardState := sdl.GetKeyboardState()
	realKey := EMULATOR_KEYS[key]
	return keyboardState[realKey] == 1
}
