// This file provides image support for both "raw" (binary) and
// "plain" (ASCII) Portable PixMap (PPM) files.

package netpbm

import (
	"bufio"
	"errors"
	"github.com/spakin/netpbm/npcolor"
	"image"
	"image/color"
	"io"
)

// An RGBM is an in-memory image whose At method returns npcolor.RGBM values.
type RGBM struct {
	// Pix holds the image's pixels, in R, G, B, M order. The pixel at (x,
	// y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent
	// pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel returns the RGBM image's color model.
func (p *RGBM) ColorModel() color.Model { return npcolor.RGBMModel }

// Bounds returns the domain for which At can return non-zero color.  The
// bounds do not necessarily contain the point (0, 0).
func (p *RGBM) Bounds() image.Rectangle { return p.Rect }

// At returns the color of the pixel at (x, y) as a color.Color.
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (p *RGBM) At(x, y int) color.Color {
	return p.RGBMAt(x, y)
}

// RGBMAt returns the color of the pixel at (x, y) as an npcolor.RGBM.
func (p *RGBM) RGBMAt(x, y int) npcolor.RGBM {
	if !(image.Point{x, y}.In(p.Rect)) {
		return npcolor.RGBM{}
	}
	i := p.PixOffset(x, y)
	return npcolor.RGBM{p.Pix[i+0], p.Pix[i+1], p.Pix[i+2], p.Pix[i+3]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGBM) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

// Set sets the pixel at (x, y) to a given color, expressed as a color.Color.
func (p *RGBM) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := npcolor.RGBMModel.Convert(c).(npcolor.RGBM)
	p.Pix[i+0] = c1.R
	p.Pix[i+1] = c1.G
	p.Pix[i+2] = c1.B
	p.Pix[i+3] = c1.M
}

// SetRGBM sets the pixel at (x, y) to a given color, expressed as an
// npcolor.RGBM.
func (p *RGBM) SetRGBM(x, y int, c npcolor.RGBM) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i+0] = c.R
	p.Pix[i+1] = c.G
	p.Pix[i+2] = c.B
	p.Pix[i+3] = c.M
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *RGBM) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &RGBM{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGBM{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *RGBM) Opaque() bool {
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

// NewRGBM returns a new RGBM with the given bounds.
func NewRGBM(r image.Rectangle) *RGBM {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 4*w*h)
	return &RGBM{pix, 4 * w, r}
}

// An RGBM64 is an in-memory image whose At method returns npcolor.RGBM64
// values.
type RGBM64 struct {
	// Pix holds the image's pixels, in R, G, B, M order and big-endian
	// format. The pixel at (x, y) starts at Pix[(y-Rect.Min.Y)*Stride +
	// (x-Rect.Min.X)*8].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent
	// pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel returns the RGBM64 image's color model.
func (p *RGBM64) ColorModel() color.Model { return npcolor.RGBM64Model }

// Bounds returns the domain for which At can return non-zero color.  The
// bounds do not necessarily contain the point (0, 0).
func (p *RGBM64) Bounds() image.Rectangle { return p.Rect }

// At returns the color of the pixel at (x, y) as a color.Color.
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (p *RGBM64) At(x, y int) color.Color {
	return p.RGBM64At(x, y)
}

// RGBM64At returns the color of the pixel at (x, y) as an npcolor.RGBM64.
func (p *RGBM64) RGBM64At(x, y int) npcolor.RGBM64 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return npcolor.RGBM64{}
	}
	i := p.PixOffset(x, y)
	return npcolor.RGBM64{
		uint16(p.Pix[i+0])<<8 | uint16(p.Pix[i+1]),
		uint16(p.Pix[i+2])<<8 | uint16(p.Pix[i+3]),
		uint16(p.Pix[i+4])<<8 | uint16(p.Pix[i+5]),
		uint16(p.Pix[i+6])<<8 | uint16(p.Pix[i+7]),
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGBM64) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*8
}

// Set sets the pixel at (x, y) to a given color, expressed as a color.Color.
func (p *RGBM64) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := npcolor.RGBM64Model.Convert(c).(npcolor.RGBM64)
	p.Pix[i+0] = uint8(c1.R >> 8)
	p.Pix[i+1] = uint8(c1.R)
	p.Pix[i+2] = uint8(c1.G >> 8)
	p.Pix[i+3] = uint8(c1.G)
	p.Pix[i+4] = uint8(c1.B >> 8)
	p.Pix[i+5] = uint8(c1.B)
	p.Pix[i+6] = uint8(c1.M >> 8)
	p.Pix[i+7] = uint8(c1.M)
}

// SetRGBM64 sets the pixel at (x, y) to a given color, expressed as an
// npcolor.RGBM.
func (p *RGBM64) SetRGBM64(x, y int, c npcolor.RGBM64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i+0] = uint8(c.R >> 8)
	p.Pix[i+1] = uint8(c.R)
	p.Pix[i+2] = uint8(c.G >> 8)
	p.Pix[i+3] = uint8(c.G)
	p.Pix[i+4] = uint8(c.B >> 8)
	p.Pix[i+5] = uint8(c.B)
	p.Pix[i+6] = uint8(c.M >> 8)
	p.Pix[i+7] = uint8(c.M)
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *RGBM64) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &RGBM64{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGBM64{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *RGBM64) Opaque() bool {
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

// NewRGBM64 returns a new RGBM64 with the given bounds.
func NewRGBM64(r image.Rectangle) *RGBM64 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 8*w*h)
	return &RGBM64{pix, 8 * w, r}
}

// decodeConfigPPM reads and parses a PPM header, either "raw" (binary) or
// "plain" (ASCII).
func decodeConfigPPM(r io.Reader) (image.Config, error) {
	// We really want a bufio.Reader.  If we were given one, use it.  If
	// not, create a new one.
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	nr := newNetpbmReader(br)

	// Parse the PPM header.
	header, ok := nr.GetNetpbmHeader()
	if !ok {
		err := nr.Err()
		if err == nil {
			err = errors.New("Invalid PPM header")
		}
		return image.Config{}, err
	}

	// Define the color model using the color channel's maximum value.
	if header.Maxval < 256 {
		header.Model = npcolor.RGBMModel
	} else {
		header.Model = npcolor.RGBM64Model
	}

	// Store and return the image configuration.
	var cfg image.Config
	cfg.Width = header.Width
	cfg.Height = header.Height
	cfg.ColorModel = header
	return cfg, nil
}

// decodePPM reads a complete "raw" (binary) PPM image.
func decodePPM(r io.Reader) (image.Image, error) {
	// Read the image header, and use it to prepare a color image.
	br := bufio.NewReader(r)
	config, err := decodeConfigPPM(br)
	if err != nil {
		return nil, err
	}

	// Create either a Color or a Color64 image.
	var img image.Image // Image to return
	var data []uint8    // RGBA image data
	var rgbData []uint8 // RGB file data
	nPixels := config.Width * config.Height
	maxVal := config.ColorModel.(netpbmHeader).Maxval // 100% white value
	if maxVal < 256 {
		rgb := NewRGBM(image.Rect(0, 0, config.Width, config.Height))
		data = rgb.Pix
		rgbData = make([]uint8, nPixels*3)
		img = rgb
	} else {
		rgb := NewRGBM64(image.Rect(0, 0, config.Width, config.Height))
		data = rgb.Pix
		rgbData = make([]uint8, nPixels*3*2)
		img = rgb
	}

	// Read RGB (no A) data into a holding buffer.
	rgbDataLeft := rgbData // RGB data left to read
	for len(rgbDataLeft) > 0 {
		nRead, err := br.Read(rgbDataLeft)
		if err != nil && err != io.EOF {
			return img, err
		}
		if nRead == 0 {
			return img, errors.New("Failed to read binary PPM data")
		}
		rgbDataLeft = rgbDataLeft[nRead:]
	}

	// Spread out RGB data into RGBA.
	nCopy := 3                       // Copy this many bytes from the input...
	nAlpha := 1                      // ...then generate this many bytes of alpha.
	opaque := []uint8{uint8(maxVal)} // Maximum opacity
	if maxVal >= 256 {
		nCopy *= 2
		nAlpha *= 2
		opaque = []uint8{uint8(maxVal >> 8), uint8(maxVal)}
	}
	for p, s, d := 0, 0, 0; p < nPixels; p++ {
		copy(data[d:d+nCopy], rgbData[s:s+nCopy])
		copy(data[d+nCopy:d+nCopy+nAlpha], opaque)
		s += nCopy
		d += nCopy + nAlpha
	}
	return img, nil
}

// decodePPMPlain reads a complete "plain" (ASCII) PPM image.
func decodePPMPlain(r io.Reader) (image.Image, error) {
	// Read the image header, and use it to prepare a color image.
	br := bufio.NewReader(r)
	config, err := decodeConfigPPM(br)
	if err != nil {
		return nil, err
	}
	var img image.Image // Image to return

	// Define a simple error handler.
	nr := newNetpbmReader(br)
	badness := func() (image.Image, error) {
		// Something went wrong.  Either we have an error code to
		// explain what or we make up a generic error message.
		err := nr.Err()
		if err == nil {
			err = errors.New("Failed to parse ASCII PPM data")
		}
		return img, err
	}

	// Create either a Color or a Color64 image.
	var data []uint8                                  // Image data
	maxVal := config.ColorModel.(netpbmHeader).Maxval // 100% white value
	if maxVal < 256 {
		rgb := NewRGBM(image.Rect(0, 0, config.Width, config.Height))
		data = rgb.Pix
		img = rgb
	} else {
		rgb := NewRGBM64(image.Rect(0, 0, config.Width, config.Height))
		data = rgb.Pix
		img = rgb
	}

	// Read ASCII base-10 integers until no more remain.
	if maxVal < 256 {
		for i := 0; i < len(data); {
			for d := 0; d < 3; d++ {
				val := nr.GetNextInt()
				switch {
				case nr.Err() != nil:
					return badness()
				case val < 0 || val > maxVal:
					return badness()
				default:
					data[i] = uint8(val)
					i++
				}
			}
			data[i] = uint8(maxVal)
			i++
		}
	} else {
		for i := 0; i < len(data); {
			for d := 0; d < 3; d++ {
				val := nr.GetNextInt()
				switch {
				case nr.Err() != nil:
					return badness()
				case val < 0 || val > maxVal:
					return badness()
				default:
					data[i] = uint8(val >> 8)
					data[i+1] = uint8(val)
					i += 2
				}
			}
			data[i] = uint8(maxVal >> 8)
			data[i+1] = uint8(maxVal)
			i += 2
		}
	}
	return img, nil
}

// Indicate that we can decode both raw and plain PPM files.
func init() {
	image.RegisterFormat("ppm", "P6", decodePPM, decodeConfigPPM)
	image.RegisterFormat("ppm", "P3", decodePPMPlain, decodeConfigPPM)
}
