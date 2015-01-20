// Package pbm encodes and decodes binary Portable BitMap (PBM) images.
package pbm

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"io"
	"unicode"
)

// A netpbmReader extends bufio.Reader with the ability to read bytes
// and numbers while skipping over comments.
type netpbmReader struct {
	*bufio.Reader        // Inherit Read, UnreadByte, etc.
	err          error  // Sticky error state
	oneByte      []byte // Buffer into which to read a single byte
}

// newNetpbmReader allocates, initializes, and returns a new netpbmReader.
func newNetpbmReader(r *bufio.Reader) *netpbmReader {
	return &netpbmReader{Reader: r, oneByte: make([]byte, 1)}
}

// Err returns the netpbmReader's current error state.
func (nr netpbmReader) Err() error {
	return nr.err
}

// GetNextByteAsRune returns the next byte, cast to a rune, or 0 on error (and
// errors are sticky).
func (nr *netpbmReader) GetNextByteAsRune() rune {
	if nr.err != nil {
		return 0
	}
	_, nr.err = nr.Read(nr.oneByte)
	if nr.err != nil {
		return 0
	}
	return rune(nr.oneByte[0])
}

// GetNextInt returns the next base-10 integer read from a netpbmReader,
// skipping whitespace and comments.
func (nr *netpbmReader) GetNextInt() int {
	// Find the first digit.
	var c rune
	for nr.err == nil && !unicode.IsDigit(c) {
		for c = nr.GetNextByteAsRune(); unicode.IsSpace(c); c = nr.GetNextByteAsRune() {
		}
		if c == '#' {
			// Comment -- discard the rest of the line.
			for c = nr.GetNextByteAsRune(); c != '\n'; c = nr.GetNextByteAsRune() {
			}
		}
	}
	if nr.err != nil {
		return -1
	}

	// Read while we have base-10 digits.  Return the resulting int.
	value := int(c - '0')
	for c = nr.GetNextByteAsRune(); unicode.IsDigit(c); c = nr.GetNextByteAsRune() {
		value = value*10 + int(c-'0')
	}
	if nr.err != nil {
		return -1
	}
	nr.err = nr.UnreadByte()
	if nr.err != nil {
		return -1
	}
	return value
}

// DecodeConfig reads and parses the PBM header.
func DecodeConfig(r io.Reader) (image.Config, error) {
	// The PBM header is ASCII.  Define a bunch of helper functions to
	// parse it.
	nr := newNetpbmReader(bufio.NewReader(r))
	var err error
	badness := func() (image.Config, error) {
		// Something went wrong.  Either we have an error code to
		// explain what or we make up a generic error message.
		if err == nil {
			err = errors.New("Invalid PBM header")
		}
		return image.Config{}, err
	}

	// A PBM file header is "P4", followed by whitespace, followed by a
	// width in pixels, followed by whitespace, followed by a height in
	// pixels, followed by a single whitespace.
	if nr.GetNextByteAsRune() != 'P' || nr.GetNextByteAsRune() != '4' || !unicode.IsSpace(nr.GetNextByteAsRune()) {
		return badness()
	}
	var cfg image.Config
	cfg.Width = nr.GetNextInt()
	cfg.Height = nr.GetNextInt()
	if nr.Err() != nil || !unicode.IsSpace(nr.GetNextByteAsRune()) {
		return badness()
	}

	// A PBM file's color map is 0=white, 1=black.
	colorMap := make(color.Palette, 2)
	colorMap[0] = color.RGBA{255, 255, 255, 255}
	colorMap[1] = color.RGBA{0, 0, 0, 255}
	cfg.ColorModel = colorMap
	return cfg, nil
}

// Decode reads a complete PBM image.
func Decode(r io.Reader) (image.Image, error) {
	// Read the image header and use it to prepare a paletted image.
	header, err := DecodeConfig(r)
	if err != nil {
		return nil, err
	}
	img := image.NewPaletted(image.Rect(0, 0, header.Width, header.Height), header.ColorModel.(color.Palette))

	// Read bits until no more remain.
	buf := make([]byte, 1<<20) // Arbitrary, large, buffer size
	bitsRemaining := header.Width * header.Height
	bitNum := 0
ReadLoop:
	for {
		var nRead int
		nRead, err = r.Read(buf)
		if nRead == 0 && err != nil {
			return nil, err
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
				if bitNum%header.Width == 0 {
					// Ignore row padding.
					break
				}
			}
		}
	}
	return img, nil
}

// Indicate that we can decode PBM files.
func init() {
	image.RegisterFormat("pbm", "P4", Decode, DecodeConfig)
}
