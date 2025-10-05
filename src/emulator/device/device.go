package device

import (
	"emulator/types"
	"emulator/utils"
	"fmt"
	"math"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
)

// 🟧 Model the 3270 device in pure go test-able code. We are handed a drawing context into which we render the datastream and any operator input. See the go3270 package for how that context is actually drawn on an HTML canvas.

// 👁️ https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf
// 👁️ http://www.prycroft6.com.au/misc/3270.html
// 👁️ http://www.tommysprinkle.com/mvs/P3270/start.htm

type Device struct {
	bus EventBus.Bus
	gg  *gg.Context

	// 👇 properties
	bgColor      string
	color        string
	cols         int
	fontHeight   float64
	fontSize     float64
	fontWidth    float64
	paddedHeight float64
	paddedWidth  float64
	rows         int
	scaleFactor  float64
	size         int

	// 👇 model the device buffer
	addr     int
	attrs    []*Attributes
	blinker  chan struct{}
	blinks   map[int]struct{}
	buffer   []uint8
	changes  *utils.Stack[int]
	command  uint8
	cursorAt int
	erase    bool
	wcc      *WCC
}

func NewDevice(
	bus EventBus.Bus,
	gg *gg.Context,
	bgColor string,
	color string,
	cols int,
	rows int,
	fontHeight float64,
	fontSize float64,
	fontWidth float64,
	paddedHeight float64,
	paddedWidth float64,
	scaleFactor float64) *Device {
	device := new(Device)
	device.bus = bus
	device.gg = gg
	// 👇 initialize properties
	device.bgColor = bgColor
	device.color = color
	device.cols = cols
	device.fontHeight = fontHeight
	device.fontSize = fontSize
	device.fontWidth = fontWidth
	device.paddedHeight = paddedHeight
	device.paddedWidth = paddedWidth
	device.rows = rows
	device.scaleFactor = scaleFactor
	device.size = int(device.cols * device.rows)
	// 👇 initialize buffer
	device.InitializeBuffer()
	return device
}

func (device *Device) BoundingBox(addr int) (float64, float64, float64, float64, float64) {
	col := addr % device.cols
	row := int(addr / device.cols)
	w := math.Round(device.fontWidth * device.paddedWidth)
	h := math.Round(device.fontHeight * device.paddedHeight)
	x := math.Round(float64(col) * w)
	y := math.Round(float64(row) * h)
	// 🔥 we could do better calculating the baseline - this is just a WAG, because an em is drawn with a significantly different height than that returned by MeasureString()
	baseline := y + h - (device.fontSize / 3 * device.scaleFactor)
	return x, y, w, h, baseline
}

func (device *Device) Close() {
	// 🔥 sorry I had to do this the hard way, here I wanted the colors
	SendMessage(Message{bus: device.bus, eventType: "log", args: []any{"%cDevice closing", "color: pink"}})
	close(device.blinker)
	device.blinker = make(chan struct{})
}

func (device *Device) InitializeBuffer() {
	device.addr = 0
	device.attrs = make([]*Attributes, device.size)
	device.blinker = make(chan struct{})
	// 👇 capacity is just a WAG, but blinking isn't common
	device.blinks = make(map[int]struct{}, 10)
	device.buffer = make([]uint8, device.size)
	device.cursorAt = 0
	device.erase = true
}

func (device *Device) MakeFramesFromBytes(bytes []uint8) []*OutboundDataStream {
	// 👇 we know there's going to be one frame, and more isn't common
	frames := make([]*OutboundDataStream, 0)
	whole := NewOutboundDataStream(&bytes)
	for {
		slice, err := whole.NextSliceUntil(types.LT)
		if len(slice) > 0 {
			frame := NewOutboundDataStream(&slice)
			frames = append(frames, frame)
		}
		if err != nil {
			break
		}
		whole.Skip(len(types.LT))
	}
	return frames
}

func (device *Device) PutBuffer(byte uint8, attrs *Attributes) {
	device.attrs[device.addr] = attrs
	if attrs.IsBlink() {
		device.blinks[device.addr] = struct{}{}
	} else {
		delete(device.blinks, device.addr)
	}
	device.buffer[device.addr] = byte
	device.changes.Push(device.addr)
	device.addr += 1
	// 👇 note wrap around
	if device.addr == device.size {
		device.addr = 0
	}
}

func (device *Device) ReceiveFromApp(bytes []uint8) {
	// 👇 reset changes stack
	close(device.blinker)
	device.blinker = make(chan struct{})
	device.changes = utils.NewStack[int](device.size)
	// 👇 data can be split into multiple frames
	frames := device.MakeFramesFromBytes(bytes)
	for ix := range frames {
		fmt.Printf("ReceiveFromApp(frame #%d)\n", ix)
		// 👇 extract command
		out := frames[ix]
		cmd, err := out.Next()
		if err != nil {
			SendMessage(Message{bus: device.bus, eventType: "panic", args: []any{fmt.Sprintf("Unable to extract write command: %s", err.Error())}})
			return
		}
		device.command = cmd
		fmt.Printf("COMMAND=%s\n", types.Command[device.command])
		// 👇 for all but WSF, extract WCC
		if device.command != types.CommandLookup["WSF"] {
			u8, err := out.Next()
			if err != nil {
				SendMessage(Message{bus: device.bus, eventType: "panic", args: []any{fmt.Sprintf("Unable to extract WCC: %s", err.Error())}})
				return
			}
			device.wcc = NewWCC(u8)
			fmt.Println(device.wcc.ToString())
		}
		// 👇 dispatch on command
		device.WriteCommands(out)
	}
	// 👇 now we can render the buffer to the drawing context --
	device.SignalStatus()
	// 🔥 after RenderBuffer is called, the "changes" stack is empty
	device.RenderBuffer(false, true)
	if device.cursorAt >= 0 {
		go device.RenderBlinkingAddrs()
	}
}

func (device *Device) RenderBlinkingAddrs() {
	for ix := 0; ; ix++ {
		select {
		case <-device.blinker:
			fmt.Println("device.BlinkingCursor() stopped")
			return
		default:
			device.changes.Push(device.cursorAt)
			for addr := range device.blinks {
				device.changes.Push(addr)
			}
			// 🔥 after RenderBuffer is called, the "changes" stack is empty
			device.RenderBuffer(true, (ix%2) == 0)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (device *Device) RenderBuffer(quiet bool, blinkOn bool) {
	if !quiet {
		defer utils.ElapsedTime(time.Now(), "RenderBuffer")
	}
	defer func() { device.erase = false }()
	// 👇 for example, EW command
	if device.erase {
		device.gg.SetHexColor(device.bgColor)
		device.gg.Clear()
	}
	// 👇 iterate over all changed cells
	for !device.changes.IsEmpty() {
		addr := device.changes.Pop()
		attrs := device.attrs[addr]
		byte := device.buffer[addr]
		// 🔥 Go idiom for XOR
		reverse := attrs.IsReverse() != ((attrs.IsBlink() || (addr == device.cursorAt)) && blinkOn)
		visible := byte != 0x00 && !attrs.IsHidden()
		x, y, w, h, baseline := device.BoundingBox(addr)
		// 👇 clear background, except when the device has already been cleared
		if (reverse && visible) || !device.erase {
			device.gg.SetHexColor(utils.Ternary(reverse, attrs.GetColor(device.color), device.bgColor))
			device.gg.DrawRectangle(x, y, w, h)
			device.gg.Fill()
		}
		// 👇 a zero byte is hidden, an SF or SFE order
		if visible {
			device.gg.SetHexColor(utils.Ternary(reverse, device.bgColor, attrs.GetColor(device.color)))
			str := string(byte)
			device.gg.DrawString(str, x, baseline)
		}
	}
}

func (device *Device) SignalStatus() {
	status := map[string]any{
		"alarm":     false,
		"cursorAt":  device.cursorAt,
		"error":     false,
		"locked":    false,
		"message":   "",
		"numeric":   false,
		"protected": false,
		"waiting":   false,
	}
	SendMessage(Message{bus: device.bus, eventType: "status", params: status})
}

func (device *Device) WriteCommands(out *OutboundDataStream) {
	defer utils.ElapsedTime(time.Now(), "WriteCommands")
	// 👇 dispatch on command
	switch device.command {
	case types.CommandLookup["RMA"]:
	case types.CommandLookup["EAU"]:
	case types.CommandLookup["EWA"]:
		device.InitializeBuffer()
		device.WriteOrdersAndData(out)
	case types.CommandLookup["W"]:
	case types.CommandLookup["RB"]:
	case types.CommandLookup["WSF"]:
	case types.CommandLookup["EW"]:
		device.InitializeBuffer()
		device.WriteOrdersAndData(out)
	case types.CommandLookup["RM"]:
	}
}

func (device *Device) WriteOrdersAndData(out *OutboundDataStream) {
	var lastAttrs *Attributes = NewAttributes([]uint8{0x00})
	for out.HasNext() {
		// 👇 look at each byte to see if it is an order
		order, _ := out.Next()
		// 👇 dispatch on order
		switch order {
		case types.OrderLookup["PT"]:
		case types.OrderLookup["GE"]:
		case types.OrderLookup["SBA"]:
			addr, _ := out.NextSlice(2)
			device.addr = utils.AddrFromBytes(addr)
			if device.addr >= device.size {
				SendMessage(Message{bus: device.bus, eventType: "panic", args: []any{"Data requires a device with a larger screen"}})
				return
			}
		case types.OrderLookup["EUA"]:
		case types.OrderLookup["IC"]:
			device.cursorAt = device.addr
		case types.OrderLookup["SF"]:
			attrs, _ := out.NextSlice(1)
			lastAttrs = NewAttributes(attrs)
			device.PutBuffer(0x00, lastAttrs)
		case types.OrderLookup["SA"]:
		case types.OrderLookup["SFE"]:
			count, _ := out.Next()
			attrs, _ := out.NextSlice(int(count) * 2)
			lastAttrs = NewAttributes(attrs)
			device.PutBuffer(0x00, lastAttrs)
		case types.OrderLookup["MF"]:
		case types.OrderLookup["RA"]:
		// 👇 if it isn't an order, it's data
		default:
			if order == 0x00 || order >= 0x40 {
				device.PutBuffer(utils.E2A([]uint8{order})[0], lastAttrs)
			}
		}
	}
}
