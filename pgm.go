// This file provides image support for both "raw" (binary) and
// "plain" (ASCII) Portable GrayMap (PGM) files.

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

// GrayM is an in-memory image whose At method returns npcolor.GrayM values.
type GrayM struct {
	// Pix holds the image's pixels as gray values. The pixel at (x, y)
	// starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
	// Model is the image's color model.
	Model npcolor.GrayMModel
}

// ColorModel returns the GrayM image's color model.
func (p *GrayM) ColorModel() color.Model { return p.Model }

// Bounds returns the domain for which At can return non-zero color.  The
// bounds do not necessarily contain the point (0, 0).
func (p *GrayM) Bounds() image.Rectangle { return p.Rect }

// At returns the color of the pixel at (x, y) as a color.Color.
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (p *GrayM) At(x, y int) color.Color {
	return p.GrayMAt(x, y)
}

// GrayMAt returns the color of the pixel at (x, y) as an npcolor.GrayM.
func (p *GrayM) GrayMAt(x, y int) npcolor.GrayM {
	if !(image.Point{x, y}.In(p.Rect)) {
		return npcolor.GrayM{}
	}
	i := p.PixOffset(x, y)
	return npcolor.GrayM{Y: p.Pix[i], M: p.Model.M}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *GrayM) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
}

// Set sets the pixel at (x, y) to a given color, expressed as a color.Color.
func (p *GrayM) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = p.Model.Convert(c).(npcolor.GrayM).Y
}

// SetGrayM sets the pixel at (x, y) to a given color, expressed as an
// npcolor.GrayM.
func (p *GrayM) SetGrayM(x, y int, c npcolor.GrayM) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	if c.M == p.Model.M {
		p.Pix[i] = c.Y
	} else {
		p.Set(x, y, c)
	}
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *GrayM) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &GrayM{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &GrayM{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *GrayM) Opaque() bool {
	return true
}

// MaxValue returns the maximum grayscale value allowed.
func (p *GrayM) MaxValue() uint16 {
	return uint16(p.Model.M)
}

// Format identifies the image as a PGM image.
func (p *GrayM) Format() Format {
	return PGM
}

// HasAlpha indicates that there is no alpha channel.
func (p *GrayM) HasAlpha() bool {
	return false
}

// PromoteToRGBM generates an 8-bit color image that looks identical to
// the given grayscale image.
func (p *GrayM) PromoteToRGBM() *RGBM {
	rgb := NewRGBM(p.Bounds(), p.Model.M)
	for i, g := range p.Pix {
		rgb.Pix[i*3+0] = g
		rgb.Pix[i*3+1] = g
		rgb.Pix[i*3+2] = g
	}
	return rgb
}

// NewGrayM returns a new GrayM with the given bounds and maximum channel
// value.
func NewGrayM(r image.Rectangle, m uint8) *GrayM {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 1*w*h)
	model := npcolor.GrayMModel{M: m}
	return &GrayM{pix, 1 * w, r, model}
}

// GrayM32 is an in-memory image whose At method returns npcolor.GrayM32 values.
type GrayM32 struct {
	// Pix holds the image's pixels, as gray values. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
	// Model is the image's color model.
	Model npcolor.GrayM32Model
}

// ColorModel returns the GrayM32 image's color model.
func (p *GrayM32) ColorModel() color.Model { return p.Model }

// Bounds returns the domain for which At can return non-zero color.  The
// bounds do not necessarily contain the point (0, 0).
func (p *GrayM32) Bounds() image.Rectangle { return p.Rect }

// At returns the color of the pixel at (x, y) as a color.Color.
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (p *GrayM32) At(x, y int) color.Color {
	return p.GrayM32At(x, y)
}

// GrayM32At returns the color of the pixel at (x, y) as an npcolor.GrayM32.
func (p *GrayM32) GrayM32At(x, y int) npcolor.GrayM32 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return npcolor.GrayM32{}
	}
	i := p.PixOffset(x, y)
	return npcolor.GrayM32{
		Y: uint16(p.Pix[i+0])<<8 | uint16(p.Pix[i+1]),
		M: p.Model.M,
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *GrayM32) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*2
}

// Set sets the pixel at (x, y) to a given color, expressed as a color.Color.
func (p *GrayM32) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := p.Model.Convert(c).(npcolor.GrayM32)
	p.Pix[i+0] = uint8(c1.Y >> 8)
	p.Pix[i+1] = uint8(c1.Y)

}

// SetGrayM32 sets the pixel at (x, y) to a given color, expressed as an
// npcolor.GrayM32.
func (p *GrayM32) SetGrayM32(x, y int, c npcolor.GrayM32) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	if c.M == p.Model.M {
		p.Pix[i+0] = uint8(c.Y >> 8)
		p.Pix[i+1] = uint8(c.Y)
	} else {
		p.Set(x, y, c)
	}
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *GrayM32) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &GrayM32{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &GrayM32{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *GrayM32) Opaque() bool {
	return true
}

// MaxValue returns the maximum grayscale value allowed.
func (p *GrayM32) MaxValue() uint16 {
	return uint16(p.Model.M)
}

// Format identifies the image as a PGM image.
func (p *GrayM32) Format() Format {
	return PGM
}

// HasAlpha indicates that there is no alpha channel.
func (p *GrayM32) HasAlpha() bool {
	return false
}

// PromoteToRGBM64 generates a 16-bit color image that looks identical to
// the given grayscale image.
func (p *GrayM32) PromoteToRGBM64() *RGBM64 {
	rgb := NewRGBM64(p.Bounds(), p.Model.M)
	for i, g := range p.Pix {
		base := i / 2
		ofs := i % 2
		rgb.Pix[base*6+ofs+0] = g
		rgb.Pix[base*6+ofs+2] = g
		rgb.Pix[base*6+ofs+4] = g
	}
	return rgb
}

// NewGrayM32 returns a new GrayM32 with the given bounds and maximum channel
// value.
func NewGrayM32(r image.Rectangle, m uint16) *GrayM32 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 2*w*h)
	model := npcolor.GrayM32Model{M: m}
	return &GrayM32{pix, 2 * w, r, model}
}

// decodeConfigPGMWithComments reads and parses a PGM header, either "raw"
// (binary) or "plain" (ASCII).  Unlike decodeConfigPGM, it also returns any
// comments appearing in the file.
func decodeConfigPGMWithComments(r io.Reader) (image.Config, []string, error) {
	// We really want a bufio.Reader.  If we were given one, use it.  If
	// not, create a new one.
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	nr := newNetpbmReader(br)

	// Parse the PGM header.
	header, ok := nr.GetNetpbmHeader()
	if !ok {
		err := nr.Err()
		if err == nil {
			err = errors.New("Invalid PGM header")
		}
		return image.Config{}, nil, err
	}

	// Store and return the image configuration.
	var cfg image.Config
	cfg.Width = header.Width
	cfg.Height = header.Height
	if header.Maxval < 256 {
		cfg.ColorModel = npcolor.GrayMModel{M: uint8(header.Maxval)}
	} else {
		cfg.ColorModel = npcolor.GrayM32Model{M: uint16(header.Maxval)}
	}
	return cfg, header.Comments, nil
}

// decodeConfigPGM reads and parses a PGM header, either "raw"
// (binary) or "plain" (ASCII).
func decodeConfigPGM(r io.Reader) (image.Config, error) {
	img, _, err := decodeConfigPGMWithComments(r)
	return img, err
}

// decodePGMWithComments reads a complete "raw" (binary) PGM image.  Unlike
// decodePGM, it also returns any comments appearing in the file.
func decodePGMWithComments(r io.Reader) (image.Image, []string, error) {
	// Read the image header, and use it to prepare a grayscale image.
	br := bufio.NewReader(r)
	config, comments, err := decodeConfigPGMWithComments(br)
	if err != nil {
		return nil, nil, err
	}

	// Create either a Gray or a Gray16 image.
	var img image.Image // Image to return
	var data []uint8    // Image data
	var maxVal uint     // 100% white value
	switch model := config.ColorModel.(type) {
	case npcolor.GrayMModel:
		maxVal = uint(model.M)
		gray := NewGrayM(image.Rect(0, 0, config.Width, config.Height), uint8(maxVal))
		data = gray.Pix
		img = gray
	case npcolor.GrayM32Model:
		maxVal = uint(model.M)
		gray := NewGrayM32(image.Rect(0, 0, config.Width, config.Height), uint16(maxVal))
		data = gray.Pix
		img = gray
	default:
		panic("Unexpected color model")
	}

	// Raw PGM images are nice because we can read directly into the image
	// data.
	for len(data) > 0 {
		nRead, err := br.Read(data)
		if err != nil && err != io.EOF {
			return img, nil, err
		}
		if nRead == 0 {
			return img, nil, errors.New("Failed to read binary PGM data")
		}
		data = data[nRead:]
	}
	return img, comments, nil
}

// decodePGM reads a complete "raw" (binary) PGM image.
func decodePGM(r io.Reader) (image.Image, error) {
	img, _, err := decodePGMWithComments(r)
	return img, err
}

// decodePGMPlainWithComments reads a complete "plain" (ASCII) PGM image.
// Unlike decodePGMPlain, it also returns any comments appearing in the file.
func decodePGMPlainWithComments(r io.Reader) (image.Image, []string, error) {
	// Read the image header, and use it to prepare a grayscale image.
	br := bufio.NewReader(r)
	config, comments, err := decodeConfigPGMWithComments(br)
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
			err = errors.New("Failed to parse ASCII PGM data")
		}
		return img, nil, err
	}

	// Create either a Gray or a Gray16 image.
	var data []uint8 // Image data
	var maxVal int   // 100% white value
	switch model := config.ColorModel.(type) {
	case npcolor.GrayMModel:
		maxVal = int(model.M)
		gray := NewGrayM(image.Rect(0, 0, config.Width, config.Height), uint8(maxVal))
		data = gray.Pix
		img = gray
	case npcolor.GrayM32Model:
		maxVal = int(model.M)
		gray := NewGrayM32(image.Rect(0, 0, config.Width, config.Height), uint16(maxVal))
		data = gray.Pix
		img = gray
	default:
		panic("Unexpected color model")
	}

	// Read ASCII base-10 integers into the image data.
	if !nr.GetASCIIData(maxVal, data) {
		return badness()
	}
	return img, comments, nil
}

// decodePGMPlain reads a complete "plain" (ASCII) PGM image.
func decodePGMPlain(r io.Reader) (image.Image, error) {
	img, _, err := decodePGMPlainWithComments(r)
	return img, err
}

// Indicate that we can decode both raw and plain PGM files.
func init() {
	image.RegisterFormat("pgm", "P5", decodePGM, decodeConfigPGM)
	image.RegisterFormat("pgm", "P2", decodePGMPlain, decodeConfigPGM)
}

// encodePGM writes an arbitrary image in PGM format.
func encodePGM(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// Write the PGM header.
	if opts.Plain {
		fmt.Fprintln(w, "P2")
	} else {
		fmt.Fprintln(w, "P5")
	}
	for _, cmt := range opts.Comments {
		cmt = strings.Replace(cmt, "\n", " ", -1)
		cmt = strings.Replace(cmt, "\r", " ", -1)
		fmt.Fprintf(w, "# %s\n", cmt)
	}
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	fmt.Fprintf(w, "%d %d\n", width, height)
	fmt.Fprintf(w, "%d\n", opts.MaxValue)

	// Write the PGM data.
	if opts.MaxValue < 256 {
		return encodeGrayData(w, img, opts)
	}
	return encodeGray32Data(w, img, opts)
}

// encodeGrayData writes image data as 8-bit samples.
func encodeGrayData(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// In the background, write each 8-bit color sample into a channel.
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	samples := make(chan uint16, width)
	go func() {
		cm := npcolor.GrayMModel{M: uint8(opts.MaxValue)}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				c := cm.Convert(img.At(x, y)).(npcolor.GrayM)
				samples <- uint16(c.Y)
			}
		}
		close(samples)
	}()

	// In the foreground, consume grayscale samples and write them to the
	// image file.
	if opts.Plain {
		return writePlainData(w, samples)
	}
	return writeRawData(w, samples, 1)
}

// encodeGray32Data writes image data as 16-bit samples.
func encodeGray32Data(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// In the background, write each 16-bit color sample into a channel.
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	samples := make(chan uint16, width)
	go func() {
		cm := npcolor.GrayM32Model{M: opts.MaxValue}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				c := cm.Convert(img.At(x, y)).(npcolor.GrayM32)
				samples <- c.Y
			}
		}
		close(samples)
	}()

	// In the foreground, consume grayscale samples and write them to the
	// image file.
	if opts.Plain {
		return writePlainData(w, samples)
	}
	return writeRawData(w, samples, 2)
}
