package netpbm

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"io"
	"unicode"

	"github.com/spakin/netpbm/npcolor"
)

// An RGBAM is an in-memory image whose At method returns npcolor.RGBAM values.
type RGBAM struct {
	// Pix holds the image's pixels, in R, G, B (no M) order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*3].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent
	// pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
	// Model is the image's color model.
	Model npcolor.RGBAMModel
}

// ColorModel returns the RGBAM image's color model.
func (p *RGBAM) ColorModel() color.Model { return p.Model }

// Bounds returns the domain for which At can return non-zero color.  The
// bounds do not necessarily contain the point (0, 0).
func (p *RGBAM) Bounds() image.Rectangle { return p.Rect }

// At returns the color of the pixel at (x, y) as a color.Color.
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (p *RGBAM) At(x, y int) color.Color {
	return p.RGBAMAt(x, y)
}

// RGBAMAt returns the color of the pixel at (x, y) as an npcolor.RGBAM.
func (p *RGBAM) RGBAMAt(x, y int) npcolor.RGBAM {
	if !(image.Point{x, y}.In(p.Rect)) {
		return npcolor.RGBAM{}
	}
	i := p.PixOffset(x, y)
	return npcolor.RGBAM{
		R: p.Pix[i+0],
		G: p.Pix[i+1],
		B: p.Pix[i+2],
		A: p.Pix[i+3],
		M: p.Model.M,
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGBAM) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

// Set sets the pixel at (x, y) to a given color, expressed as a color.Color.
func (p *RGBAM) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := p.Model.Convert(c).(npcolor.RGBAM)
	p.Pix[i+0] = c1.R
	p.Pix[i+1] = c1.G
	p.Pix[i+2] = c1.B
	p.Pix[i+3] = c1.A
}

// SetRGBAM sets the pixel at (x, y) to a given color, expressed as an
// npcolor.RGBAM.
func (p *RGBAM) SetRGBAM(x, y int, c npcolor.RGBAM) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	if c.M == p.Model.M {
		p.Pix[i+0] = c.R
		p.Pix[i+1] = c.G
		p.Pix[i+2] = c.B
		p.Pix[i+3] = c.A
	} else {
		p.Set(x, y, c)
	}
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *RGBAM) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &RGBAM{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGBAM{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *RGBAM) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}
	i0, i1 := 3, p.Rect.Dx()*4
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		for i := i0; i < i1; i += 4 {
			if p.Pix[i] != 0xff {
				return false
			}
		}
		i0 += p.Stride
		i1 += p.Stride
	}
	return true
}

// MaxValue returns the maximum value allowed on any color channel.
func (p *RGBAM) MaxValue() uint16 {
	return uint16(p.Model.M)
}

// Format identifies the image as a PPM image.
func (p *RGBAM) Format() Format {
	return PAM
}

// NewRGBAM returns a new RGBAM with the given bounds and maximum channel value.
func NewRGBAM(r image.Rectangle, m uint8) *RGBAM {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 4*w*h)
	model := npcolor.RGBAMModel{M: m}
	return &RGBAM{pix, 4 * w, r, model}
}

// An RGBAM64 is an in-memory image whose At method returns npcolor.RGBAM64
// values.
type RGBAM64 struct {
	// Pix holds the image's pixels, in R, G, B, M order and big-endian
	// format. The pixel at (x, y) starts at Pix[(y-Rect.Min.Y)*Stride +
	// (x-Rect.Min.X)*8].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent
	// pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
	// Model is the image's color model.
	Model npcolor.RGBAM64Model
}

// ColorModel returns the RGBAM64 image's color model.
func (p *RGBAM64) ColorModel() color.Model { return p.Model }

// Bounds returns the domain for which At can return non-zero color.  The
// bounds do not necessarily contain the point (0, 0).
func (p *RGBAM64) Bounds() image.Rectangle { return p.Rect }

// At returns the color of the pixel at (x, y) as a color.Color.
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (p *RGBAM64) At(x, y int) color.Color {
	return p.RGBAM64At(x, y)
}

// RGBAM64At returns the color of the pixel at (x, y) as an npcolor.RGBAM64.
func (p *RGBAM64) RGBAM64At(x, y int) npcolor.RGBAM64 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return npcolor.RGBAM64{}
	}
	i := p.PixOffset(x, y)
	return npcolor.RGBAM64{
		R: uint16(p.Pix[i+0])<<8 | uint16(p.Pix[i+1]),
		G: uint16(p.Pix[i+2])<<8 | uint16(p.Pix[i+3]),
		B: uint16(p.Pix[i+4])<<8 | uint16(p.Pix[i+5]),
		A: uint16(p.Pix[i+6])<<8 | uint16(p.Pix[i+7]),
		M: p.Model.M,
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGBAM64) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*8
}

// Set sets the pixel at (x, y) to a given color, expressed as a color.Color.
func (p *RGBAM64) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := p.Model.Convert(c).(npcolor.RGBAM64)
	p.Pix[i+0] = uint8(c1.R >> 8)
	p.Pix[i+1] = uint8(c1.R)
	p.Pix[i+2] = uint8(c1.G >> 8)
	p.Pix[i+3] = uint8(c1.G)
	p.Pix[i+4] = uint8(c1.B >> 8)
	p.Pix[i+5] = uint8(c1.B)
	p.Pix[i+6] = uint8(c1.A >> 8)
	p.Pix[i+7] = uint8(c1.A)
}

// SetRGBAM64 sets the pixel at (x, y) to a given color, expressed as an
// npcolor.RGBAM.
func (p *RGBAM64) SetRGBAM64(x, y int, c npcolor.RGBAM64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	if c.M == p.Model.M {
		p.Pix[i+0] = uint8(c.R >> 8)
		p.Pix[i+1] = uint8(c.R)
		p.Pix[i+2] = uint8(c.G >> 8)
		p.Pix[i+3] = uint8(c.G)
		p.Pix[i+4] = uint8(c.B >> 8)
		p.Pix[i+5] = uint8(c.B)
		p.Pix[i+6] = uint8(c.A >> 8)
		p.Pix[i+7] = uint8(c.A)
	} else {
		p.Set(x, y, c)
	}
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *RGBAM64) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &RGBAM64{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGBAM64{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *RGBAM64) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}
	i0, i1 := 6, p.Rect.Dx()*8
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		for i := i0; i < i1; i += 8 {
			if p.Pix[i+0] != 0xff || p.Pix[i+1] != 0xff {
				return false
			}
		}
		i0 += p.Stride
		i1 += p.Stride
	}
	return true
}

// MaxValue returns the maximum value allowed on any color channel.
func (p *RGBAM64) MaxValue() uint16 {
	return uint16(p.Model.M)
}

// Format identifies the image as a PPM image.
func (p *RGBAM64) Format() Format {
	return PAM
}

// NewRGBAM64 returns a new RGBAM64 with the given bounds and maximum channel
// value.
func NewRGBAM64(r image.Rectangle, m uint16) *RGBAM64 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 8*w*h)
	model := npcolor.RGBAM64Model{M: m}
	return &RGBAM64{pix, 8 * w, r, model}
}

// GetPamHeader parses the entire header of a PAM file (raw or
// plain) and returns it as a netpbmHeader (plus a success value).
func (nr *netpbmReader) GetPamHeader() (netpbmHeader, bool) {
	var header netpbmHeader

	// Read the magic value and skip the following whitespace.
	rune1 := nr.GetNextByteAsRune()
	if rune1 != 'P' {
		return netpbmHeader{}, false
	}
	rune2 := nr.GetNextByteAsRune()
	if rune2 != '7' {
		return netpbmHeader{}, false
	}
	if !unicode.IsSpace(nr.GetNextByteAsRune()) {
		return netpbmHeader{}, false
	}
	header.Magic = string(rune1) + string(rune2)

	// Read the width and height.
	header.Width = nr.GetNextInt()
	header.Height = nr.GetNextInt()
	header.Depth = nr.GetNextInt()
	header.Maxval = nr.GetNextInt()

	if fieldName := nr.GetNextString(); fieldName != "TUPLTYPE" {
		return netpbmHeader{}, false
	}
	header.TupleType = nr.GetNextString()

	if fieldName := nr.GetNextString(); fieldName != "ENDHDR" {
		return netpbmHeader{}, false
	}

	if nr.Err() != nil || !unicode.IsSpace(nr.GetNextByteAsRune()) ||
		header.Maxval < 1 || header.Maxval > 65535 {
		return netpbmHeader{}, false
	}

	// Return the header and a success code.
	return header, true
}

// decodeConfigPAM reads and parses a PAM header, either "raw" (binary) or
// "plain" (ASCII).
func decodeConfigPAM(r io.Reader) (image.Config, error) {
	// We really want a bufio.Reader.  If we were given one, use it.  If
	// not, create a new one.
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	nr := newNetpbmReader(br)

	// Parse the PAM header.
	header, ok := nr.GetPamHeader()
	if !ok {
		err := nr.Err()
		if err == nil {
			err = errors.New("Invalid PAM header")
		}
		return image.Config{}, err
	}

	// Store and return the image configuration.
	var cfg image.Config
	cfg.Width = header.Width
	cfg.Height = header.Height
	if header.Maxval < 256 {
		cfg.ColorModel = npcolor.RGBAMModel{M: uint8(header.Maxval)}
	} else {
		cfg.ColorModel = npcolor.RGBAM64Model{M: uint16(header.Maxval)}
	}
	return cfg, nil
}

// decodePAM reads a complete "raw" (binary) PAM image.
func decodePAM(r io.Reader) (image.Image, error) {
	// Read the image header, and use it to prepare a color image.
	br := bufio.NewReader(r)
	config, err := decodeConfigPAM(br)
	if err != nil {
		return nil, err
	}

	// Create either a Color or a Color64 image.
	var img image.Image // Image to return
	var data []uint8    // RGB (no M) image data
	var maxVal uint     // 100% white value
	switch model := config.ColorModel.(type) {
	case npcolor.RGBAMModel:
		maxVal = uint(model.M)
		rgb := NewRGBAM(image.Rect(0, 0, config.Width, config.Height), uint8(maxVal))
		data = rgb.Pix
		img = rgb
	case npcolor.RGBM64Model:
		maxVal = uint(model.M)
		rgb := NewRGBAM64(image.Rect(0, 0, config.Width, config.Height), uint16(maxVal))
		data = rgb.Pix
		img = rgb
	default:
		panic("Unexpected color model")
	}

	// Raw PAM images are nice because we can read directly into the image
	// data.
	for len(data) > 0 {
		nRead, err := br.Read(data)
		if err != nil && err != io.EOF {
			return img, err
		}
		if nRead == 0 {
			return img, errors.New("Failed to read binary PPM data")
		}
		data = data[nRead:]
	}
	return img, nil
}
