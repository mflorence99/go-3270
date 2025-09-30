package device

import (
	"emulator/types"
	"emulator/utils"
	"math"
	"math/rand"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/fogleman/gg"
)

type Device struct {
	bus          EventBus.Bus
	color        string
	cols         float64
	gg           *gg.Context
	fontHeight   float64
	fontSize     float64
	fontWidth    float64
	paddedHeight float64
	paddedWidth  float64
	rows         float64
	scaleFactor  float64
}

func NewDevice(bus EventBus.Bus,
	color string,
	cols float64,
	gg *gg.Context,
	fontHeight float64,
	fontSize float64,
	fontWidth float64,
	paddedHeight float64,
	paddedWidth float64,
	rows float64,
	scaleFactor float64) *Device {
	device := &Device{}
	device.bus = bus
	device.color = color
	device.cols = cols
	device.gg = gg
	device.fontHeight = fontHeight
	device.fontSize = fontSize
	device.fontWidth = fontWidth
	device.paddedHeight = paddedHeight
	device.paddedWidth = paddedWidth
	device.rows = rows
	device.scaleFactor = scaleFactor
	return device
}

func (device *Device) Close() {
	device.bus.Publish("go3270-log", "%cDevice closing", "color: cadetblue")
}

func (device *Device) ReceiveFromApp(bytes []uint8) {
}

// ///////////////////////////////////////////////////////////////////////////
// 🔥 EVERYTHING BELOW HERE IS JUST TEST CODE
// ///////////////////////////////////////////////////////////////////////////

func (device *Device) TestPattern() {
	defer utils.ElapsedTime(time.Now(), "TestPattern")
	device.gg.SetHexColor(types.CLUT[0xf0][0]) /* 👈 ragged fonts if draw on transparent! */
	device.gg.Clear()
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[{]};:,<.>/?"
	chs := []rune(str)
	for col := 0.0; col < device.cols; col++ {
		for row := 0.0; row < device.rows; row++ {
			x, _, _, _, baseline := device.boundingBox(col, row)
			// 👇 choose colors from the CLUT, using the base color if out of range
			ix := uint8(math.Floor(col/10) + 0xf1)
			bright := device.color
			color := device.color
			if ix <= 0xf7 {
				bright = types.CLUT[ix][0]
				color = types.CLUT[ix][1]
			}
			// 👇 alternate high intensity, normal
			if int(row)%2 == 0 {
				device.gg.SetHexColor(bright)
			} else {
				device.gg.SetHexColor(color)
			}
			ich := rand.Intn(len(chs))
			ch := string(chs[ich])
			device.gg.DrawString(ch, x, baseline)
		}
	}
	device.bus.Publish("go3270-render")
}

// 👇 Helpers

func (device *Device) boundingBox(col, row float64) (float64, float64, float64, float64, float64) {
	w := device.fontWidth * device.paddedWidth
	h := device.fontHeight * device.paddedHeight
	x := col * w
	y := row * h
	// 🔥 we could do better calculating the baseline - this is just a WAG, because an em is drawn with a significantly different height than that returned by MeasureString()
	baseline := y + h - (device.fontSize / 3 * device.scaleFactor)
	return x, y, w, h, baseline
}
