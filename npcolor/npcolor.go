// Package npcolor implements the color types and models used by Netpbm.
package npcolor

import (
	"image/color"
)

// RGBM represents a 24-bit color and the value used for 100% of a color
// channel.  Because RGBM does not support alpha channels it does make sense to
// describe it as either "alpha-premultiplied" or "non-alpha-premultiplied".
type RGBM struct {
	R, G, B, M uint8
}

// RGBA converts an RGBM to alpha-premultiplied R, G, B, and A.
func (c RGBM) RGBA() (r, g, b, a uint32) {
	m := uint32(c.M)
	r = (uint32(c.R) * 0xffff) / m
	g = (uint32(c.G) * 0xffff) / m
	b = (uint32(c.B) * 0xffff) / m
	a = 0xffff
	return
}

// rgbmModel converts an arbitrary color to an RGBM with maximum value 255.
func rgbmModel(c color.Color) color.Color {
	if _, ok := c.(RGBM); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return RGBM{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), M: 0xff}
}

// An RGBMModel converts an arbitrary color to an RGBM with maximum value 255.
var RGBMModel color.Model = color.ModelFunc(rgbmModel)

// RGBM64 represents a 48-bit color and the value used for 100% of a color
// channel.  Because RGBM64 does not support alpha channels it does make sense
// to describe it as either "alpha-premultiplied" or "non-alpha-premultiplied".
type RGBM64 struct {
	R, G, B, M uint16
}

// RGBA converts an RGBM16 to alpha-premultiplied R, G, B, and A.
func (c RGBM64) RGBA() (r, g, b, a uint32) {
	m := uint32(c.M)
	r = (uint32(c.R) * 0xffff) / m
	g = (uint32(c.G) * 0xffff) / m
	b = (uint32(c.B) * 0xffff) / m
	a = 0xffff
	return
}

// rgbm64Model converts an arbitrary color to an RGBM64 with maximum value
// 65535.
func rgbm64Model(c color.Color) color.Color {
	if _, ok := c.(RGBM64); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return RGBM64{R: uint16(r), G: uint16(g), B: uint16(b), M: 0xffff}
}

// An RGBM64Model converts an arbitrary color to an RGBM64 with maximum value
// 65535.
var RGBM64Model color.Model = color.ModelFunc(rgbm64Model)
