package chip8

import (
	"fmt"
	"os"
)

const MEM_SIZE = 0xFFF

var NUMBER_SPRITES = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type CPU struct {
	Mem [MEM_SIZE]byte // available CPU memory
	Pc  uint16         // program counter

	Sp    uint16     // stack pointer
	Stack [16]uint16 // 16 bits register

	V  [16]uint8 // general purpose register
	I  uint16    // special direction I
	Dt uint8     // delay timer
	St uint8     // sound timer

	Gfx [32 * 64]byte // graphic memory, each value ( 1 or 0 ) represents the state of a pixel (on or off)

	WaitKey int
}

func (c *CPU) ResetGM() {
	for i := 0; i < len(c.Gfx); i++ {
		c.Gfx[i] = 0
	}
}

func LoadRom(cpu *CPU, path string) error {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading content from: \"%s\"\n%s", path, err.Error())
	}

	if len(fileContent) > MEM_SIZE-0x200 {
		return fmt.Errorf("rom size exeeds available memory")
	}

	copy(cpu.Mem[0x200:], fileContent)
	return nil
}

func InitializeCPU(cpu *CPU) {
	cpu.Pc = 0x200
	cpu.Sp = 0
	cpu.Dt = 0
	cpu.St = 0
	cpu.I = 0
	cpu.WaitKey = -1

	for i := 0; i < MEM_SIZE; i++ {
		cpu.Mem[i] = 0
	}

	for i := 0; i < 16; i++ {
		cpu.Stack[i] = 0
		cpu.V[i] = 0
	}

	n := copy(cpu.Mem[0x50:], NUMBER_SPRITES)
	if n == 0 {
		println("error loading number sprites")
	}
	cpu.ResetGM()
}
