package device

import (
	"emulator/utils"
	"fmt"
	"strings"
)

// 🟧 Handle keystrokes as fordwaed by UI

func (device *Device) HandleKeystroke(code string, key string, alt bool, ctrl bool, shift bool) {
	fmt.Printf("HandleKeystroke(code=%s key=%s alt=%t ctrl=%t shift=%t)\n", code, key, alt, ctrl, shift)
	// 👇 pre-analyze the key semantics
	attrs := device.attrs[device.addr]
	isData := len(key) == 1
	keyInProtected := isData && attrs.IsProtected()
	alphaInNumeric := isData && !strings.Contains("0123456789.", key) && attrs.IsNumeric()
	// 👇 we may be trying to go where no man is supposed to go!
	if isData && (keyInProtected || alphaInNumeric) {
		device.alarm = true
		// 👇 we can move the cursor anywhere we want to
	} else if strings.HasPrefix(code, "Arrow") {
		device.MoveCursor(code)
	} else if isData {
		u8 := utils.A2E([]byte(key))[0]
		device.UpdateByteAtCursor(u8)
	}
	// 👇 post-analyze the key semantics
	device.StatusForAttributes(device.attrs[device.addr])
	device.SignalStatus()
	device.RenderBuffer(RenderBufferOpts{blinkOn: true, quiet: true})
}

func (device *Device) MoveCursor(code string) {
	// 👇 reset changes stack
	device.changes = utils.NewStack[int](2)
	device.changes.Push(device.cursorAt)
	var cursorTo int
	switch code {
	case "ArrowDown":
		cursorTo = device.cursorAt + device.cols
		if cursorTo >= device.size {
			cursorTo = device.cursorAt % device.cols
		}
	case "ArrowLeft":
		cursorTo = device.cursorAt - 1
		if cursorTo < 0 {
			cursorTo = device.size - 1
		}
	case "ArrowRight":
		cursorTo = device.cursorAt + 1
		if cursorTo >= device.size {
			cursorTo = 0
		}
	case "ArrowUp":
		cursorTo = device.cursorAt - device.cols
		if cursorTo < 0 {
			cursorTo = (device.cursorAt % device.cols) + device.size - device.cols
		}
	}
	device.cursorAt = cursorTo
	device.addr = device.cursorAt
	device.changes.Push(device.cursorAt)
}
