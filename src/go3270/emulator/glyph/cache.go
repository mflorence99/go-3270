package glyph

import (
	"go3270/emulator/pubsub"
	"go3270/emulator/screen"
	"go3270/emulator/utils"
	"image"

	"github.com/fogleman/gg"
)

type Cache struct {
	cache map[Glyph]image.Image
	cfg   pubsub.Config
}

func NewCache(cfg pubsub.Config) *Cache {
	c := new(Cache)
	c.cache = make(map[Glyph]image.Image)
	c.cfg = cfg
	return c
}

func (c *Cache) ImageFor(g Glyph) image.Image {
	img, ok := c.cache[g]
	if ok {
		return img
	}
	// 👇 cache miss: draw the glyph in a temporary context
	println("🔥 glyph cache miss", g.Char)
	box := screen.NewBox(0, 0, c.cfg)
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
	// 👇 now cache and bitblt the glyph
	c.cache[g] = temp.Image()
	return c.cache[g]
}
