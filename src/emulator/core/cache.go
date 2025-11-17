package core

import (
	"emulator/conv"
	"emulator/utils"
	"image"

	"github.com/fogleman/gg"
)

// 🟧 Cache of glyphs as drawn from the buffer

type Cache struct {
	cache map[Glyph]image.Image

	emu *Emulator // 👈 back pointer to all common components
}

// 🟦 Constructor

func NewCache(emu *Emulator) *Cache {
	c := new(Cache)
	c.emu = emu
	// 👇 subscriptions
	c.emu.Bus.SubInit(c.init)
	// 🔥 we never reset the glyph cache!
	// c.emu.Bus.SubReset(c.reset)
	return c
}

func (c *Cache) init() {
	c.cache = make(map[Glyph]image.Image)
}

// 🟦 Public functions

func (c *Cache) ImageFor(g Glyph, box Box) image.Image {
	img, ok := c.cache[g]
	if !ok {
		// 👇 cache miss: draw the glyph in a temporary context
		rgba := image.NewRGBA(image.Rect(0, 0, int(box.W), int(box.H)))
		temp := gg.NewContextForRGBA(rgba)
		temp.SetFontFace(utils.Ternary(g.Highlight, *c.emu.Cfg.BoldFace, *c.emu.Cfg.NormalFace))
		// 👇 clear background
		temp.SetHexColor(utils.Ternary(g.Reverse, g.Color, c.emu.Cfg.BgColor))
		temp.DrawRectangle(0, 0, box.W, box.H)
		temp.Fill()
		// 👇 render the byte
		temp.SetHexColor(utils.Ternary(g.Reverse, c.emu.Cfg.BgColor, g.Color))
		temp.DrawString(string(conv.E2Rune(g.LCID, g.Char)), 0, box.Baseline-box.Y)
		// 👇 lines for outline/underscore
		if g.Underscore || g.Outline.Bottom {
			temp.SetLineWidth(1)
			temp.DrawLine(0, box.H, box.W, box.H)
			temp.Stroke()
		}
		if g.Outline.Right {
			temp.SetLineWidth(1)
			temp.DrawLine(box.W, 0, box.W, box.H)
			temp.Stroke()
		}
		if g.Outline.Top {
			temp.SetLineWidth(1)
			temp.DrawLine(0, 0, box.W, 0)
			temp.Stroke()
		}
		if g.Outline.Left {
			temp.SetLineWidth(1)
			temp.DrawLine(0, 0, 0, box.H)
			temp.Stroke()
		}
		// 👇 now cache the glyph
		img = temp.Image()
		c.cache[g] = img
	}
	return img
}
