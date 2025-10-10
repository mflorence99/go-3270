package glyph

import (
	"image"
)

type Cache struct {
	cache map[Glyph]image.Image
}

func (c *Cache) Get(glyph Glyph) (image.Image, bool) {
	if c.cache == nil {
		return nil, false
	}
	img, ok := c.cache[glyph]
	return img, ok
}

func (c *Cache) Set(glyph Glyph, img image.Image) {
	if c.cache == nil {
		c.cache = make(map[Glyph]image.Image)
	}
	c.cache[glyph] = img
}
