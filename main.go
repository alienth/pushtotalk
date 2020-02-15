package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
// #include <X11/X.h>
// #include <stdlib.h>
import "C"

import (
	"flag"
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

var muted = true

var pulseClient *pulse.Client

func main() {
	micIndexFlag := flag.Int("mic-index", -1, "Source index of mic.")
	micNameFlag := flag.String("mic-name", "", "Source name of mic.")
	keyCodeFlag := flag.Int("key-code", 134, "Key code of PTT key.")

	flag.Parse()

	micIndex := *micIndexFlag
	micName := *micNameFlag
	keyCode := *keyCodeFlag
	if (micIndex == -1 && micName == "") || (micIndex != -1 && micName != "") {
		fmt.Println("Must specify one of --mic-index or --mic-name.")
		os.Exit(-1)
	}

	var err error
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	soundPath := fmt.Sprintf("%s/%s", execDir, "ptt.wav")

	pulseClient, err = pulse.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	var muteReq proto.SetSourceMute
	if micIndex != -1 {
		muteReq.SourceIndex = uint32(micIndex)
	} else {
		muteReq.SourceName = micName
	}
	setMute := func(mute bool) {
		if mute != muted {
			muteReq.Mute = mute
			err := pulseClient.RawRequest(&muteReq, nil)
			if err != nil {
				log.Println(err)
			}

			cmd := exec.Command("aplay", soundPath)
			cmd.Run()
			muted = mute
		}
	}

	watchForKey(keyCode, setMute)
}

func watchForKey(pttKey int, callback func(bool)) {
	display := C.XOpenDisplay(nil)
	pttKeyByte := pttKey / 8
	pttKeyBit := pttKey % 8
	pttKeyMask := byte(1 << uint(pttKeyBit))

	keys := [32]C.char{}
	for {
		C.XQueryKeymap(display, &keys[0])
		keyArr := C.GoBytes(unsafe.Pointer(&keys), 32)

		if (pttKeyMask & keyArr[pttKeyByte]) == pttKeyMask {
			callback(false)
		} else {
			callback(true)
		}
		time.Sleep(time.Millisecond * 10)
	}

}
