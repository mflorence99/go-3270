package device

import (
	"bytes"
	"fmt"
	"time"
)

// 🟧 Process the "outbound" datastream --- ie from app to 3270

func (device *Device) EraseBuffer() {
	device.addr = 0
	clear(device.attrs)
	// 👇 initialize with ptotected fields
	for ix := range device.attrs {
		device.attrs[ix] = NewAttribute(0b00100000)
	}
	device.blinker = make(chan struct{})
	clear(device.blinks)
	clear(device.buffer)
	device.cursorAt = 0
	device.erase = true
}

func (device *Device) MakeFramesFromBytes(u8s []byte) []*OutboundDataStream {
	slices := bytes.Split(u8s, LT)
	frames := make([]*OutboundDataStream, 0)
	for ix := range slices {
		if len(slices[ix]) > 0 {
			frame := NewOutboundDataStream(&slices[ix])
			frames = append(frames, frame)
		}
	}
	return frames
}

func (device *Device) ProcessCommands(out *OutboundDataStream) {
	defer ElapsedTime(time.Now(), "ProcessCommands")
	// 👇 dispatch on command
	switch device.command {
	case CommandLookup["RMA"]:
	case CommandLookup["EAU"]:
	case CommandLookup["EWA"]:
		device.EraseBuffer()
		device.ProcessWCC(out)
		device.ProcessOrdersAndData(out)
	case CommandLookup["W"]:
		device.ProcessWCC(out)
		device.ProcessOrdersAndData(out)
	case CommandLookup["RB"]:
	case CommandLookup["WSF"]:
	case CommandLookup["EW"]:
		device.EraseBuffer()
		device.ProcessWCC(out)
		device.ProcessOrdersAndData(out)
	case CommandLookup["RM"]:
	}
}

func (device *Device) ProcessOrdersAndData(out *OutboundDataStream) {
	defer ElapsedTime(time.Now(), "ProcessOrdersAndData")
	var lastAttrs *Attributes = NewAttribute(0b00000000)
	for out.HasNext() {
		// 👇 look at each order to see if it is an order
		order, _ := out.Next()
		// 👇 dispatch on order
		switch order {
		case OrderLookup["PT"]:
		case OrderLookup["GE"]:
		case OrderLookup["SBA"]:
			addr, _ := out.NextSlice(2)
			device.addr = AddrFromBytes(addr)
			if device.addr >= device.size {
				device.SendGo3270Message(Go3270Message{eventType: "panic", args: []any{"Data requires a device with a larger screen"}})
				return
			}
		case OrderLookup["EUA"]:
		case OrderLookup["IC"]:
			device.cursorAt = device.addr
		case OrderLookup["SF"]:
			attrs, _ := out.NextSlice(1)
			lastAttrs = NewAttributes(attrs)
			device.PutByteIntoBuffer(OrderLookup["SF"], lastAttrs)
		case OrderLookup["SA"]:
		case OrderLookup["SFE"]:
			count, _ := out.Next()
			attrs, _ := out.NextSlice(int(count) * 2)
			lastAttrs = NewAttributes(attrs)
			device.PutByteIntoBuffer(OrderLookup["SF"], lastAttrs)
		case OrderLookup["MF"]:
		case OrderLookup["RA"]:
		// 👇 if it isn't an order, it's data
		// 🔥 let's not convert the EBCDIC byte to ASCII until we actually need to, as we'll cache glyphs by their EDCDIC value
		default:
			if order == 0x00 || order >= 0x40 {
				device.PutByteIntoBuffer(order, lastAttrs)
			}
		}
	}
	// 👇 leave the buffer address at the last cursor position
	device.addr = device.cursorAt
}

func (device *Device) ProcessWCC(out *OutboundDataStream) {
	byte, err := out.Next()
	if err != nil {
		device.SendGo3270Message(Go3270Message{eventType: "panic", args: []any{fmt.Sprintf("Unable to extract WCC: %s", err.Error())}})
		return
	}
	wcc := NewWCC(byte)
	println(wcc.ToString())
	// 👇 honor WCC instructions
	if wcc.DoAlarm() {
		device.locked = false
	}
	if wcc.DoUnlock() {
		device.locked = false
	}
	if wcc.DoReset() {
		// TODO implement DoReset
	}
	if wcc.DoResetMDT() {
		// TODO implement DoReset
	}
}

func (device *Device) PutByteIntoBuffer(u8 byte, attrs *Attributes) {
	device.attrs[device.addr] = attrs
	if attrs.Blink() {
		device.blinks[device.addr] = struct{}{}
	} else {
		delete(device.blinks, device.addr)
	}
	device.buffer[device.addr] = u8
	device.changes.Push(device.addr)
	device.addr += 1
	// 👇 note wrap around
	if device.addr == device.size {
		device.addr = 0
	}
}

func (device *Device) ReceiveFromApp(u8s []byte) {
	// 👇 reset any binking
	if device.blinker != nil {
		close(device.blinker)
		device.blinker = nil
	}
	// 👇 reset changes stack
	device.changes = NewStack[int](device.size)
	// 👇 data can be split into multiple frames
	frames := device.MakeFramesFromBytes(u8s)
	for ix := range frames {
		fmt.Printf("ReceiveFromApp(frame #%d)\n", ix)
		// 👇 extract command
		out := frames[ix]
		cmd, err := out.Next()
		if err != nil {
			device.SendGo3270Message(Go3270Message{eventType: "panic", args: []any{fmt.Sprintf("Unable to extract write command: %s", err.Error())}})
			return
		}
		device.command = cmd
		println("COMMAND=", Command[device.command])
		// 👇 dispatch on command
		device.ProcessCommands(out)
	}
	// 👇 broadcast status
	device.SignalStatus()
	// 👇 now we can render the buffer to the drawing context
	// 🔥 after RenderBuffer is called, the "changes" stack is empty
	device.RenderBuffer(RenderBufferOpts{blinkOn: true})
	// 👇 start any blinking
	device.blinker = make(chan struct{})
	go device.RenderBlinkingAddrs(device.blinker)
}
