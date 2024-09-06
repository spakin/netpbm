/*
Package npcolor implements the color types and models used by Netpbm.

RGBAM and RGBAM64 are analogous to the color package's RGBA and RGBA64
in that they represent, respectively, colors with 8 bits and 16 bits per
color channel.  They additionally store a maximum channel value.
Consequently, while a color.RGBA has a hard-wired upper bound of 255 per
color channel, npcolor.RGBAM supports any upper bound from 1–255.
Likewise, while a color.RGBA64 has a hard-wired upper bound of 65,535
per color channel, npcolor.RGBAM64 supports any upper bound from
1–65,535.

RGBM and RGBM64 are also analogous to the color package's RGBA and
RGBA64 in that they represent, respectively, colors with 8 bits and 16
bits per color channel.  While color.RGBA and color.RGBA64 store red,
green, blue, and alpha channels, npcolor.RGBM and npcolor.RGBM64 lack
alpha channels.  However, they store a maximum channel value.
Consequently, while a color.RGBA has a hard-wired upper bound of 255 per
color channel, npcolor.RGBM supports any upper bound from 1–255.
Likewise, while a color.RGBA64 has a hard-wired upper bound of 65,535
per color channel, npcolor.RGBM64 supports any upper bound from
1–65,535.

GrayM and GrayM32 are analogous to the color package's Gray and Gray16
in that they represent, respectively, 8-bit and 16-bit grayscale
values.  However, while a color.Gray value has a hard-wired upper
bound of 255, npcolor.GrayM supports any upper bound from 1–255.
Likewise, while a color.Gray16 value has a hard-wired upper bound of
65,535, npcolor.GrayM32 supports any upper bound from 1–65,535.

GrayAM and GrayAM32 have no analogue in the color package.  They
represent, respectively, 8-bit and 16-bit grayscale values with an
alpha channel.  Like the other colors represented in this package,
these support variable maximum channel values.  npcolor.GrayAM
supports any upper bound from 1–255, and npcolor.GrayM32 supports any
upper bound from 1–65,535.
*/
package npcolor

import (
	"image/color"
)

// GrayM represents an 8-bit grayscale value and the value to represent 100%
// white.  Because GrayM does not support alpha channels it does make sense to
// describe it as either "alpha-premultiplied" or "non-alpha-premultiplied".
type GrayM struct {
	Y, M uint8
}

// RGBA converts a GrayM to alpha-premultiplied R, G, B, and A.
func (c GrayM) RGBA() (r, g, b, a uint32) {
	if c.M == 0 {
		return
	}
	m := uint32(c.M)
	y := (uint32(c.Y)*0xffff + m/2) / m
	return y, y, y, 0xffff
}

// A GrayMModel represents the maximum value of a GrayM (0-255).
type GrayMModel struct {
	M uint8 // Maximum value of the luminance channel
}

// Convert converts an arbitrary color to a GrayM.
func (model GrayMModel) Convert(c color.Color) color.Color {
	if gray, ok := c.(GrayM); ok && gray.M == model.M {
		return c
	}
	r, g, b, _ := c.RGBA()
	y := (299*r + 587*g + 114*b + 500) / 1000
	m := uint32(model.M)
	y = (y*m + 0xffff/2) / 0xffff
	return GrayM{Y: uint8(y), M: uint8(m)}
}

// GrayM32 represents a 16-bit grayscale value and the value to represent 100%
// white.  Because GrayM16 does not support alpha channels it does make sense
// to describe it as either "alpha-premultiplied" or "non-alpha-premultiplied".
type GrayM32 struct {
	Y, M uint16
}

// RGBA converts a GrayM32 to alpha-premultiplied R, G, B, and A.
func (c GrayM32) RGBA() (r, g, b, a uint32) {
	if c.M == 0 {
		return
	}
	m := uint32(c.M)
	y := (uint32(c.Y)*0xffff + m/2) / m
	return y, y, y, 0xffff
}

// A GrayM32Model represents the maximum value of a GrayM32 (0-65535).
type GrayM32Model struct {
	M uint16 // Maximum value of the luminance channel
}

// Convert converts an arbitrary color to a GrayM32.
func (model GrayM32Model) Convert(c color.Color) color.Color {
	if gray, ok := c.(GrayM32); ok && gray.M == model.M {
		return c
	}
	r, g, b, _ := c.RGBA()
	y := (299*r + 587*g + 114*b + 500) / 1000
	m := uint32(model.M)
	y = (y*m + 0xffff/2) / 0xffff
	return GrayM32{Y: uint16(y), M: uint16(m)}
}

// RGBM represents a 24-bit color and the value used for 100% of a color
// channel.  Because RGBM does not support alpha channels it does not make
// sense to describe it as either "alpha-premultiplied" or
// "non-alpha-premultiplied".
type RGBM struct {
	R, G, B, M uint8
}

// RGBA converts an RGBM to alpha-premultiplied R, G, B, and A.
func (c RGBM) RGBA() (r, g, b, a uint32) {
	if c.M == 0 {
		return
	}
	m := uint32(c.M)
	r = (uint32(c.R)*0xffff + m/2) / m
	g = (uint32(c.G)*0xffff + m/2) / m
	b = (uint32(c.B)*0xffff + m/2) / m
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
	const half = 0xffff / 2
	r = (r*m + half) / 0xffff
	g = (g*m + half) / 0xffff
	b = (b*m + half) / 0xffff
	return RGBM{R: uint8(r), G: uint8(g), B: uint8(b), M: uint8(m)}
}

// RGBM64 represents a 48-bit color and the value used for 100% of a color
// channel.
type RGBM64 struct {
	R, G, B, M uint16
}

// RGBA converts an RGBM64 to alpha-premultiplied R, G, B, and A.
func (c RGBM64) RGBA() (r, g, b, a uint32) {
	if c.M == 0 {
		return
	}
	m := uint32(c.M)
	r = (uint32(c.R)*0xffff + m/2) / m
	g = (uint32(c.G)*0xffff + m/2) / m
	b = (uint32(c.B)*0xffff + m/2) / m
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
	const half = 0xffff / 2
	r = (r*m + half) / 0xffff
	g = (g*m + half) / 0xffff
	b = (b*m + half) / 0xffff
	return RGBM64{R: uint16(r), G: uint16(g), B: uint16(b), M: uint16(m)}
}

// GrayAM represents an 8-bit grayscale value with an alpha channel and the
// value to represent 100% white.
type GrayAM struct {
	Y, A, M uint8
}

// RGBA converts a GrayAM to alpha-premultiplied R, G, B, and A.
func (c GrayAM) RGBA() (r, g, b, a uint32) {
	if c.M == 0 {
		return
	}
	m := uint32(c.M)
	a = (uint32(c.A)*0xffff + m/2) / m
	y := (uint32(c.Y)*a + m/2) / m
	return y, y, y, a
}

// A GrayAMModel represents the maximum value of a GrayAM (0-255).
type GrayAMModel struct {
	M uint8 // Maximum value of the luminance channel
}

// Convert converts an arbitrary color to a GrayAM.
func (model GrayAMModel) Convert(c color.Color) color.Color {
	if gray, ok := c.(GrayAM); ok && gray.M == model.M {
		return c
	}
	r, g, b, a := c.RGBA()
	y := (299*r + 587*g + 114*b + 500) / 1000
	m := uint32(model.M)
	const half = 0xffff / 2
	y = (y*m + half) / 0xffff
	a = (a*m + half) / 0xffff
	return GrayAM{Y: uint8(y), A: uint8(a), M: uint8(m)}
}

// GrayAM48 represents a 16-bit grayscale value with an alpha channel and the
// value to represent 100% white.
type GrayAM48 struct {
	Y, A, M uint16
}

// RGBA converts a GrayAM48 to alpha-premultiplied R, G, B, and A.
func (c GrayAM48) RGBA() (r, g, b, a uint32) {
	if c.M == 0 {
		return
	}
	m := uint32(c.M)
	a = (uint32(c.A)*0xffff + m/2) / m
	y := (uint32(c.Y)*a + m/2) / m
	return y, y, y, a
}

// A GrayAM48Model represents the maximum value of a GrayAM48 (0-65535).
type GrayAM48Model struct {
	M uint16 // Maximum value of the luminance channel
}

// Convert converts an arbitrary color to a GrayAM48.
func (model GrayAM48Model) Convert(c color.Color) color.Color {
	if gray, ok := c.(GrayAM48); ok && gray.M == model.M {
		return c
	}
	r, g, b, a := c.RGBA()
	y := (299*r + 587*g + 114*b + 500) / 1000
	m := uint32(model.M)
	const half = 0xffff / 2
	y = (y*m + half) / 0xffff
	a = (a*m + half) / 0xffff
	return GrayAM48{Y: uint16(y), A: uint16(a), M: uint16(m)}
}

// RGBAM represents a 32-bit color and the value used for 100% of a color
// channel.
type RGBAM struct {
	R, G, B, A, M uint8
}

// RGBA converts an RGBAM to alpha-premultiplied R, G, B, and A.
func (c RGBAM) RGBA() (r, g, b, a uint32) {
	if c.M == 0 {
		return
	}
	m := uint32(c.M)
	a = (uint32(c.A)*0xffff + m/2) / m
	r = (uint32(c.R)*a + m/2) / m
	g = (uint32(c.G)*a + m/2) / m
	b = (uint32(c.B)*a + m/2) / m
	return
}

// An RGBAMModel represents the maximum value of an RGBM (0-255).
type RGBAMModel struct {
	M uint8 // Maximum value of each color channel
}

// Convert converts an arbitrary color to an RGBM.
func (model RGBAMModel) Convert(c color.Color) color.Color {
	if rgba, ok := c.(RGBAM); ok && rgba.M == model.M {
		return c
	}
	m := uint32(model.M)
	r, g, b, a := c.RGBA()
	const half = 0xffff / 2
	r = (r*m + half) / 0xffff
	g = (g*m + half) / 0xffff
	b = (b*m + half) / 0xffff
	a = (a*m + half) / 0xffff
	return RGBAM{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a), M: uint8(m)}
}

// RGBAM64 represents a 48-bit color and the value used for 100% of a color
// channel.
type RGBAM64 struct {
	R, G, B, A, M uint16
}

// RGBA converts an RGBAM64 to alpha-premultiplied R, G, B, and A.
func (c RGBAM64) RGBA() (r, g, b, a uint32) {
	if c.M == 0 {
		return
	}
	m := uint32(c.M)
	a = (uint32(c.A)*0xffff + m/2) / m
	r = (uint32(c.R)*a + m/2) / m
	g = (uint32(c.G)*a + m/2) / m
	b = (uint32(c.B)*a + m/2) / m
	return
}

// An RGBAM64Model represents the maximum value of an RGBAM64 (0-65535).
type RGBAM64Model struct {
	M uint16 // Maximum value of each color channel
}

// Convert converts an arbitrary color to an RGBAM64.
func (model RGBAM64Model) Convert(c color.Color) color.Color {
	if rgba, ok := c.(RGBAM64); ok && rgba.M == model.M {
		return c
	}
	m := uint32(model.M)
	r, g, b, a := c.RGBA()
	const half = 0xffff / 2
	r = (r*m + half) / 0xffff
	g = (g*m + half) / 0xffff
	b = (b*m + half) / 0xffff
	a = (a*m + half) / 0xffff
	return RGBAM64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a), M: uint16(m)}
}
