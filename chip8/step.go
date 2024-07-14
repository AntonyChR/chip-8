package chip8

import (
	"math/rand/v2"

	"github.com/veandco/go-sdl2/sdl"
)

func Step(cpu *CPU){
		opcode := (uint16(cpu.Mem[cpu.Pc]) << 8) | uint16(cpu.Mem[cpu.Pc+1])
		cpu.Pc = (cpu.Pc + 2) & 0xFFF

		// Get params
		nnn := opcode & 0x0FFF
		kk := uint8(opcode & 0xFF)
		x := uint8((opcode >> 8) & 0xF)
		y := uint8((opcode >> 4) & 0xF)

		// last 4 bits
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
				if(cpu.St != 0){
					println("sound on!!")
					sdl.PauseAudio(false)
				}
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

