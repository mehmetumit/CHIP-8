package chip8

// typedef unsigned char Uint8;
// void OnAudioPlayback(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"reflect"
	"unsafe"
)

var (
	audio  []byte
	offset int // We use this to keep track of which part of audio to play
	spec   *sdl.AudioSpec
	dev    sdl.AudioDeviceID
)

//export OnAudioPlayback
func OnAudioPlayback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]byte)(unsafe.Pointer(&hdr))
	for i := 0; i < n; i++ {
		buf[i] = audio[offset]
		offset = (offset + 1) % len(audio) // Increase audio offset and loop when it reaches the end
	}
}

func init() {
	var err error

	audio, spec = sdl.LoadWAV("./sounds/beep.wav")
	if spec == nil {
		log.Fatal("Audio load error!", sdl.GetError())
	}
	spec.Callback = sdl.AudioCallback(C.OnAudioPlayback)
	// Open default playback device
	if dev, err = sdl.OpenAudioDevice("", false, spec, nil, 0); err != nil {
		log.Fatal("Audio device couldn't open!", err)
	}
}
func ClouseAudio() {
	sdl.CloseAudioDevice(dev)
}
func PlayAudio() {
	log.Print("Playing audio...")
	// Start playback audio of device
	sdl.PauseAudioDevice(dev, false)
}
func PauseAudio() {
	log.Print("Audio paused!")
	// Stop playback audio of device
	sdl.PauseAudioDevice(dev, true)
}
