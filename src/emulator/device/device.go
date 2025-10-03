package device

import (
	"emulator/types"
	"emulator/utils"
	"fmt"
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
	device.initializeBuffer()
	return device
}

func (device *Device) Close() {
	fmt.Println("device.Close()")
}

func (device *Device) MessageUI(eventType string, bytes []uint8, params map[string]any, args ...any) {
	device.bus.Publish("go3270", eventType, bytes, params, args)
}

func (device *Device) ReceiveFromApp(bytes []uint8) {
	// 👇 reset changes stack
	device.changes = utils.NewStack[int](device.size)
	// 👇 data can be split into multiple frames
	frames := device.MakeFramesFromBytes(bytes)
	for ix := range frames {
		fmt.Printf("device.ReceiveFromApp(frame #%d)\n", ix)
		// 👇 extract command
		out := frames[ix]
		cmd, err := out.Next()
		if err != nil {
			panic(fmt.Sprintf("unable to extraact write command: %s", err.Error()))
		}
		device.command = cmd
		fmt.Printf("COMMAND=%s\n", types.Command[device.command])
		// 👇 for all but WSF, extract WCC
		if device.command != types.CommandLookup["WSF"] {
			u8, err := out.Next()
			if err != nil {
				panic(fmt.Sprintf("unable to extract WCC: %s", err.Error()))
			}
			device.wcc = NewWCC(u8)
			fmt.Println(device.wcc.ToString())
		}
		// 👇 dispatch on command
		device.writeCommands(out)
	}
	// 👇 now we can render the buffer to the drawing context
	device.renderBuffer()
}

// 👇 Helpers - they need to be public to be tested

func (device *Device) MakeFramesFromBytes(bytes []uint8) []*OutboundDataStream {
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

// 👇 helpers

func (device *Device) boundingBox(addr int) (float64, float64, float64, float64, float64) {
	col := addr % device.cols
	row := int(addr / device.cols)
	w := device.fontWidth * device.paddedWidth
	h := device.fontHeight * device.paddedHeight
	x := float64(col) * w
	y := float64(row) * h
	// 🔥 we could do better calculating the baseline - this is just a WAG, because an em is drawn with a significantly different height than that returned by MeasureString()
	baseline := y + h - (device.fontSize / 3 * device.scaleFactor)
	return x, y, w, h, baseline
}

func (device *Device) initializeBuffer() {
	device.addr = 0
	device.attrs = make([]*Attributes, device.size)
	device.buffer = make([]uint8, device.size)
	device.cursorAt = 0
	device.erase = true
}

func (device *Device) putBuffer(byte uint8) {
	device.buffer[device.addr] = byte
	device.changes.Push(device.addr)
	device.addr += 1
	// 👇 note wrap around
	if device.addr > device.size {
		device.addr = 0
	}
}

func (device *Device) renderBuffer() {
	defer utils.ElapsedTime(time.Now(), "renderBuffer")
	// 👇 ragged fonts if drawn on transparent!
	if device.erase {
		device.gg.SetHexColor(device.bgColor)
		device.gg.Clear()
		device.erase = false
	}
	for !device.changes.IsEmpty() {
		addr := device.changes.Pop()
		// 🔥 TEMPORARY
		x, _, _, _, baseline := device.boundingBox(addr)
		device.gg.SetHexColor(device.color)
		str := string(device.buffer[addr])
		device.gg.DrawString(str, x, baseline)

	}
}

func (device *Device) writeCommands(out *OutboundDataStream) {
	// 👇 dispatch on command
	switch device.command {
	case types.CommandLookup["RMA"]:
	case types.CommandLookup["EAU"]:
	case types.CommandLookup["EWA"]:
	case types.CommandLookup["W"]:
	case types.CommandLookup["RB"]:
	case types.CommandLookup["WSF"]:
	case types.CommandLookup["EW"]:
		device.initializeBuffer()
		device.writeOrdersAndData(out)
	case types.CommandLookup["RM"]:
	}
}

func (device *Device) writeOrdersAndData(out *OutboundDataStream) {
	for out.HasNext() {
		// 👇 look at each byte to see if it is an order
		order, err := out.Next()
		if err != nil {
			panic(fmt.Sprintf("unexpected EOF: %s", err.Error()))
		}
		// 👇 dispatch on order
		switch order {
		case types.OrderLookup["PT"]:
		case types.OrderLookup["GE"]:
		case types.OrderLookup["SBA"]:
		case types.OrderLookup["EUA"]:
		case types.OrderLookup["IC"]:
		case types.OrderLookup["SF"]:
		case types.OrderLookup["SA"]:
		case types.OrderLookup["SFE"]:
		case types.OrderLookup["MF"]:
		case types.OrderLookup["RA"]:
		// 👇 if it isn't an order, it's data
		default:
			if order == 0x00 || order >= 0x40 {
				device.putBuffer(utils.E2A([]uint8{order})[0])
			}
		}
	}
}
