/*

	Package netpbm implements image decoders and encoders for the
	Netpbm image formats.

	The Netpbm home page is at http://netpbm.sourceforge.net/.
*/
package netpbm

import (
	"bufio"
	"image/color"
	"unicode"
)

// A netpbmReader extends bufio.Reader with the ability to read bytes
// and numbers while skipping over comments.
type netpbmReader struct {
	*bufio.Reader       // Inherit Read, UnreadByte, etc.
	err           error // Sticky error state
}

// newNetpbmReader allocates, initializes, and returns a new netpbmReader.
func newNetpbmReader(r *bufio.Reader) *netpbmReader {
	return &netpbmReader{Reader: r}
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
	var b byte
	b, nr.err = nr.ReadByte()
	if nr.err != nil {
		return 0
	}
	return rune(b)
}

// GetNextInt returns the next base-10 integer read from a netpbmReader,
// skipping preceding whitespace and comments.
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

// A netpbmHeader encapsulates the components of an image header.
type netpbmHeader struct {
	Magic  string      // Two-character magic value (e.g., "P6" for PPM)
	Width  int         // Image width in pixels
	Height int         // Image height in pixels
	Maxval int         // Maximum channel value (0-65535)
	Model  color.Model // Color model represented by this image
}

// We let netpbmHeader implement color.Model.  This lets us piggyback all of
// our image metadata into an image.Config.
func (nh netpbmHeader) Convert(c color.Color) color.Color {
	return nh.Model.Convert(c)
}

// GetNetpbmHeader parses the entire header (PBM, PGM, or PPM; raw or
// plain) and returns it as a netpbmHeader (plus a success value).
func (nr *netpbmReader) GetNetpbmHeader() (netpbmHeader, bool) {
	var header netpbmHeader

	// Read the magic value and skip the following whitespace.
	rune1 := nr.GetNextByteAsRune()
	if rune1 != 'P' {
		return netpbmHeader{}, false
	}
	rune2 := nr.GetNextByteAsRune()
	if rune2 < '1' || rune2 > '6' {
		return netpbmHeader{}, false
	}
	if !unicode.IsSpace(nr.GetNextByteAsRune()) {
		return netpbmHeader{}, false
	}
	header.Magic = string(rune1) + string(rune2)

	// Read the width and height.
	header.Width = nr.GetNextInt()
	header.Height = nr.GetNextInt()

	// PBM files (raw or plain) don't specify a maximum channel.  All other
	// formats do.
	switch header.Magic {
	case "P1", "P4":
		header.Maxval = 1
	default:
		header.Maxval = nr.GetNextInt()
	}
	if nr.Err() != nil || !unicode.IsSpace(nr.GetNextByteAsRune()) ||
		header.Maxval < 1 || header.Maxval > 65535 {
		return netpbmHeader{}, false
	}

	// Return the header and a success code.
	return header, true
}
