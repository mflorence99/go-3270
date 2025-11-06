package screen

import (
	"go3270/emulator/conv"
	"go3270/emulator/pubsub"
	"go3270/emulator/utils"
	"image"

	"github.com/fogleman/gg"
)

// 🟧 Cache of glyphs as drawn from the buffer

type Cache struct {
	bus   *pubsub.Bus
	cache map[Glyph]image.Image
	cfg   pubsub.Config
}

// 🟦 Constructor

func NewCache(bus *pubsub.Bus) *Cache {
	c := new(Cache)
	c.bus = bus
	// 👇 subscriptions
	c.bus.SubConfig(c.configure)
	return c
}

func (c *Cache) configure(cfg pubsub.Config) {
	c.cfg = cfg
	c.cache = make(map[Glyph]image.Image)
}

// 🟦 Public functions

func (c *Cache) ImageFor(g Glyph, box Box) image.Image {
	img, ok := c.cache[g]
	if !ok {
		// 👇 cache miss: draw the glyph in a temporary context
		rgba := image.NewRGBA(image.Rect(0, 0, int(box.W), int(box.H)))
		temp := gg.NewContextForRGBA(rgba)
		temp.SetFontFace(utils.Ternary(g.Highlight, *c.cfg.BoldFace, *c.cfg.NormalFace))
		// 👇 clear background
		temp.SetHexColor(utils.Ternary(g.Reverse, g.Color, c.cfg.BgColor))
		temp.DrawRectangle(0, 0, box.W, box.H)
		temp.Fill()
		// 👇 render the byte
		temp.SetHexColor(utils.Ternary(g.Reverse, c.cfg.BgColor, g.Color))
		temp.DrawString(string(conv.E2A(g.Char)), 0, box.Baseline-box.Y)
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
