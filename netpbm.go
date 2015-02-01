/*

	Package netpbm implements image decoders and encoders for the Netpbm
	image formats (PBM, PGM, and PPM), in both "raw" (binary) and "plain"
	(ASCII) versions.

	The Netpbm home page is at http://netpbm.sourceforge.net/.
*/
package netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"io"
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
	Magic  string // Two-character magic value (e.g., "P6" for PPM)
	Width  int    // Image width in pixels
	Height int    // Image height in pixels
	Maxval int    // Maximum channel value (0-65535)
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

// An Image extends image.Image to include a MaxValue method.
type Image interface {
	image.Image
	MaxValue() uint16
}

// A Format represents a specific Netpbm format.
type Format int

// Define a symbol for each supported Netpbm format.
const (
	PBM Format = 1 << iota
	PGM
	PPM
)

// DecodeOptions represents a list of options for decoding a Netpbm file.
type DecodeOptions struct {
	Allowed Format // Bit mask of Netpbm formats allowed (0 = all)
}

// DecodeConfig returns image metadata without decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	// Peek at the file's magic number.
	rr, ok := r.(*bufio.Reader)
	if !ok {
		rr = bufio.NewReader(r)
	}
	magic, err := rr.Peek(2)
	if err != nil {
		return image.Config{}, err
	}

	// Invoke the decode function corresponding to the magic number.
	if magic[0] != 'P' {
		return image.Config{}, errors.New("Not a Netpbm image")
	}
	switch magic[1] {
	case '1', '4':
		// PBM
		return decodeConfigPBM(rr)
	case '2', '5':
		// PGM
		return decodeConfigPGM(rr)
	case '3', '6':
		// PPM
		return decodeConfigPPM(rr)
	default:
		// None of the above
		return image.Config{}, fmt.Errorf("Unrecognized magic sequence %q", string(magic))
	}
}

// Decode reads a Netpbm image from r and returns it as an Image.
func Decode(r io.Reader, opt *DecodeOptions) (Image, error) {
	// Determine the set of all formats allowed.
	var allowed Format
	if opt == nil || opt.Allowed == 0 {
		allowed = PBM | PGM | PPM
	} else {
		allowed = opt.Allowed
	}

	// Peek at the file's magic number.
	rr, ok := r.(*bufio.Reader)
	if !ok {
		rr = bufio.NewReader(r)
	}
	magic, err := rr.Peek(2)
	if err != nil {
		return nil, err
	}

	// Invoke the decode function corresponding to the magic number.
	if magic[0] != 'P' {
		return nil, errors.New("Not a Netpbm image")
	}
	var img image.Image // Image to return
	switch magic[1] {
	case '1':
		// Plain PBM
		if allowed&PBM != PBM {
			return nil, errors.New("PBM rejected by Decode options")
		}
		img, err = decodePBMPlain(rr)
	case '2':
		// Plain PGM
		if allowed&PGM != PGM {
			return nil, errors.New("PGM rejected by Decode options")
		}
		img, err = decodePGMPlain(rr)
	case '3':
		// Plain PPM
		if allowed&PPM != PPM {
			return nil, errors.New("PPM rejected by Decode options")
		}
		img, err = decodePPMPlain(rr)
	case '4':
		// Raw PBM
		if allowed&PBM != PBM {
			return nil, errors.New("PBM rejected by Decode options")
		}
		img, err = decodePBM(rr)
	case '5':
		// Raw PGM
		if allowed&PGM != PGM {
			return nil, errors.New("PGM rejected by Decode options")
		}
		img, err = decodePGM(rr)
	case '6':
		// Raw PPM
		if allowed&PPM != PPM {
			return nil, errors.New("PPM rejected by Decode options")
		}
		img, err = decodePPM(rr)
	default:
		// None of the above
		return nil, fmt.Errorf("Unrecognized magic sequence %q", string(magic))
	}
	return img.(Image), err
}

// EncodeOptions represents a list of options for writing a Netpbm file.
type EncodeOptions struct {
	Format   Format // Netpbm format
	MaxValue uint16 // Maximum value for each color channel
	Plain    bool   // true="plain" (ASCII); false="raw" (binary)
	Comment  string // Header comment
}

// Encode writes an arbitrary image in any of the Netpbm formats.
func Encode(w io.Writer, img image.Image, opts *EncodeOptions) error {
	var o EncodeOptions
	if opts == nil {
		// Select some reasonable default options.
		o = EncodeOptions{
			Format:   PPM,
			MaxValue: 255,
		}
	} else {
		// Ensure the provided options are sensible.
		o = *opts
		if o.MaxValue == 0 {
			return errors.New("MaxValue must be greater than 0")
		}
	}
	switch o.Format {
	case PPM:
		return encodePPM(w, img, &o)
	case PGM:
		return encodePGM(w, img, &o)
	case PBM:
		return encodePBM(w, img, &o)
	default:
		return errors.New("Invalid Netpbm format specified")
	}
}

// writePlainData writes numbers read from a channel as base-10 strings, 70
// characters per line.
func writePlainData(w io.Writer, ch chan uint16) error {
	var line string
	for s := range ch {
		word := fmt.Sprintf("%d ", s)
		if len(line)+len(word) <= 70 {
			line += word
		} else {
			lineBytes := []byte(line)
			lineBytes[len(lineBytes)-1] = '\n'
			_, err := w.Write(lineBytes)
			if err != nil {
				return err
			}
			line = word
		}

	}
	if line != "" {
		lineBytes := []byte(line)
		lineBytes[len(lineBytes)-1] = '\n'
		_, err := w.Write(lineBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeRawData writes numbers read from a channel as bytes, either uint8s (if
// wd = 1) or uint16s (if wd = 2).
func writeRawData(w io.Writer, ch chan uint16, wd int) error {
	var err error
	wb, ok := w.(*bufio.Writer)
	if !ok {
		wb = bufio.NewWriter(w)
	}
	switch wd {
	case 1:
		for s := range ch {
			err = wb.WriteByte(uint8(s))
			if err != nil {
				return err
			}
		}
	case 2:
		for s := range ch {
			err = wb.WriteByte(uint8(s >> 8))
			if err != nil {
				return err
			}
			err = wb.WriteByte(uint8(s))
			if err != nil {
				return err
			}
		}
	default:
		panic("writeRawData was given an invalid byte width")
	}
	wb.Flush()
	return nil
}
