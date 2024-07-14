package chip8

// typedef unsigned char Uint8;
// void SineWave(void *userdata, Uint8 *stream, int len);
import "C"

import (
	"math"
	"reflect"
	"unsafe"

	sdl "github.com/veandco/go-sdl2/sdl"
)

const (
	PI      = 3.14159
	TONE    = 1000 // Hz
	SAMPLES = 44100
	D_PHASE = 2 * PI * TONE / SAMPLES
)

var phase float64

//export SineWave
func SineWave(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	for i := 0; i < n; i++ {
		buf[i] = C.Uint8((math.Sin(phase) + 0.9999) * 128)
		phase += D_PHASE
	}
}

func CreateAudioSpec() *sdl.AudioSpec {
	spec := &sdl.AudioSpec{
		Freq:     SAMPLES,
		Format:   sdl.AUDIO_U8,
		Channels: 1,
		Samples:  SAMPLES / 100,
		Callback: sdl.AudioCallback(C.SineWave),
	}
	return spec
}
