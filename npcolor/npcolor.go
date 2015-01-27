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

// An RGBMModel represents the maximum value of an RGBM (0-255).
type RGBMModel struct {
	M uint8 // Maximum value of each color channel
}

// Convert converts an arbitrary color to an RGBM.
func (model RGBMModel) Convert(c color.Color) color.Color {
	if rgb, ok := c.(RGBM); ok && rgb.M == model.M {
		return c
	}
	m := uint32(model.M)
	r, g, b, _ := c.RGBA()
	r = (r * m) / 0xffff
	g = (g * m) / 0xffff
	b = (b * m) / 0xffff
	return RGBM{R: uint8(r), G: uint8(g), B: uint8(b), M: uint8(m)}
}

// RGBM64 represents a 48-bit color and the value used for 100% of a color
// channel.  Because RGBM64 does not support alpha channels it does make sense
// to describe it as either "alpha-premultiplied" or "non-alpha-premultiplied".
type RGBM64 struct {
	R, G, B, M uint16
}

// RGBA converts an RGBM64 to alpha-premultiplied R, G, B, and A.
func (c RGBM64) RGBA() (r, g, b, a uint32) {
	m := uint32(c.M)
	r = (uint32(c.R) * 0xffff) / m
	g = (uint32(c.G) * 0xffff) / m
	b = (uint32(c.B) * 0xffff) / m
	a = 0xffff
	return
}

// An RGBM64Model represents the maximum value of an RGBM64 (0-65535).
type RGBM64Model struct {
	M uint16 // Maximum value of each color channel
}

// Convert converts an arbitrary color to an RGBM64.
func (model RGBM64Model) Convert(c color.Color) color.Color {
	if rgb, ok := c.(RGBM64); ok && rgb.M == model.M {
		return c
	}
	m := uint32(model.M)
	r, g, b, _ := c.RGBA()
	r = (r * m) / 0xffff
	g = (g * m) / 0xffff
	b = (b * m) / 0xffff
	return RGBM64{R: uint16(r), G: uint16(g), B: uint16(b), M: uint16(m)}
}
