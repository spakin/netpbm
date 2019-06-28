// This file provides image support for both "raw" (binary) and
// "plain" (ASCII) Portable PixMap (PPM) files.

package netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"

	"github.com/spakin/netpbm/npcolor"
)

// An RGBM is an in-memory image whose At method returns npcolor.RGBM values.
type RGBM struct {
	// Pix holds the image's pixels, in R, G, B (no M) order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*3].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent
	// pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
	// Model is the image's color model.
	Model npcolor.RGBMModel
}

// ColorModel returns the RGBM image's color model.
func (p *RGBM) ColorModel() color.Model { return p.Model }

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
	return npcolor.RGBM{
		R: p.Pix[i+0],
		G: p.Pix[i+1],
		B: p.Pix[i+2],
		M: p.Model.M,
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGBM) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

// Set sets the pixel at (x, y) to a given color, expressed as a color.Color.
func (p *RGBM) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := p.Model.Convert(c).(npcolor.RGBM)
	p.Pix[i+0] = c1.R
	p.Pix[i+1] = c1.G
	p.Pix[i+2] = c1.B
}

// SetRGBM sets the pixel at (x, y) to a given color, expressed as an
// npcolor.RGBM.
func (p *RGBM) SetRGBM(x, y int, c npcolor.RGBM) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	if c.M == p.Model.M {
		p.Pix[i+0] = c.R
		p.Pix[i+1] = c.G
		p.Pix[i+2] = c.B
	} else {
		p.Set(x, y, c)
	}
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
	return true
}

// MaxValue returns the maximum value allowed on any color channel.
func (p *RGBM) MaxValue() uint16 {
	return uint16(p.Model.M)
}

// Format identifies the image as a PPM image.
func (p *RGBM) Format() Format {
	return PPM
}

// HasAlpha indicates that there is no alpha channel.
func (p *RGBM) HasAlpha() bool {
	return false
}

// NewRGBM returns a new RGBM with the given bounds and maximum channel value.
func NewRGBM(r image.Rectangle, m uint8) *RGBM {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 3*w*h)
	model := npcolor.RGBMModel{M: m}
	return &RGBM{pix, 3 * w, r, model}
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
	// Model is the image's color model.
	Model npcolor.RGBM64Model
}

// ColorModel returns the RGBM64 image's color model.
func (p *RGBM64) ColorModel() color.Model { return p.Model }

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
		R: uint16(p.Pix[i+0])<<8 | uint16(p.Pix[i+1]),
		G: uint16(p.Pix[i+2])<<8 | uint16(p.Pix[i+3]),
		B: uint16(p.Pix[i+4])<<8 | uint16(p.Pix[i+5]),
		M: p.Model.M,
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGBM64) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*6
}

// Set sets the pixel at (x, y) to a given color, expressed as a color.Color.
func (p *RGBM64) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := p.Model.Convert(c).(npcolor.RGBM64)
	p.Pix[i+0] = uint8(c1.R >> 8)
	p.Pix[i+1] = uint8(c1.R)
	p.Pix[i+2] = uint8(c1.G >> 8)
	p.Pix[i+3] = uint8(c1.G)
	p.Pix[i+4] = uint8(c1.B >> 8)
	p.Pix[i+5] = uint8(c1.B)
}

// SetRGBM64 sets the pixel at (x, y) to a given color, expressed as an
// npcolor.RGBM.
func (p *RGBM64) SetRGBM64(x, y int, c npcolor.RGBM64) {
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
	} else {
		p.Set(x, y, c)
	}
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
	return true
}

// MaxValue returns the maximum value allowed on any color channel.
func (p *RGBM64) MaxValue() uint16 {
	return uint16(p.Model.M)
}

// Format identifies the image as a PPM image.
func (p *RGBM64) Format() Format {
	return PPM
}

// HasAlpha indicates that there is no alpha channel.
func (p *RGBM64) HasAlpha() bool {
	return false
}

// NewRGBM64 returns a new RGBM64 with the given bounds and maximum channel
// value.
func NewRGBM64(r image.Rectangle, m uint16) *RGBM64 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 6*w*h)
	model := npcolor.RGBM64Model{M: m}
	return &RGBM64{pix, 6 * w, r, model}
}

// decodeConfigPPMWithComments reads and parses a PPM header, either "raw"
// (binary) or "plain" (ASCII).  Unlike decodeConfigPPM, it also returns any
// comments appearing in the file.
func decodeConfigPPMWithComments(r io.Reader) (image.Config, []string, error) {
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
		return image.Config{}, nil, err
	}

	// Store and return the image configuration.
	var cfg image.Config
	cfg.Width = header.Width
	cfg.Height = header.Height
	if header.Maxval < 256 {
		cfg.ColorModel = npcolor.RGBMModel{M: uint8(header.Maxval)}
	} else {
		cfg.ColorModel = npcolor.RGBM64Model{M: uint16(header.Maxval)}
	}
	return cfg, header.Comments, nil
}

// decodeConfigPPM reads and parses a PPM header, either "raw"
// (binary) or "plain" (ASCII).
func decodeConfigPPM(r io.Reader) (image.Config, error) {
	img, _, err := decodeConfigPPMWithComments(r)
	return img, err
}

// decodePPMWithComments reads a complete "raw" (binary) PPM image.  Unlike
// decodePPM, it also returns any comments appearing in the file.
func decodePPMWithComments(r io.Reader) (image.Image, []string, error) {
	// Read the image header, and use it to prepare a color image.
	br := bufio.NewReader(r)
	config, comments, err := decodeConfigPPMWithComments(br)
	if err != nil {
		return nil, nil, err
	}

	// Create either a Color or a Color64 image.
	var img image.Image // Image to return
	var data []uint8    // RGB (no M) image data
	var maxVal uint     // 100% white value
	switch model := config.ColorModel.(type) {
	case npcolor.RGBMModel:
		maxVal = uint(model.M)
		rgb := NewRGBM(image.Rect(0, 0, config.Width, config.Height), uint8(maxVal))
		data = rgb.Pix
		img = rgb
	case npcolor.RGBM64Model:
		maxVal = uint(model.M)
		rgb := NewRGBM64(image.Rect(0, 0, config.Width, config.Height), uint16(maxVal))
		data = rgb.Pix
		img = rgb
	default:
		panic("Unexpected color model")
	}

	// Raw PPM images are nice because we can read directly into the image
	// data.
	for len(data) > 0 {
		nRead, err := br.Read(data)
		if err != nil && err != io.EOF {
			return img, nil, err
		}
		if nRead == 0 {
			return img, nil, errors.New("Failed to read binary PPM data")
		}
		data = data[nRead:]
	}
	return img, comments, nil
}

// decodePPM reads a complete "raw" (binary) PPM image.
func decodePPM(r io.Reader) (image.Image, error) {
	img, _, err := decodePPMWithComments(r)
	return img, err
}

// decodePPMPlainWithComments reads a complete "plain" (ASCII) PPM image.
// Unlike decodePPMPlain, it also returns any comments appearing in the file.
func decodePPMPlainWithComments(r io.Reader) (image.Image, []string, error) {
	// Read the image header, and use it to prepare a color image.
	br := bufio.NewReader(r)
	config, comments, err := decodeConfigPPMWithComments(br)
	if err != nil {
		return nil, nil, err
	}
	var img image.Image // Image to return

	// Define a simple error handler.
	nr := newNetpbmReader(br)
	badness := func() (image.Image, []string, error) {
		// Something went wrong.  Either we have an error code to
		// explain what or we make up a generic error message.
		err := nr.Err()
		if err == nil {
			err = errors.New("Failed to parse ASCII PPM data")
		}
		return img, nil, err
	}

	// Create either a Color or a Color64 image.
	var data []uint8 // Image data
	var maxVal int   // 100% white value
	switch model := config.ColorModel.(type) {
	case npcolor.RGBMModel:
		maxVal = int(model.M)
		rgb := NewRGBM(image.Rect(0, 0, config.Width, config.Height), uint8(maxVal))
		data = rgb.Pix
		img = rgb
	case npcolor.RGBM64Model:
		maxVal = int(model.M)
		rgb := NewRGBM64(image.Rect(0, 0, config.Width, config.Height), uint16(maxVal))
		data = rgb.Pix
		img = rgb
	default:
		panic("Unexpected color model")
	}

	// Read ASCII base-10 integers until no more remain.
	if !nr.GetASCIIData(maxVal, data) {
		return badness()
	}
	return img, comments, nil
}

// decodePPMPlain reads a complete "plain" (ASCII) PPM image.
func decodePPMPlain(r io.Reader) (image.Image, error) {
	img, _, err := decodePPMPlainWithComments(r)
	return img, err
}

// Indicate that we can decode both raw and plain PPM files.
func init() {
	image.RegisterFormat("ppm", "P6", decodePPM, decodeConfigPPM)
	image.RegisterFormat("ppm", "P3", decodePPMPlain, decodeConfigPPM)
}

// encodePPM writes an arbitrary image in PPM format.
func encodePPM(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// Write the PPM header.
	if opts.Plain {
		fmt.Fprintln(w, "P3")
	} else {
		fmt.Fprintln(w, "P6")
	}
	for _, cmt := range opts.Comments {
		cmt = strings.ReplaceAll(cmt, "\n", " ")
		cmt = strings.ReplaceAll(cmt, "\r", " ")
		fmt.Fprintf(w, "# %s\n", cmt)
	}
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	fmt.Fprintf(w, "%d %d\n", width, height)
	fmt.Fprintf(w, "%d\n", opts.MaxValue)

	// Write the PPM data.
	if opts.MaxValue < 256 {
		return encodeRGBData(w, img, opts)
	}
	return encodeRGB64Data(w, img, opts)
}

// encodeRGBData writes image data as 8-bit samples.
func encodeRGBData(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// In the background, write each 8-bit color sample into a channel.
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	samples := make(chan uint16, width*3)
	go func() {
		cm := npcolor.RGBMModel{M: uint8(opts.MaxValue)}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				c := cm.Convert(img.At(x, y)).(npcolor.RGBM)
				samples <- uint16(c.R)
				samples <- uint16(c.G)
				samples <- uint16(c.B)
			}
		}
		close(samples)
	}()

	// In the foreground, consume color samples and write them to the image
	// file.
	if opts.Plain {
		return writePlainData(w, samples)
	}
	return writeRawData(w, samples, 1)
}

// encodeRGB64Data writes image data as 16-bit samples.
func encodeRGB64Data(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// In the background, write each 16-bit color sample into a channel.
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	samples := make(chan uint16, width*3)
	go func() {
		cm := npcolor.RGBM64Model{M: opts.MaxValue}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				c := cm.Convert(img.At(x, y)).(npcolor.RGBM64)
				samples <- c.R
				samples <- c.G
				samples <- c.B
			}
		}
		close(samples)
	}()

	// In the foreground, consume color samples and write them to the image
	// file.
	if opts.Plain {
		return writePlainData(w, samples)
	}
	return writeRawData(w, samples, 2)
}
