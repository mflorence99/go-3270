package device

import (
	"emulator/utils"
	"image"
	"math"
	"time"

	"github.com/fogleman/gg"
)

// 🟧 Render the 3270 buffer into a graphics context

// 👁️ go3270.go for how pixels actually get drawn on the screen

type Glyph struct {
	u8         byte
	color      string
	reverse    bool
	underscore bool
}

type RenderBufferOpts struct {
	quiet   bool
	blinkOn bool
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

func (device *Device) RenderBuffer(opts RenderBufferOpts) {
	defer utils.ElapsedTime(time.Now(), "RenderBuffer", opts.quiet)
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
			"color":  "coral",
			"ebcdic": true,
			"title":  "RenderBuffer",
		}
		device.SendMessage(Message{eventType: "dumpBytes", params: params, u8s: device.buffer})
	}
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
			u8:         cell,
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
