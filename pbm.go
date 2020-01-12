// This file provides image support for both "raw" (binary) and
// "plain" (ASCII) Portable BitMap (PBM) files.

package netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"
	"unicode"
)

// A BW is simply an alias for an image.Paletted.  However, it is intended to
// represent images containing only white and black in their color palette.
type BW struct{ *image.Paletted }

// MaxValue returns the maximum index value allowed.
func (p *BW) MaxValue() uint16 {
	return 1
}

// Format identifies the image as a PBM image.
func (p *BW) Format() Format {
	return PBM
}

// HasAlpha indicates that there is no alpha channel.
func (p *BW) HasAlpha() bool {
	return false
}

// NewBW returns a new black-and-white image with the given bounds and
// maximum value.
func NewBW(r image.Rectangle) *BW {
	colorMap := make(color.Palette, 2)
	colorMap[0] = color.RGBA{255, 255, 255, 255}
	colorMap[1] = color.RGBA{0, 0, 0, 255}
	return &BW{image.NewPaletted(r, colorMap)}
}

// PromoteToGrayM generates an 8-bit grayscale image that looks identical to
// the given black-and-white image.  It takes as input a maximum channel value.
func (p *BW) PromoteToGrayM(m uint8) *GrayM {
	gray := NewGrayM(p.Bounds(), m)
	for i, bw := range p.Pix {
		gray.Pix[i] = (1 - bw) * m // PBM defines 0=white, 1=black.
	}
	return gray
}

// PromoteToGrayM32 generates an 16-bit grayscale image that looks identical to
// the given black-and-white image.  It takes as input a maximum channel value.
func (p *BW) PromoteToGrayM32(m uint16) *GrayM32 {
	gray := NewGrayM32(p.Bounds(), m)
	for i, bw := range p.Pix {
		g := uint16(1-bw) * m // PBM defines 0=white, 1=black.
		gray.Pix[i*2+0] = uint8(g >> 8)
		gray.Pix[i*2+1] = uint8(g)
	}
	return gray
}

// decodeConfigPBMWithComments reads and parses a PBM header, either "raw"
// (binary) or "plain" (ASCII).  Unlike decodeConfigPBM, it also returns any
// comments appearing in the file.
func decodeConfigPBMWithComments(r io.Reader) (image.Config, []string, error) {
	// We really want a bufio.Reader.  If we were given one, use it.  If
	// not, create a new one.
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	nr := newNetpbmReader(br)

	// Parse the PBM header.
	header, ok := nr.GetNetpbmHeader()
	if !ok {
		err := nr.Err()
		if err == nil {
			err = errors.New("Invalid PBM header")
		}
		return image.Config{}, nil, err
	}

	// Store the image configuration.
	var cfg image.Config
	cfg.Width = header.Width
	cfg.Height = header.Height

	// A PBM file's color map is 0=white, 1=black.
	colorMap := make(color.Palette, 2)
	colorMap[0] = color.RGBA{255, 255, 255, 255}
	colorMap[1] = color.RGBA{0, 0, 0, 255}
	cfg.ColorModel = colorMap
	return cfg, header.Comments, nil
}

// decodeConfigPBM reads and parses a PBM header, either "raw"
// (binary) or "plain" (ASCII).
func decodeConfigPBM(r io.Reader) (image.Config, error) {
	img, _, err := decodeConfigPBMWithComments(r)
	return img, err
}

// decodePBMWithComments reads a complete "raw" (binary) PBM image.  Unlike
// decodePBM, it also returns any comments appearing in the file.
func decodePBMWithComments(r io.Reader) (image.Image, []string, error) {
	// Read the image header, and use it to prepare a B&W image.
	br := bufio.NewReader(r)
	config, comments, err := decodeConfigPBMWithComments(br)
	if err != nil {
		return nil, nil, err
	}
	img := NewBW(image.Rect(0, 0, config.Width, config.Height))

	// Read bits until no more remain.
	nr := newNetpbmReader(br)
	buf := make([]byte, 1<<20) // Arbitrary, large, buffer size
	bitsRemaining := config.Width * config.Height
	bitNum := 0
ReadLoop:
	for {
		var nRead int
		nRead, err = nr.Read(buf)
		if nRead == 0 && err != nil {
			return nil, nil, err
		}
		for _, oneByte := range buf[:nRead] {
			for i := 7; i >= 0; i-- {
				img.Pix[bitNum] = uint8((oneByte >> uint8(i)) & 1)
				bitNum++
				bitsRemaining--
				if bitsRemaining == 0 {
					// We've read the entire image.
					break ReadLoop
				}
				if bitNum%config.Width == 0 {
					// Ignore row padding.
					break
				}
			}
		}
	}
	return img, comments, nil
}

// decodePBM reads a complete "raw" (binary) PBM image.
func decodePBM(r io.Reader) (image.Image, error) {
	img, _, err := decodePBMWithComments(r)
	return img, err
}

// decodePBMPlainWithComments reads a complete "plain" (ASCII) PBM image.
// Unlike decodePBMPlain, it also returns any comments appearing in the file.
func decodePBMPlainWithComments(r io.Reader) (image.Image, []string, error) {
	// Read the image header, and use it to prepare a B&W image.
	br := bufio.NewReader(r)
	config, comments, err := decodeConfigPBMWithComments(br)
	if err != nil {
		return nil, nil, err
	}
	img := NewBW(image.Rect(0, 0, config.Width, config.Height))

	// Define a simple error handler.
	nr := newNetpbmReader(br)
	badness := func() (image.Image, []string, error) {
		// Something went wrong.  Either we have an error code to
		// explain what or we make up a generic error message.
		err := nr.Err()
		if err == nil {
			err = errors.New("Failed to parse ASCII PBM data")
		}
		return img, nil, err
	}

	// Read bits (ASCII "0" or "1") until no more remain.
	totalBits := config.Width * config.Height
	for i := 0; i < totalBits; {
		ch := nr.GetNextByteAsRune()
		switch {
		case nr.Err() != nil:
			return badness()
		case unicode.IsSpace(ch):
			continue
		case ch == '0' || ch == '1':
			img.Pix[i] = uint8(ch - '0')
			i++
		default:
			return badness()
		}
	}
	return img, comments, nil
}

// decodePBMPlain reads a complete "plain" (ASCII) PBM image.
func decodePBMPlain(r io.Reader) (image.Image, error) {
	img, _, err := decodePBMPlainWithComments(r)
	return img, err
}

// Indicate that we can decode both raw and plain PBM files.
func init() {
	image.RegisterFormat("pbm", "P4", decodePBM, decodeConfigPBM)
	image.RegisterFormat("pbm", "P1", decodePBMPlain, decodeConfigPBM)
}

// encodePBM writes an arbitrary image in PBM format.
func encodePBM(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// Write the PBM header.
	if opts.Plain {
		fmt.Fprintln(w, "P1")
	} else {
		fmt.Fprintln(w, "P4")
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

	// Write the PBM data.
	return encodeBWData(w, img, opts)
}

// encodeBWData writes image data as 1-bit samples.
func encodeBWData(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// In the background, write each index value into a channel.
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	samples := make(chan uint16, width)
	go func() {
		bwImage := NewBW(image.ZR)
		cm := bwImage.ColorModel().(color.Palette)
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				samples <- uint16(cm.Index(img.At(x, y)))
			}
		}
		close(samples)
	}()

	// In the foreground, consume index values (either 0 or 1) and
	// write them to the image file as individual bits.  Pack 8
	// bits to a byte, pad each row, and output.
	if opts.Plain {
		return writePlainData(w, samples)
	}
	wb, ok := w.(*bufio.Writer)
	if !ok {
		wb = bufio.NewWriter(w)
	}
	var b byte      // Next byte to write
	var bLen uint   // Valid bits in b
	var rowBits int // Bits written to the current row
	for s := range samples {
		b = b<<1 | byte(s)
		bLen++
		rowBits++
		if rowBits == width {
			// Pad the last byte in the row.
			b <<= 8 - bLen
			bLen = 8
			rowBits = 0
		}
		if bLen == 8 {
			// Write a full byte to the output.
			if err := wb.WriteByte(b); err != nil {
				return err
			}
			b = 0
			bLen = 0
		}
	}
	wb.Flush()
	return nil
}
