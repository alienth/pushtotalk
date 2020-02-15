package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
// #include <X11/X.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/jfreymuth/pulse"
	"github.com/jfreymuth/pulse/proto"
)

// Play a sound when coming out of idle.
// Directly mute / unmute.
// Have an 'idle', where we only check once per second.

// Super_R
const pttKey = 134

const micSourceIndex = 3

var muted = true

var pulseClient *pulse.Client

func main() {
	display := C.XOpenDisplay(nil)

	var err error
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	soundPath = fmt.Sprintf("%s/%s", execDir, "ptt.wav")

	pulseClient, err = pulse.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	pttKeyByte := pttKey / 8
	pttKeyBit := pttKey % 8
	pttKeyMask := byte(1 << uint(pttKeyBit))

	keys := [32]C.char{}
	for {
		C.XQueryKeymap(display, &keys[0])
		keyArr := C.GoBytes(unsafe.Pointer(&keys), 32)

		if (pttKeyMask & keyArr[pttKeyByte]) == pttKeyMask {
			muteSource(micSourceIndex, false)
		} else {
			muteSource(micSourceIndex, true)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

var soundPath string

func muteSource(source int, mute bool) {
	if mute != muted {
		muteReq := proto.SetSourceMute{SourceIndex: micSourceIndex, Mute: mute}
		err := pulseClient.RawRequest(&muteReq, nil)
		if err != nil {
			log.Println(err)
		}

		cmd := exec.Command("aplay", soundPath)
		cmd.Run()
		muted = mute
	}
}
