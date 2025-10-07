package device

import (
	"emulator/types"
	"emulator/utils"
	"fmt"
	"strings"
)

// 🟧 Handle keystrokes as forwartded by UI

func (device *Device) FindFieldStartFromCursor() int {
	addr := device.cursorAt
	// 👇 prevent eendless loop if no field
	for ix := 0; ix < 2; ix++ {
		if device.buffer[addr] == types.OrderLookup["SF"] ||
			device.buffer[addr] == types.OrderLookup["SFE"] {
			return addr
		}
		addr -= 1
		if addr < 0 {
			addr = device.size - 1
		}
	}
	return -1
}

func (device *Device) Focussed(focussed bool) {
	device.changes = utils.NewStack[int](1)
	device.changes.Push(device.cursorAt)
	device.error = !focussed
	device.focussed = focussed
	device.message = utils.Ternary(focussed, "", "LOCKED")
	device.SignalStatus()
	device.RenderBuffer(RenderBufferOpts{blinkOn: focussed, quiet: true})
}

func (device *Device) Keystroke(code string, key string, alt bool, ctrl bool, shift bool) {
	fmt.Printf("Keystroke(code=%s key=%s alt=%t ctrl=%t shift=%t)\n", code, key, alt, ctrl, shift)
	device.changes = utils.NewStack[int](0)
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
		device.KeystrokeToMoveCursor(code)
	} else if isData {
		device.KeystrokeToSetByteAtCursor(key)
	}
	// 👇 post-analyze the key semantics
	device.StatusForAttributes(device.attrs[device.addr])
	device.SignalStatus()
	device.RenderBuffer(RenderBufferOpts{blinkOn: true, quiet: true})
}

func (device *Device) KeystrokeToMoveCursor(code string) {
	// 👇 reset changes stack
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

func (device *Device) KeystrokeToSetByteAtCursor(key string) {
	u8 := utils.A2E([]byte(key))[0]
	device.addr = device.cursorAt
	device.buffer[device.addr] = u8
	device.changes.Push(device.addr)
	device.addr += 1
	// 👇 note wrap around
	if device.addr == device.size {
		device.addr = 0
	}
	device.cursorAt = device.addr
}
