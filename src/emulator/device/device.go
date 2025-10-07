package device

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"strings"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// 🟧 Model the 3270 device in pure go test-able code. We are handed a drawing context into which we render the datastream and any operator input. See the go3270 package for how that context is actually drawn on an HTML canvas.

// 👁️ https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf
// 👁️ http://www.prycroft6.com.au/misc/3270.html
// 👁️ http://www.tommysprinkle.com/mvs/P3270/start.htm

type Device struct {
	bus  EventBus.Bus
	dc   *gg.Context
	face font.Face

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
	size         int

	// 👇 model the 3270 internals
	addr      int
	alarm     bool
	attrs     []*Attributes
	blinker   chan struct{}
	blinks    map[int]struct{}
	buffer    []byte
	changes   *Stack[int]
	command   byte
	cursorAt  int
	erase     bool
	error     bool
	focussed  bool
	locked    bool
	message   string
	numeric   bool
	protected bool
	waiting   bool

	// 👇 the glyph cache
	glyphs map[Glyph]image.Image
}

// 👁️ go3270.go for how pixels actually get drawn on the screen

type Glyph struct {
	u8         byte
	color      string
	reverse    bool
	underscore bool
}

// 👁️ go3270.go subscribes to them on the bus, then uses dispatchEvent to route them to mediator.ts in the UI. All this is necessary so that go test-able code can communicate with the ui without using syscall/js

type Go3270Message struct {
	args      []any
	u8s       []byte
	eventType string
	params    map[string]any
}

func NewDevice(
	bus EventBus.Bus,
	rgba *image.RGBA,
	face font.Face,
	bgColor string,
	color string,
	cols int,
	rows int,
	fontHeight float64,
	fontSize float64,
	fontWidth float64,
	paddedHeight float64,
	paddedWidth float64) *Device {
	device := new(Device)
	// 👇 initialize inherited properties
	device.bgColor = bgColor
	device.color = color
	device.cols = cols
	device.fontHeight = fontHeight
	device.fontSize = fontSize
	device.fontWidth = fontWidth
	device.paddedHeight = paddedHeight
	device.paddedWidth = paddedWidth
	device.rows = rows
	device.size = int(device.cols * device.rows)
	// 👇 initialize underlying data structures
	device.addr = 0
	device.attrs = make([]*Attributes, device.size)
	device.blinker = make(chan struct{})
	device.blinks = make(map[int]struct{}, 10)
	device.buffer = make([]byte, device.size)
	device.bus = bus
	device.dc = gg.NewContextForRGBA(rgba)
	device.face = face
	device.glyphs = make(map[Glyph]image.Image)
	// 👇 reset device status
	device.ResetStatus()
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
	baseline := y + h - (device.fontSize / 3)
	return x, y, w, h, baseline
}

func (device *Device) Close() {
	if device.blinker != nil {
		close(device.blinker)
		device.blinker = nil
	}
}

func (device *Device) EraseBuffer() {
	device.addr = 0
	clear(device.attrs)
	// 👇 initialize with protected fields
	for ix := range device.attrs {
		device.attrs[ix] = NewAttribute(0b00100000)
	}
	device.blinker = make(chan struct{})
	clear(device.blinks)
	clear(device.buffer)
	device.cursorAt = 0
	device.erase = true
}

func (device *Device) Focussed(focussed bool) {
	device.changes = NewStack[int](1)
	device.changes.Push(device.cursorAt)
	device.error = !focussed
	device.focussed = focussed
	device.message = Ternary(focussed, "", "LOCKED")
	device.SignalStatus()
	device.RenderBuffer(RenderBufferOpts{blinkOn: focussed, quiet: true})
}

func (device *Device) Keystroke(code string, key string, alt bool, ctrl bool, shift bool) {
	device.changes = NewStack[int](0)
	// 👇 pre-analyze the key semantics
	attrs := device.attrs[device.addr]
	isData := len(key) == 1
	keyInProtected := isData && attrs.Protected()
	alphaInNumeric := isData && !strings.Contains("0123456789.", key) && attrs.Numeric()
	switch {
	// 👇 we may be trying to go where no man is supposed to go!
	case isData && (keyInProtected || alphaInNumeric):
		device.alarm = true
		// 👇 we can move the cursor anywhere we want to
	case strings.HasPrefix(code, "Arrow"):
		device.KeystrokeToMoveCursor(code)
		// 👇 just data
	case isData:
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
	u8 := A2E([]byte(key))[0]
	device.addr = device.cursorAt
	device.buffer[device.addr] = u8
	device.changes.Push(device.addr)
	device.addr++
	// 👇 note wrap around
	if device.addr == device.size {
		device.addr = 0
	}
	device.cursorAt = device.addr
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
		device.SendGo3270Message(Go3270Message{eventType: "panic", args: []any{"🔥 RMA not handled"}})
	case CommandLookup["EAU"]:
		device.SendGo3270Message(Go3270Message{eventType: "panic", args: []any{"🔥 EAU not handled"}})
	case CommandLookup["EWA"]:
		device.EraseBuffer()
		device.ProcessWCC(out)
		device.ProcessOrdersAndData(out)
	case CommandLookup["W"]:
		device.ProcessWCC(out)
		device.ProcessOrdersAndData(out)
	case CommandLookup["RB"]:
		device.SendGo3270Message(Go3270Message{eventType: "panic", args: []any{"🔥 RB not handled"}})
	case CommandLookup["WSF"]:
		device.SendGo3270Message(Go3270Message{eventType: "panic", args: []any{"🔥 WSF not handled"}})
	case CommandLookup["EW"]:
		device.EraseBuffer()
		device.ProcessWCC(out)
		device.ProcessOrdersAndData(out)
	case CommandLookup["RM"]:
		device.SendGo3270Message(Go3270Message{eventType: "panic", args: []any{"🔥 RM not handled"}})
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
			println("🔥 PT not handled")
		case OrderLookup["GE"]:
			println("🔥 GE not handled")
		case OrderLookup["SBA"]:
			addr, _ := out.NextSlice(2)
			device.addr = AddrFromBytes(addr)
			if device.addr >= device.size {
				device.SendGo3270Message(Go3270Message{eventType: "panic", args: []any{"Data requires a device with a larger screen"}})
				return
			}
		case OrderLookup["EUA"]:
			println("🔥 EUA not handled")
		case OrderLookup["IC"]:
			device.cursorAt = device.addr
		case OrderLookup["SF"]:
			attrs, _ := out.NextSlice(1)
			lastAttrs = NewAttributes(attrs)
			device.PutByteIntoBuffer(OrderLookup["SF"], lastAttrs)
		case OrderLookup["SA"]:
			println("🔥 SA not handled")
		case OrderLookup["SFE"]:
			count, _ := out.Next()
			attrs, _ := out.NextSlice(int(count) * 2)
			lastAttrs = NewAttributes(attrs)
			device.PutByteIntoBuffer(OrderLookup["SF"], lastAttrs)
		case OrderLookup["MF"]:
			println("🔥 MF not handled")
		case OrderLookup["RA"]:
			println("🔥 RA not handled")
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
	if wcc.Alarm() {
		device.locked = false
	}
	if wcc.Unlock() {
		device.locked = false
	}
	if wcc.Reset() {
		// TODO implement DoReset
	}
	if wcc.ResetMDT() {
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
	device.addr++
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

func (device *Device) RenderBlinkingAddrs(quit <-chan struct{}) {
	for ix := 0; ; ix++ {
		select {
		case <-quit:
			return
		default:
			device.changes.Push(device.cursorAt)
			for addr := range device.blinks {
				device.changes.Push(addr)
			}
			// 🔥 after RenderBuffer is called, the "changes" stack is empty
			device.RenderBuffer(RenderBufferOpts{blinkOn: (ix % 2) == 0, quiet: true})
			time.Sleep(500 * time.Millisecond)
		}
	}
}

type RenderBufferOpts struct {
	quiet   bool
	blinkOn bool
}

func (device *Device) RenderBuffer(opts RenderBufferOpts) {
	defer ElapsedTime(time.Now(), "RenderBuffer", opts.quiet)
	// 👇 for example, EW command
	if device.erase {
		device.dc.SetHexColor(device.bgColor)
		device.dc.Clear()
	}
	// 🔥 don't do this until we're done because we need the flag
	defer func() { device.erase = false }()
	// 👇 if requested, dump the buffer contents
	if !opts.quiet {
		params := map[string]any{
			"color":  "cyan",
			"ebcdic": true,
			"title":  "RenderBuffer",
		}
		device.SendGo3270Message(Go3270Message{eventType: "dumpBytes", params: params, u8s: device.buffer})
	}
	// 👇 iterate over all changed cells
	for !device.changes.IsEmpty() {
		addr := device.changes.Pop()
		attrs := device.attrs[addr]
		cell := device.buffer[addr]
		color := attrs.Color(device.color)
		underscore := attrs.Underscore()
		visible := cell != 0x00 && !attrs.Hidden()
		// 👇 quick exit: if not visible, and we've already cleared the device, we don't have to do anything
		if !visible && device.erase {
			break
		}
		// 🔥 != here is the Go idiom for XOR
		showCursor := (addr == device.cursorAt) && device.focussed
		blinkMe := (attrs.Blink() || showCursor) && opts.blinkOn
		reverse := attrs.Reverse() != blinkMe
		x, y, w, h, baseline := device.BoundingBox(addr)
		// 👇 lookup the glyph in the cache
		glyph := Glyph{
			u8:         cell,
			color:      color,
			reverse:    reverse,
			underscore: underscore,
		}
		if img, ok := device.glyphs[glyph]; ok {
			// 👇 cache hit: just bitblt the glyph
			device.dc.DrawImage(img, int(x), int(y))
		} else {
			println("🔥 cache miss", cell)
			// 👇 cache miss: draw the glyph in a temporary context
			rgba := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
			temp := gg.NewContextForRGBA(rgba)
			temp.SetFontFace(device.face)
			// 👇 clear background
			temp.SetHexColor(Ternary(reverse, color, device.bgColor))
			temp.Clear()
			// 👇 render the byte
			temp.SetHexColor(Ternary(reverse, device.bgColor, color))
			str := string(E2A([]byte{cell}))
			temp.DrawString(str, 0, baseline-y)
			if underscore {
				temp.SetLineWidth(2)
				temp.MoveTo(0, h-1)
				temp.LineTo(w, h-1)
				temp.Stroke()
			}
			// 👇 now cache and bitblt the glyph
			device.glyphs[glyph] = temp.Image()
			device.dc.DrawImage(temp.Image(), int(x), int(y))
		}
	}
}

func (device *Device) ResetStatus() {
	device.alarm = false
	device.cursorAt = 0
	device.erase = false
	device.error = false
	device.focussed = true
	device.locked = true
	device.message = ""
	device.numeric = false
	device.protected = false
	device.waiting = false
}

func (device *Device) SendGo3270Message(msg Go3270Message) {
	device.bus.Publish("go3270", msg.eventType, msg.u8s, msg.params, msg.args)
}

func (device *Device) SignalStatus() {
	status := map[string]any{
		"alarm":     device.alarm,
		"cursorAt":  device.cursorAt,
		"error":     device.error,
		"locked":    device.locked,
		"message":   device.message,
		"numeric":   device.numeric,
		"protected": device.protected,
		"waiting":   device.waiting,
	}
	device.SendGo3270Message(Go3270Message{eventType: "status", params: status})
	device.alarm = false
}

func (device *Device) StatusForAttributes(attrs *Attributes) {
	device.numeric = attrs.Numeric()
	device.protected = attrs.Protected()
}
