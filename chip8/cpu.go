package chip8

import (
	"fmt"
	"os"
)
//	Memory Map:
//	+---------------+= 0xFFF (4095) End of Chip-8 RAM
//	|               |
//	|               |
//	|               |
//	|               |
//	|               |
//	| 0x200 to 0xFFF|
//	|     Chip-8    |
//	| Program / Data|
//	|     Space     |
//	|               |
//	|               |
//	|               |
//	+- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
//	|               |
//	|               |
//	|               |
//	+---------------+= 0x200 (512) Start of most Chip-8 programs
//	| 0x000 to 0x1FF|
//	| Reserved for  |
//	|  interpreter  |
//	+---------------+= 0x000 (0) Start of Chip-8 RAM

const MEM_SIZE = 0xFFF

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

// Reset graphic memory, each pixel turn off
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
