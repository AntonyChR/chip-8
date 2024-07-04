package chip8

import (
	"os"
	"testing"

)

func TestLoadRom(t *testing.T) {
	PONG_ROM_PATH := "../roms/PONG"
	cpu := CPU{}
	err := LoadRom(&cpu,PONG_ROM_PATH)

	if err != nil {
		t.Errorf("error reading rom data, %s", err.Error())
	}

	if cpu.Mem[0x200] == 0 {
		t.Error("memory allocation error, rom data must be start at 0x200")
	}

	romContent, err := os.ReadFile(PONG_ROM_PATH)

	if err != nil {
		t.Error("Error reading test rom, ", err.Error())
	}

	if romContent[0] != cpu.Mem[0x200] {
		t.Errorf("Error loading ROM, inconsistent data, Mem[0x200] should be %x, instead %x", romContent[0], cpu.Mem[0x200])
	}

}
