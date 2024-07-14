package chip8

import (
	"fmt"
	"math/rand/v2"
	"os"
)

/**

	Memory Map:
	+---------------+= 0xFFF (4095) End of Chip-8 RAM
	|               |
	|               |
	|               |
	|               |
	|               |
	| 0x200 to 0xFFF|
	|     Chip-8    |
	| Program / Data|
	|     Space     |
	|               |
	|               |
	|               |
	+- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
	|               |
	|               |
	|               |
	+---------------+= 0x200 (512) Start of most Chip-8 programs
	| 0x000 to 0x1FF|
	| Reserved for  |
	|  interpreter  |
	+---------------+= 0x000 (0) Start of Chip-8 RAM

**/

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

func (c *CPU) Step() {
	if c.WaitKey != -1 {
		var key uint8
		for key = 0; key <= 0xF; key++ {
			if IsKeyPressed(key) {
				c.V[c.WaitKey] = key
				c.WaitKey = -1
				break
			}
		}
		if c.WaitKey != 1 {
			return
		}
	}
	c.processOpcode()
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

// Take the instruction (opcode) from memory location cpu.Pc (program counter)
// each opcode is processed according to the specification: http://devernay.free.fr/hacks/chip8/C8TECH10.HTM
func (cpu *CPU) processOpcode() {
	opcode := (uint16(cpu.Mem[cpu.Pc]) << 8) | uint16(cpu.Mem[cpu.Pc+1])
	cpu.Pc = (cpu.Pc + 2) & 0xFFF

	// 0x|p|H|H|H| 16 bit opcode
	//   | |n|n|n| or addr - A 12-bit value, the lowest 12 bits of the instruction
	//   | | | |n| or nibble - A 4-bit value, the lowest 4 bits of the instruction
	//   | |x| | | - A 4-bit value, the lower 4 bits of the high byte of the instruction
	//   | | |y| | - A 4-bit value, the upper 4 bits of the low byte of the instruction
	//   | | |k|k| or byte - An 8-bit value, the lowest 8 bits of the instruction
	nnn := opcode & 0x0FFF
	kk := uint8(opcode & 0xFF)
	x := uint8((opcode >> 8) & 0xF)
	y := uint8((opcode >> 4) & 0xF)

	// To process the opcodes we take the last 4 bits on the left as the identifier "p",
	// in the same way if several opcodes match the same identifier "p", we take "n" (the last 4 bits on the right)
	n := uint8(opcode & 0xF)
	p := uint8(opcode >> 12)

	switch p {
	case 0:
		if opcode == 0x00E0 {
			// CLS - clear display
			cpu.ResetGM()
		} else if opcode == 0x00EE {
			// return from a subroutine
			if cpu.Sp > 0 {
				cpu.Sp--
				cpu.Pc = cpu.Stack[cpu.Sp]
			}
		}
	case 1:
		// JP: jump program counter to nnn
		cpu.Pc = nnn
	case 2:
		// 2nnn, increments the stack pointer, then puts the current Pc on the top of the stack
		if cpu.Sp < 16 {
			cpu.Stack[cpu.Sp] = cpu.Pc
			cpu.Sp++
		}
		cpu.Pc = nnn
	case 3:
		// 3xkk - SE V[x], Skip next instruction if V[x] = kk.
		if cpu.V[x] == kk {
			cpu.Pc = (cpu.Pc + 2) & 0xFFF
		}
	case 4:
		// 4xkk - SNE V[x], Skip next instruction if V[x] != kk
		if cpu.V[x] != kk {
			cpu.Pc = (cpu.Pc + 2) & 0xFFF
		}
	case 5:
		// 5xy0 - SE V[x], V[y], Skip next instruction if V[x] == V[y]
		if cpu.V[x] == cpu.V[y] {
			cpu.Pc = (cpu.Pc + 2) & 0xFFF
		}
	case 6:
		// 6xkk - LD V[x], Set V[x] = kk
		cpu.V[x] = kk
	case 7:
		// 7xkk - ADD V[x], Set V[x] = (V[x] + kk) & 0xFF
		cpu.V[x] = (cpu.V[x] + kk) & 0xFF
	case 8:
		// several instructions contain 8 as identifier,
		// but we can use the last 4 bits to determinate the instruction
		switch n {
		case 0:
			// 8xy0 - LD V[x], V[y], Set V[x] = V[y]
			cpu.V[x] = cpu.V[y]

		case 1:
			// 8xy1 - OR V[x], V[y], Set V[x] = V[x] | V[y]
			cpu.V[x] |= cpu.V[y]

		case 2:
			// 8xy2 - AND V[x], V[y], Set V[x] = V[x] & V[y]
			cpu.V[x] &= cpu.V[y]

		case 3:
			// 8xy3 - AND V[x], V[y], Set V[x] = V[x] ^ V[y]
			cpu.V[x] ^= cpu.V[y]

		case 4:
			// 8xy4 - ADD V[x], V[y], Set V[x] = V[x] + V[y], set V[f] = carry
			//if cpu.V[x] > cpu.V[x]+cpu.V[y]{
			if cpu.V[x]+cpu.V[y] > 0xFF {
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[x] = (cpu.V[x] + cpu.V[y]) & 0xFF
			//cpu.V[x] += cpu.V[y]
		case 5:
			// 8xy5 - SUB V[x], V[y], Set V[x] = V[x] - V[y], if V[x] > V[y] then V[0xF] is set to 1 otherwise 0

			if cpu.V[x] > cpu.V[y] {
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}

			cpu.V[x] = cpu.V[x] - cpu.V[y]

		case 6:
			// 8xy6 - SHR, set V[x] = V[x] >> 1, If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
			cpu.V[0xF] = (cpu.V[x] & 1)
			cpu.V[x] >>= 1

		case 7:
			// 8xy7 - SUB, V[x], V[y], Set V[x] = V[y] - V[x], if V[x] < V[y] then V[0xF] is set to 1 otherwise 0
			if cpu.V[x] < cpu.V[y] {
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[x] = cpu.V[y] - cpu.V[x]

		case 0xE:
			// 8xyE - SHL V[x], V[x] = V[x] << 1, If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2
			msb := cpu.V[x] >> 7
			//msb := cpu.V[x] & 0x80
			if msb != 0 {
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[x] <<= 1
		}

	case 9:
		// 9xy0 - SNE v[x], V[y], Skip next instruction if V[x] != V[y]
		if cpu.V[x] != cpu.V[y] {
			cpu.Pc = (cpu.Pc + 2) & 0xFFF
		}

	case 0xA:
		// Annn - LD I, the value of register I is set to nnn
		cpu.I = nnn
	case 0xB:
		// Bnnn - JP V[0], the program counter is set to nnn plus the value of V[0]
		cpu.Pc = (uint16(cpu.V[0]) + nnn) & 0xFFF
	case 0xC:
		// Cxkk - RND V[x], set V[x] = (random number 0-255) AND kk
		cpu.V[x] = uint8(rand.IntN(0x100)) & kk
	case 0xD:
		// Dxyn - DRW V[x], V[y],
		cpu.V[15] = 0
		var j uint8
		var i uint8

		// We take the sprite and calculate its position and pixel
		// value based on the size of the graphics memory (32x64)
		for j = 0; j < n; j++ {
			sprite := cpu.Mem[cpu.I+uint16(j)]
			for i = 0; i < 8; i++ {
				px := (cpu.V[x] + i) & 63
				py := (cpu.V[y] + j) & 31
				pos := 64*uint16(py) + uint16(px)

				var pixel uint8
				if (sprite & (1 << (7 - i))) != 0 {
					pixel = 1
				} else {
					pixel = 0
				}
				cpu.V[15] |= cpu.Gfx[pos] & pixel
				cpu.Gfx[pos] ^= pixel
			}
		}
	case 0xE:
		if kk == 0x9E {
			if IsKeyPressed(cpu.V[x]) {
				cpu.Pc = (cpu.Pc + 2) & 0xFFF
			}
		} else if kk == 0xA1 {
			if !IsKeyPressed(cpu.V[x]) {
				cpu.Pc = (cpu.Pc + 2) & 0xFFF
			}
		}
	case 0xF:
		switch kk {
		case 0x07:
			// Fx07 - LD V[x], V[x] = DT (delay timer value)
			cpu.V[x] = cpu.Dt
		case 0x0A:
			cpu.WaitKey = int(x)
		case 0x15:
			// Fx15 - LD DT, V[x], DT(delay timer) = V[x]
			cpu.Dt = cpu.V[x]
		case 0x18:
			// Fx07 - LD ST V[x], ST (sound timer)= V[x]
			cpu.St = cpu.V[x]
		case 0x1E:
			cpu.I += uint16(cpu.V[x])
		case 0x29:
			cpu.I = 0x50 + (uint16(cpu.V[x])&0xF)*5
		case 0x33:
			cpu.Mem[cpu.I+2] = cpu.V[x] % 10
			cpu.Mem[cpu.I+1] = (cpu.V[x] / 10) % 10
			cpu.Mem[cpu.I] = cpu.V[x] / 100
		case 0x55:
			var reg uint16
			for reg = 0; reg <= uint16(x); reg++ {
				cpu.Mem[cpu.I+reg] = cpu.V[reg]
			}
		case 0x65:
			for reg := 0; reg <= int(x); reg++ {
				cpu.V[reg] = cpu.Mem[cpu.I+uint16(reg)]
			}
		}
	}
}
