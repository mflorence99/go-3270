package device

import (
	"emulator/types"
	"emulator/utils"
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
	changes   *utils.Stack[int]
	command   byte
	cursorAt  int
	erase     bool
	error     bool
	locked    bool
	message   string
	numeric   bool
	protected bool
	waiting   bool

	// 👇 the glyph cache
	glyphs map[Glyph]image.Image
}

type Glyph struct {
	byte       byte
	color      string
	reverse    bool
	underscore bool
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
	// 👇 initialize underlying data structures
	device.bus = bus
	device.dc = gg.NewContextForRGBA(rgba)
	device.face = face
	device.glyphs = make(map[Glyph]image.Image)
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
	// 👇 initial device status
	device.addr = 0
	device.alarm = false
	device.attrs = make([]*Attributes, device.size)
	device.blinker = make(chan struct{})
	device.blinks = make(map[int]struct{}, 10)
	device.buffer = make([]byte, device.size)
	device.cursorAt = 0
	device.erase = false
	device.error = false
	device.locked = true
	device.message = ""
	device.numeric = false
	device.protected = false
	device.waiting = false
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
	// 🔥 sorry I had to do this the hard way, here I wanted the colors
	SendMessage(Message{bus: device.bus, eventType: "log", args: []any{"%cDevice closing", "color: pink"}})
	if device.blinker != nil {
		close(device.blinker)
		device.blinker = nil
	}
}

func (device *Device) EraseBuffer() {
	device.addr = 0
	device.attrs = make([]*Attributes, device.size)
	// 👇 initialize with ptotected fields
	for ix := range device.attrs {
		device.attrs[ix] = NewAttributes([]byte{0b00100000})
	}
	device.blinker = make(chan struct{})
	device.blinks = make(map[int]struct{})
	device.buffer = make([]byte, device.size)
	device.cursorAt = 0
	device.erase = true
}

func (device *Device) HandleKeystroke(code string, key string, alt bool, ctrl bool, shift bool) {
	fmt.Printf("HandleKeystroke(code=%s key=%s alt=%t ctrl=%t shift=%t)\n", code, key, alt, ctrl, shift)
	// 👇 pre-analyze the key semantics
	var attrs *Attributes
	attrs = device.attrs[device.addr]
	isData := len(key) == 1
	keyInProtected := isData && attrs.IsProtected()
	alphaInNumeric := isData && !strings.Contains("0123456789.", key) && attrs.IsNumeric()
	// 👇 we may be trying to go where no man is supposed to go!
	if device.locked || keyInProtected || alphaInNumeric {
		device.alarm = true
		device.error = true
		device.message = "LOCKED"
		// 👇 we can move the cursor anywhere we want to
	} else if strings.HasPrefix(code, "Arrow") {
		device.MoveCursor(code)
	}
	// 👇 post-analyze the key semantics
	attrs = device.attrs[device.addr]
	cursorInNumeric := attrs.IsNumeric()
	cursorInProtected := attrs.IsProtected()
	// 👇 broadcast status
	device.numeric = cursorInNumeric
	device.protected = cursorInProtected
	device.SignalStatus()
}

func (device *Device) MakeFramesFromBytes(bytes []byte) []*OutboundDataStream {
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

// TODO 🔥 experimental
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
	device.RenderBuffer(RenderBufferOpts{blinkOn: true})
}

func (device *Device) ProcessCommands(out *OutboundDataStream) {
	defer utils.ElapsedTime(time.Now(), "ProcessCommands")
	// 👇 dispatch on command
	switch device.command {
	case types.CommandLookup["RMA"]:
	case types.CommandLookup["EAU"]:
	case types.CommandLookup["EWA"]:
		device.EraseBuffer()
		device.ProcessWCC(out)
		device.ProcessOrdersAndData(out)
	case types.CommandLookup["W"]:
		device.ProcessWCC(out)
		device.ProcessOrdersAndData(out)
	case types.CommandLookup["RB"]:
	case types.CommandLookup["WSF"]:
	case types.CommandLookup["EW"]:
		device.EraseBuffer()
		device.ProcessWCC(out)
		device.ProcessOrdersAndData(out)
	case types.CommandLookup["RM"]:
	}
}

func (device *Device) ProcessOrdersAndData(out *OutboundDataStream) {
	defer utils.ElapsedTime(time.Now(), "ProcessOrdersAndData")
	var lastAttrs *Attributes = NewAttributes([]byte{0x00})
	for out.HasNext() {
		// 👇 look at each byte to see if it is an order
		byte, _ := out.Next()
		// 👇 dispatch on order
		switch byte {
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
		// 🔥 let's not convert the EBCDIC byte to ASCII until we actually need to, as we'll cache glyphs by their EDCDIC value
		default:
			if byte == 0x00 || byte >= 0x40 {
				device.PutBuffer(byte, lastAttrs)
			}
		}
	}
	// 👇 leave the buffer address at the last cursor position
	device.addr = device.cursorAt
}

func (device *Device) ProcessWCC(out *OutboundDataStream) {
	byte, err := out.Next()
	if err != nil {
		SendMessage(Message{bus: device.bus, eventType: "panic", args: []any{fmt.Sprintf("Unable to extract WCC: %s", err.Error())}})
		return
	}
	wcc := NewWCC(byte)
	fmt.Println(wcc.ToString())
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

func (device *Device) PutBuffer(byte byte, attrs *Attributes) {
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

func (device *Device) ReceiveFromApp(bytes []byte) {
	// 👇 reset any binking
	if device.blinker != nil {
		close(device.blinker)
		device.blinker = nil
	}
	// 👇 reset changes stack
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
	defer utils.ElapsedTime(time.Now(), "RenderBuffer", opts.quiet)
	// 👇 for example, EW command
	if device.erase {
		device.dc.SetHexColor(device.bgColor)
		device.dc.Clear()
	}
	// 🔥 don't do this until we're done because we need the flag
	defer func() { device.erase = false }()
	// 👇 iterate over all changed cells
	for !device.changes.IsEmpty() {
		addr := device.changes.Pop()
		attrs := device.attrs[addr]
		cell := device.buffer[addr]
		color := attrs.GetColor(device.color)
		underscore := attrs.IsUnderscore()
		visible := cell != 0x00 && !attrs.IsHidden()
		// 👇 quick exit: if not visible, and we've already cleared the device, we don't have to do anything
		if !visible && device.erase {
			break
		}
		// 🔥 != here is the Go idiom for XOR
		blinkMe := (attrs.IsBlink() || (addr == device.cursorAt)) && opts.blinkOn
		reverse := attrs.IsReverse() != blinkMe
		x, y, w, h, baseline := device.BoundingBox(addr)
		// 👇 lookup the glyph in the cache
		glyph := Glyph{
			byte:       cell,
			color:      color,
			reverse:    reverse,
			underscore: underscore,
		}
		if img, ok := device.glyphs[glyph]; ok {
			// 👇 cache hit: just bitblt the glyph
			device.dc.DrawImage(img, int(x), int(y))
		} else {
			// 👇 cache hit: draw the glyph in a temporary context
			rgba := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
			temp := gg.NewContextForRGBA(rgba)
			temp.SetFontFace(device.face)
			// 👇 clear background
			temp.SetHexColor(utils.Ternary(reverse, color, device.bgColor))
			temp.Clear()
			// 👇 render the byte
			temp.SetHexColor(utils.Ternary(reverse, device.bgColor, color))
			str := string(utils.E2A([]byte{cell}))
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
	SendMessage(Message{bus: device.bus, eventType: "status", params: status})
}
