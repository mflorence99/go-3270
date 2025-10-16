package glyph

import (
	"go3270/emulator/pubsub"
	"go3270/emulator/utils"
	"image"

	"github.com/fogleman/gg"
)

type Cache struct {
	bus   *pubsub.Bus
	cache map[Glyph]image.Image
	cfg   pubsub.Config
}

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

func (c *Cache) ImageFor(g Glyph, box pubsub.Box) image.Image {
	img, ok := c.cache[g]
	if !ok {
		// 👇 cache miss: draw the glyph in a temporary context
		// println("🔥 glyph cache miss", g.Char)
		rgba := image.NewRGBA(image.Rect(0, 0, int(box.W), int(box.H)))
		temp := gg.NewContextForRGBA(rgba)
		temp.SetFontFace(utils.Ternary(g.Highlight, *c.cfg.BoldFace, *c.cfg.NormalFace))
		// 👇 clear background
		temp.SetHexColor(utils.Ternary(g.Reverse, g.Color, c.cfg.BgColor))
		temp.Clear()
		// 👇 render the byte
		temp.SetHexColor(utils.Ternary(g.Reverse, c.cfg.BgColor, g.Color))
		temp.DrawString(string(g.Char), 0, box.Baseline-box.Y)
		if g.Underscore {
			temp.SetLineWidth(2)
			temp.MoveTo(0, box.H-1)
			temp.LineTo(box.W, box.H-1)
			temp.Stroke()
		}
		// 👇 now cache the glyph
		img = temp.Image()
		c.cache[g] = img
	}
	return img
}
