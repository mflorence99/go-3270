package device

import (
	"image"

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

func (device *Device) Close() {
	if device.blinker != nil {
		close(device.blinker)
		device.blinker = nil
	}
}
