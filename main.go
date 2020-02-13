package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
// #include <X11/X.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

// Play a sound when coming out of idle.
// Directly mute / unmute.
// Have an 'idle', where we only check once per second.

// Super_R
const pttKey = 134

func main() {
	display := C.XOpenDisplay(nil)

	pttKeyByte := pttKey / 8
	pttKeyBit := pttKey % 8

	keys := [32]C.char{}
	for {
		C.XQueryKeymap(display, &keys[0])
		keyArr := C.GoBytes(unsafe.Pointer(&keys), 32)

		k := keyArr[pttKeyByte]
		// for i, k := range keyArr {
		mask := byte(1 << uint(pttKeyBit))
		if (mask & k) == mask {
			fmt.Printf("key %d is pressed\n", pttKeyBit+(pttKeyByte*8))
		}
		// }
		time.Sleep(time.Millisecond * 10)
	}
}
