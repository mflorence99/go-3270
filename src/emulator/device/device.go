package device

import (
	"emulator/types"
	"emulator/utils"
	"fmt"

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
	cols         float64
	fontHeight   float64
	fontSize     float64
	fontWidth    float64
	paddedHeight float64
	paddedWidth  float64
	rows         float64
	scaleFactor  float64

	// 👇 model the device buffer
	addr     int
	attrs    []*Attributes
	buffer   []uint8
	changed  []bool
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
	cols float64,
	rows float64,
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
	frames := device.MakeFramesFromBytes(bytes)
	for ix := range frames {
		fmt.Printf("device.ReceiveFromApp(frame #%d)\n", ix)
		// 👇 extract command
		out := frames[ix]
		device.command, _ = out.Next()
		fmt.Printf("COMMAND=%s\n", types.Command[device.command])
		// 👇 for all but WSF, extract WCC
		if device.command != types.CommandLookup["WSF"] {
			u8, _ := out.Next()
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

func (device *Device) initializeBuffer() {
	size := int(device.cols * device.rows)
	device.addr = 0
	device.attrs = make([]*Attributes, size)
	device.buffer = make([]uint8, size)
	device.changed = make([]bool, size)
	device.cursorAt = 0
	device.erase = true
}

func (device *Device) renderBuffer() {
	if device.erase {
		// 👈 ragged fonts if draw on transparent!
		device.gg.SetHexColor(device.bgColor)
		device.gg.Clear()
		device.erase = false
	}
}

func (device *Device) writeCommands(out *OutboundDataStream) {
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
		order, _ := out.Next()
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
		default:
			// 👇 if it isn't an order, it's data
			if order == 0x00 || order >= 0x40 {
				device.buffer[device.addr] = utils.E2A([]uint8{order})[0]
				device.changed[device.addr] = true
			}
		}
	}
}
