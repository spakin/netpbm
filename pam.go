// This file provides image support for Portable Arbitrary Map (PAM) files.

package netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spakin/netpbm/npcolor"
	"image"
	"image/color"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// Define a type representing a known tuple type.
type pamTupleType int

// These are the allowed values of a pamTupleType.
const (
	pamBlackAndWhite pamTupleType = iota
	pamBlackAndWhiteAlpha
	pamGrayscale
	pamGrayscaleAlpha
	pamColor
	pamColorAlpha
)

// ttToInt maps a PAM tuple type from a string to an integer.
var ttToInt = map[string]pamTupleType{
	"BLACKANDWHITE":       pamBlackAndWhite,
	"BLACKANDWHITE_ALPHA": pamBlackAndWhiteAlpha,
	"GRAYSCALE":           pamGrayscale,
	"GRAYSCALE_ALPHA":     pamGrayscaleAlpha,
	"RGB":                 pamColor,
	"RGB_ALPHA":           pamColorAlpha,
}

// A PAMRGBM is identical to an RGBM but is of PAM format, not PPM format.
type PAMRGBM struct{ RGBM }

// NewPAMRGBM returns a new PAMRGBM with the given bounds and maximum channel
// value.
func NewPAMRGBM(r image.Rectangle, m uint8) *PAMRGBM {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 3*w*h)
	model := npcolor.RGBMModel{M: m}
	pam := &PAMRGBM{}
	pam.Pix = pix
	pam.Stride = 3 * w
	pam.Rect = r
	pam.Model = model
	return pam
}

// Format identifies the image as a PAM image.
func (p *PAMRGBM) Format() Format {
	return PAM
}

// A PAMRGBM64 is identical to an RGBM64 but is of PAM format, not PPM format.
type PAMRGBM64 struct{ RGBM64 }

// NewPAMRGBM64 returns a new PAMRGBM64 with the given bounds and maximum
// channel value.
func NewPAMRGBM64(r image.Rectangle, m uint16) *PAMRGBM64 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 6*w*h)
	model := npcolor.RGBM64Model{M: m}
	pam := &PAMRGBM64{}
	pam.Pix = pix
	pam.Stride = 6 * w
	pam.Rect = r
	pam.Model = model
	return pam
}

// Format identifies the image as a PAM image.
func (p *PAMRGBM64) Format() Format {
	return PAM
}

// A PAMRGBAM is an in-memory image whose At method returns npcolor.RGBAM
// values.
type PAMRGBAM struct {
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

// ColorModel returns the PAMRGBAM image's color model.
func (p *PAMRGBAM) ColorModel() color.Model { return p.Model }

// Bounds returns the domain for which At can return non-zero color.  The
// bounds do not necessarily contain the point (0, 0).
func (p *PAMRGBAM) Bounds() image.Rectangle { return p.Rect }

// At returns the color of the pixel at (x, y) as a color.Color.
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (p *PAMRGBAM) At(x, y int) color.Color {
	return p.RGBAMAt(x, y)
}

// RGBAMAt returns the color of the pixel at (x, y) as an npcolor.RGBAM.
func (p *PAMRGBAM) RGBAMAt(x, y int) npcolor.RGBAM {
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
func (p *PAMRGBAM) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

// Set sets the pixel at (x, y) to a given color, expressed as a color.Color.
func (p *PAMRGBAM) Set(x, y int, c color.Color) {
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
func (p *PAMRGBAM) SetRGBAM(x, y int, c npcolor.RGBAM) {
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
func (p *PAMRGBAM) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &PAMRGBAM{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &PAMRGBAM{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *PAMRGBAM) Opaque() bool {
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
func (p *PAMRGBAM) MaxValue() uint16 {
	return uint16(p.Model.M)
}

// Format identifies the image as a PPM image.
func (p *PAMRGBAM) Format() Format {
	return PAM
}

// NewPAMRGBAM returns a new PAMRGBAM with the given bounds and maximum channel
// value.
func NewPAMRGBAM(r image.Rectangle, m uint8) *PAMRGBAM {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 4*w*h)
	model := npcolor.RGBAMModel{M: m}
	return &PAMRGBAM{pix, 4 * w, r, model}
}

// A PAMRGBAM64 is an in-memory image whose At method returns npcolor.RGBAM64
// values.
type PAMRGBAM64 struct {
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

// ColorModel returns the PAMRGBAM64 image's color model.
func (p *PAMRGBAM64) ColorModel() color.Model { return p.Model }

// Bounds returns the domain for which At can return non-zero color.  The
// bounds do not necessarily contain the point (0, 0).
func (p *PAMRGBAM64) Bounds() image.Rectangle { return p.Rect }

// At returns the color of the pixel at (x, y) as a color.Color.
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (p *PAMRGBAM64) At(x, y int) color.Color {
	return p.RGBAM64At(x, y)
}

// RGBAM64At returns the color of the pixel at (x, y) as an npcolor.RGBAM64.
func (p *PAMRGBAM64) RGBAM64At(x, y int) npcolor.RGBAM64 {
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
func (p *PAMRGBAM64) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*8
}

// Set sets the pixel at (x, y) to a given color, expressed as a color.Color.
func (p *PAMRGBAM64) Set(x, y int, c color.Color) {
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
func (p *PAMRGBAM64) SetRGBAM64(x, y int, c npcolor.RGBAM64) {
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
func (p *PAMRGBAM64) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &PAMRGBAM64{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &PAMRGBAM64{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *PAMRGBAM64) Opaque() bool {
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
func (p *PAMRGBAM64) MaxValue() uint16 {
	return uint16(p.Model.M)
}

// Format identifies the image as a PPM image.
func (p *PAMRGBAM64) Format() Format {
	return PAM
}

// NewPAMRGBAM64 returns a new PAMRGBAM64 with the given bounds and maximum
// channel value.
func NewPAMRGBAM64(r image.Rectangle, m uint16) *PAMRGBAM64 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 8*w*h)
	model := npcolor.RGBAM64Model{M: m}
	return &PAMRGBAM64{pix, 8 * w, r, model}
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

	// Process each line in turn.
ReadLoop:
	for {
		// Read a line.
		kv := nr.GetLineAsKeyValue()
		if nr.Err() != nil {
			return netpbmHeader{}, false
		}
		if len(kv) == 0 {
			continue
		}
		if len(kv) == 1 && kv[0] != "ENDHDR" {
			return netpbmHeader{}, false
		}

		// Parse the line.
		var err error
		k, v := kv[0], kv[1]
		switch k {
		case "ENDHDR":
			break ReadLoop
		case "HEIGHT":
			header.Height, err = strconv.Atoi(v)
		case "WIDTH":
			header.Width, err = strconv.Atoi(v)
		case "DEPTH":
			header.Depth, err = strconv.Atoi(v)
		case "MAXVAL":
			header.Maxval, err = strconv.Atoi(v)
		case "TUPLTYPE":
			if header.TupleType != "" {
				header.TupleType += " "
			}
			header.TupleType += v
		case "#":
			header.Comments = append(header.Comments, v)
		default:
			return netpbmHeader{}, false
		}
		if err != nil {
			return netpbmHeader{}, false
		}
	}
	if header.Maxval < 1 || header.Maxval > 65535 {
		return netpbmHeader{}, false
	}

	// Return the header and a success code.
	return header, true
}

// decodeConfigPAMWithComments reads and parses a PAM header.  Unlike
// decodeConfigPAM, it also returns any comments appearing in the file.
func decodeConfigPAMWithComments(r io.Reader) (image.Config, []string, error) {
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
		return image.Config{}, nil, err
	}

	// Store and return the image configuration.
	var cfg image.Config
	cfg.Width = header.Width
	cfg.Height = header.Height
	ttype, ok := ttToInt[header.TupleType]
	if !ok {
		return image.Config{}, nil, fmt.Errorf("Unsupported tuple type %q", header.TupleType)
	}
	if header.Maxval < 256 {
		switch ttype {
		case pamColorAlpha:
			cfg.ColorModel = npcolor.RGBAMModel{M: uint8(header.Maxval)}
		case pamColor:
			cfg.ColorModel = npcolor.RGBMModel{M: uint8(header.Maxval)}
		case pamGrayscaleAlpha:
			// TODO: Implement grayscale + alpha
			panic("Grayscale + alpha is not currently supported")
		case pamGrayscale:
			cfg.ColorModel = npcolor.GrayMModel{M: uint8(header.Maxval)}
		case pamBlackAndWhiteAlpha:
			// TODO: Implement BW + alpha
			panic("Black & white + alpha is not currently supported")
		case pamBlackAndWhite:
			// Define a color map with 0=black and 1=white.
			colorMap := make(color.Palette, 2)
			colorMap[0] = color.RGBA{0, 0, 0, 255}
			colorMap[1] = color.RGBA{255, 255, 255, 255}
			cfg.ColorModel = colorMap
		default:
			panic(fmt.Sprintf("Internal error processing tuple type %q", header.TupleType))
		}
	} else {
		switch ttype {
		case pamColorAlpha:
			cfg.ColorModel = npcolor.RGBAM64Model{M: uint16(header.Maxval)}
		case pamColor:
			cfg.ColorModel = npcolor.RGBM64Model{M: uint16(header.Maxval)}
		case pamGrayscaleAlpha:
			// TODO: Implement grayscale + alpha
			panic("Grayscale + alpha is not currently supported")
		case pamGrayscale:
			cfg.ColorModel = npcolor.GrayM32Model{M: uint16(header.Maxval)}
		case pamBlackAndWhiteAlpha:
			// TODO: Implement BW + alpha
			panic("Black & white + alpha is not currently supported")
		case pamBlackAndWhite:
			// Define a color map with 0=black and 1=white.
			colorMap := make(color.Palette, 2)
			colorMap[0] = color.RGBA{0, 0, 0, 255}
			colorMap[1] = color.RGBA{255, 255, 255, 255}
			cfg.ColorModel = colorMap
		default:
			panic(fmt.Sprintf("Internal error processing tuple type %q", header.TupleType))
		}
	}
	return cfg, header.Comments, nil
}

// decodeConfigPAM reads and parses a PAM header.
func decodeConfigPAM(r io.Reader) (image.Config, error) {
	img, _, err := decodeConfigPAMWithComments(r)
	return img, err
}

// decodePAMWithComments reads a complete PAM image.  Unlike decodePAM, it also
// returns any comments appearing in the file.
func decodePAMWithComments(r io.Reader) (image.Image, []string, error) {
	// Read the image header, and use it to prepare a color image.
	br := bufio.NewReader(r)
	config, comments, err := decodeConfigPAMWithComments(br)
	if err != nil {
		return nil, nil, err
	}

	// Create an appropriate image type.
	var img image.Image // Image to return
	var data []uint8    // RGB (no M) image data
	var maxVal uint     // 100% white value
	switch model := config.ColorModel.(type) {
	case npcolor.RGBAMModel:
		maxVal = uint(model.M)
		pImg := NewPAMRGBAM(image.Rect(0, 0, config.Width, config.Height), uint8(maxVal))
		data = pImg.Pix
		img = pImg

	case npcolor.RGBMModel:
		maxVal = uint(model.M)
		pImg := NewPAMRGBM(image.Rect(0, 0, config.Width, config.Height), uint8(maxVal))
		data = pImg.Pix
		img = pImg

	case npcolor.GrayMModel:
		maxVal = uint(model.M)
		pImg := NewGrayM(image.Rect(0, 0, config.Width, config.Height), uint8(maxVal))
		data = pImg.Pix
		img = pImg

	case npcolor.RGBAM64Model:
		maxVal = uint(model.M)
		pImg := NewPAMRGBAM64(image.Rect(0, 0, config.Width, config.Height), uint16(maxVal))
		data = pImg.Pix
		img = pImg

	case npcolor.RGBM64Model:
		maxVal = uint(model.M)
		pImg := NewPAMRGBM64(image.Rect(0, 0, config.Width, config.Height), uint16(maxVal))
		data = pImg.Pix
		img = pImg

	case npcolor.GrayM32Model:
		maxVal = uint(model.M)
		pImg := NewGrayM32(image.Rect(0, 0, config.Width, config.Height), uint16(maxVal))
		data = pImg.Pix
		img = pImg

	default:
		panic("Unexpected color model")
	}

	// PAM images are nice because we can read directly into the image
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

// decodePAM reads a complete PAM image.
func decodePAM(r io.Reader) (image.Image, error) {
	img, _, err := decodePAMWithComments(r)
	return img, err
}

// Indicate that we can decode PAM files.
func init() {
	image.RegisterFormat("pam", "P7", decodePAM, decodeConfigPAM)
}

// encodePAM writes an arbitrary image in PAM format.
func encodePAM(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// Determine the depth from the tuple type.
	var depth int
	ttype, ok := ttToInt[opts.TupleType]
	if !ok {
		return fmt.Errorf("Unsupported tuple type %q", opts.TupleType)
	}
	switch ttype {
	case pamColorAlpha:
		depth = 4
	case pamColor:
		depth = 3
	case pamGrayscaleAlpha:
		depth = 2
	case pamGrayscale:
		depth = 1
	case pamBlackAndWhiteAlpha:
		depth = 2
	case pamBlackAndWhite:
		depth = 1
	default:
		panic(fmt.Sprintf("Internal error processing tuple type %q", opts.TupleType))
	}

	// Write the PAM header.
	fmt.Fprintln(w, "P7")
	for _, cmt := range opts.Comments {
		cmt = strings.ReplaceAll(cmt, "\n", " ")
		cmt = strings.ReplaceAll(cmt, "\r", " ")
		fmt.Fprintf(w, "# %s\n", cmt)
	}
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	fmt.Fprintf(w, "WIDTH %d\n", width)
	fmt.Fprintf(w, "HEIGHT %d\n", height)
	fmt.Fprintf(w, "DEPTH %d\n", depth)
	fmt.Fprintf(w, "MAXVAL %d\n", opts.MaxValue)
	fmt.Fprintf(w, "TUPLTYPE %s\n", opts.TupleType)
	fmt.Fprintf(w, "ENDHDR\n")

	// Write the PPM data.
	if opts.MaxValue < 256 {
		switch ttype {
		case pamColorAlpha:
			return encodeRGBAData(w, img, opts)
		case pamColor:
			return encodeRGBData(w, img, opts)
		case pamGrayscaleAlpha:
			// TODO: Implement grayscale + alpha
			panic("Grayscale + alpha is not currently supported")
		case pamGrayscale:
			return encodeGrayData(w, img, opts)
		case pamBlackAndWhiteAlpha:
			// TODO: Implement BW + alpha
			panic("Black & white + alpha is not currently supported")
		case pamBlackAndWhite:
			return encodeBWData(w, img, opts)
		default:
			panic(fmt.Sprintf("Internal error processing tuple type %q", opts.TupleType))
		}
	} else {
		switch ttype {
		case pamColorAlpha:
			return encodeRGBA64Data(w, img, opts)
		case pamColor:
			return encodeRGB64Data(w, img, opts)
		case pamGrayscaleAlpha:
			// TODO: Implement 16-bit grayscale + alpha
			panic("16-bit grayscale + alpha is not currently supported")
		case pamGrayscale:
			return encodeGray32Data(w, img, opts)
		case pamBlackAndWhiteAlpha:
			// TODO: Implement 16-bit BW + alpha
			panic("16-bit Black & white + alpha is not currently supported")
		case pamBlackAndWhite:
			return encodeBWData(w, img, opts)
		default:
			panic(fmt.Sprintf("Internal error processing tuple type %q", opts.TupleType))
		}
	}
}

// encodeRGBAData writes image data as 8-bit samples.
func encodeRGBAData(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// In the background, write each 8-bit color sample into a channel.
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	samples := make(chan uint16, width*3)
	go func() {
		cm := npcolor.RGBAMModel{M: uint8(opts.MaxValue)}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				c := cm.Convert(img.At(x, y)).(npcolor.RGBAM)
				samples <- uint16(c.R)
				samples <- uint16(c.G)
				samples <- uint16(c.B)
				samples <- uint16(c.A)
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

// encodeRGBA64Data writes image data as 16-bit samples.
func encodeRGBA64Data(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// In the background, write each 16-bit color sample into a channel.
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	samples := make(chan uint16, width*3)
	go func() {
		cm := npcolor.RGBAM64Model{M: opts.MaxValue}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				c := cm.Convert(img.At(x, y)).(npcolor.RGBAM64)
				samples <- c.R
				samples <- c.G
				samples <- c.B
				samples <- c.A
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
