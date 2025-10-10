package glyph

import (
	"image"
)

type Cache struct {
	cache map[Glyph]image.Image
}

func (cache *Cache) Get(glyph Glyph) (image.Image, bool) {
	if cache.cache == nil {
		return nil, false
	}
	img, ok := cache.cache[glyph]
	return img, ok
}

func (cache *Cache) Set(glyph Glyph, img image.Image) {
	if cache.cache == nil {
		cache.cache = make(map[Glyph]image.Image)
	}
	cache.cache[glyph] = img
}
