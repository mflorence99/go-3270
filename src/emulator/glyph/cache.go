package glyph

import (
	"image"
)

type Cache struct {
	cache map[Glyph]image.Image
}

func NewCache() *Cache {
	c := new(Cache)
	c.cache = make(map[Glyph]image.Image)
	return c
}

func (c *Cache) Get(glyph Glyph) (image.Image, bool) {
	img, ok := c.cache[glyph]
	return img, ok
}

func (c *Cache) Set(glyph Glyph, img image.Image) {
	c.cache[glyph] = img
}
