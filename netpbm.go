/*

Package netpbm implements image decoders and encoders for the Netpbm image
formats: PBM (black and white), PGM (grayscale), PPM (color), and PAM (black
and white, grayscale, or color, as indicated by the image header).  Both "raw"
(binary) and "plain" (ASCII) files can be read and written.  Both 8-bit and
16-bit color channels are supported.

The netpbm package is fully compatible with the image package in the standard
library but additionally reproduces the Netpbm library's ability to promote
formats during decode.  That is, a program that expects to read a grayscale
image can also be given a black-and-white image, and a program that expects to
read a color image can also be given either a grayscale or a black-and-white
image.

The Netpbm home page is at http://netpbm.sourceforge.net/.

*/
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

// GetLineAsKeyValue returns the next line, split into a space-separated key
// and a value or nil on error.  Errors are sticky.
func (nr *netpbmReader) GetLineAsKeyValue() []string {
	// Read a line.
	if nr.err != nil {
		return nil
	}
	var s string
	s, nr.err = nr.ReadString('\n')
	if nr.err != nil {
		return nil
	}

	// Split the string into a key and a value.  As a special case "#"
	// counts as a key, and everything following it is a comment.
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if s[0] == '#' {
		return []string{"#", strings.TrimSpace(s[1:])}
	}
	fs := strings.SplitN(s, " ", 2)
	if len(fs) == 1 {
		return []string{fs[0], ""}
	}
	return []string{fs[0], strings.TrimSpace(fs[1])}
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
	if nr.err != nil && nr.err != io.EOF {
		return -1
	}
	nr.err = nr.UnreadByte()
	if nr.err != nil {
		return -1
	}
	return value
}

// GetNextString returns the next string from the Netpbm file.
func (nr *netpbmReader) GetNextString() string {
	var c rune
	for nr.err == nil && !unicode.IsLetter(c) {
		for c = nr.GetNextByteAsRune(); unicode.IsSpace(c); c = nr.GetNextByteAsRune() {
		}
		if c == '#' {
			// Comment -- discard the rest of the line.
			for c = nr.GetNextByteAsRune(); c != '\n'; c = nr.GetNextByteAsRune() {
			}
		}
	}
	if nr.err != nil {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteRune(c)

	for c = nr.GetNextByteAsRune(); unicode.IsLetter(c) || unicode.IsPunct(c); c = nr.GetNextByteAsRune() {
		_, nr.err = sb.WriteRune(c)
	}
	if nr.err != nil && nr.err != io.EOF {
		return ""
	}
	nr.err = nr.UnreadByte()
	if nr.err != nil {
		return ""
	}
	return sb.String()
}

// GetIntsAndComments returns n integers and a list of comments encountered
// along the way.  Comments discard the initial "#" and up to one subsequent
// whitespace character as well as the final carriage return and/or line feed.
func (nr *netpbmReader) GetIntsAndComments(n int) ([]int, []string, error) {
	// Initialize a state machine.
	num := 0                         // Current number
	numbers := make([]int, 0, n)     // All numbers
	cmt := make([]rune, 0, 100)      // Current comment
	comments := make([]string, 0, 1) // All strings
	const (
		InSpace = iota
		InDigit
		InComment
	)
	state := InSpace  // Current state
	var prevState int // Previous state (restored after a comment)

	// Process one rune at a time.
	var c rune
RuneLoop:
	for {
		switch state {
		case InSpace:
			// Skip spaces.
			for c = nr.GetNextByteAsRune(); unicode.IsSpace(c); c = nr.GetNextByteAsRune() {
			}
			if c >= '0' && c <= '9' {
				state = InDigit
				num = int(c - '0')
			} else if c == '#' {
				state = InComment
				prevState = InSpace
			} else if c == 0 {
				return nil, nil, errors.New("Unexpected EOF in Netpbm header")
			} else {
				return nil, nil, fmt.Errorf("Unexpected header character %q", c)
			}

		case InDigit:
			// Read a base-10 number.
			for c = nr.GetNextByteAsRune(); c >= '0' && c <= '9'; c = nr.GetNextByteAsRune() {
				num = num*10 + int(c-'0')
			}
			if unicode.IsSpace(c) {
				state = InSpace
				numbers = append(numbers, num)
				if len(numbers) == n {
					return numbers, comments, nil
				}
			} else if c == '#' {
				state = InComment
				prevState = InDigit
			} else if c == 0 {
				return nil, nil, errors.New("Unexpected EOF in Netpbm header")
			} else {
				return nil, nil, fmt.Errorf("Unexpected header character %q", c)
			}

		case InComment:
			// Append to the current comment until we reach the end
			// of the line.
			for c = nr.GetNextByteAsRune(); c != '\n' && c != '\r'; c = nr.GetNextByteAsRune() {
				cmt = append(cmt, c)
			}
			if len(cmt) > 0 && unicode.IsSpace(cmt[0]) {
				cmt = cmt[1:]
			}
			comments = append(comments, string(cmt))
			cmt = cmt[:0]
			state = prevState
			prevState = InComment

		default:
			break RuneLoop
		}
	}
	return nil, nil, errors.New("Unexpected EOF in Netpbm header")
}

// GetASCIIData reads ASCII base-10 integers until the input array is filled.
// It returns a success code.
func (nr *netpbmReader) GetASCIIData(maxVal int, data []uint8) bool {
	// Read ASCII base-10 integers until no more remain.
	if maxVal < 256 {

		for i := 0; i < len(data); i++ {
			val := nr.GetNextInt()
			switch {
			case nr.Err() != nil:
				return false
			case val < 0 || val > maxVal:
				return false
			default:
				data[i] = uint8(val)
			}
		}
	} else {
		for i := 0; i < len(data); i += 2 {
			val := nr.GetNextInt()
			switch {
			case nr.Err() != nil:
				return false
			case val < 0 || val > maxVal:
				return false
			default:
				data[i] = uint8(val >> 8)
				data[i+1] = uint8(val)
			}
		}
	}
	return true
}

// A netpbmHeader encapsulates the components of an image header.
type netpbmHeader struct {
	Magic     string   // Two-character magic value (e.g., "P6" for PPM)
	Width     int      // Image width in pixels
	Height    int      // Image height in pixels
	Depth     int      // Image pixel depth in bytes
	Maxval    int      // Maximum channel value (0-65535)
	TupleType string   // Image tuple type ("RGB_ALPHA", etc.)
	Comments  []string // Aggregated list of comment lines
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

	// PBM files (raw or plain) don't specify a maximum channel.  All other
	// formats do.
	switch header.Magic {
	case "P1", "P4":
		header.Maxval = 1
		nums, comments, err := nr.GetIntsAndComments(2)
		if err != nil {
			return netpbmHeader{}, false
		}
		header.Width = nums[0]
		header.Height = nums[1]
		header.Comments = comments

	default:
		nums, comments, err := nr.GetIntsAndComments(3)
		if err != nil {
			return netpbmHeader{}, false
		}
		header.Width = nums[0]
		header.Height = nums[1]
		header.Maxval = nums[2]
		header.Comments = comments
	}
	if nr.Err() != nil || header.Maxval < 1 || header.Maxval > 65535 {
		return netpbmHeader{}, false
	}

	// Return the header and a success code.
	return header, true
}

// An Image extends image.Image to include a few extra methods.
type Image interface {
	image.Image                             // At, Bounds, and ColorModel
	MaxValue() uint16                       // Maximum value on each color channel
	Format() Format                         // Netpbm format
	Opaque() bool                           // Report whether the image is fully opaque
	PixOffset(x, y int) int                 // Find (x, y) in pixel data
	Set(x, y int, c color.Color)            // Set a pixel to a color
	SubImage(r image.Rectangle) image.Image // Portion of the image visible through r
}

// A Format represents a specific Netpbm format.
type Format int

// Define a symbol for each supported Netpbm format.
const (
	PNM Format = iota // Portable Any Map (any of PBM, PGM, or PPM)
	PBM               // Portable Bit Map (black and white)
	PGM               // Portable Gray Map (grayscale)
	PPM               // Portable Pix Map (color)
	PAM               // Portable Arbitrary Map (alpha)
)

// String outputs the name of a Netpbm format.
func (f Format) String() string {
	switch f {
	case PNM:
		return "PNM"
	case PBM:
		return "PBM"
	case PGM:
		return "PGM"
	case PPM:
		return "PPM"
	case PAM:
		return "PAM"
	default:
		return fmt.Sprintf("%%!s(netpbm.Format=%d)", f)
	}
}

// DecodeOptions represents a list of options for decoding a Netpbm file.
type DecodeOptions struct {
	Target      Format // Netpbm format to return
	Exact       bool   // true=allow only Target; false=promote lesser formats
	PBMMaxValue uint16 // Maximum channel value to use when promoting a PBM image (0=default)
}

// DecodeConfigWithComments returns image metadata without decoding the entire
// image.  Unlike Decode, it also returns any comments appearing in the file.
// Pass in a bufio.Reader if you intend to read data following the image
// header.
func DecodeConfigWithComments(r io.Reader) (image.Config, []string, error) {
	// Peek at the file's magic number.
	rr, ok := r.(*bufio.Reader)
	if !ok {
		rr = bufio.NewReader(r)
	}
	magic, err := rr.Peek(2)
	if err != nil {
		return image.Config{}, nil, err
	}

	// Invoke the decode function corresponding to the magic number.
	if magic[0] != 'P' {
		return image.Config{}, nil, errors.New("Not a Netpbm image")
	}
	switch magic[1] {
	case '1', '4':
		// PBM
		return decodeConfigPBMWithComments(rr)
	case '2', '5':
		// PGM
		return decodeConfigPGMWithComments(rr)
	case '3', '6':
		// PPM
		return decodeConfigPPMWithComments(rr)
	case '7':
		// PAM
		return decodeConfigPAMWithComments(rr)
	default:
		// None of the above
		return image.Config{}, nil, fmt.Errorf("Unrecognized magic sequence %q", string(magic))
	}
}

// DecodeConfig returns image metadata without decoding the entire image.  Pass
// in a bufio.Reader if you intend to read data following the image header.
func DecodeConfig(r io.Reader) (image.Config, error) {
	cfg, _, err := DecodeConfigWithComments(r)
	return cfg, err
}

// DecodeWithComments reads a Netpbm image from r and returns it as an Image.
// Unlike Decode, it also returns any comments appearing in the file.  Pass in
// a bufio.Reader if you intend to read data following the image.
func DecodeWithComments(r io.Reader, opts *DecodeOptions) (Image, []string, error) {
	// Peek at the file's magic number.
	rr, ok := r.(*bufio.Reader)
	if !ok {
		rr = bufio.NewReader(r)
	}
	magic, err := rr.Peek(2)
	if err != nil {
		return nil, nil, err
	}
	if magic[0] != 'P' {
		return nil, nil, errors.New("Not a Netpbm image")
	}

	// Provide default options.
	var o DecodeOptions
	if opts != nil {
		o = *opts
	}
	if o.PBMMaxValue == 0 {
		o.PBMMaxValue = 255
	}
	if o.Exact && o.Target == PNM {
		// PNM isn't its own format so it doesn't make sense to try to
		// read exactly a PNM file.
		return nil, nil, errors.New("Exact=true is incompatible with Target=PNM")
	}

	// Invoke the decode function corresponding to the magic number.
	var img image.Image   // Image to return
	var comments []string // Comments appearing in the image header
	switch magic[1] {
	case '1':
		// Plain PBM
		if o.Exact && o.Target != PBM {
			return nil, nil, errors.New("PBM rejected by Decode options")
		}
		img, comments, err = decodePBMPlainWithComments(rr)
	case '2':
		// Plain PGM
		if o.Exact && o.Target != PGM {
			return nil, nil, errors.New("PGM rejected by Decode options")
		}
		img, comments, err = decodePGMPlainWithComments(rr)
	case '3':
		// Plain PPM
		if o.Exact && o.Target != PPM {
			return nil, nil, errors.New("PPM rejected by Decode options")
		}
		img, comments, err = decodePPMPlainWithComments(rr)
	case '4':
		// Raw PBM
		if o.Exact && o.Target != PBM {
			return nil, nil, errors.New("PBM rejected by Decode options")
		}
		img, comments, err = decodePBMWithComments(rr)
	case '5':
		// Raw PGM
		if o.Exact && o.Target != PGM {
			return nil, nil, errors.New("PGM rejected by Decode options")
		}
		img, comments, err = decodePGMWithComments(rr)
	case '6':
		// Raw PPM
		if o.Exact && o.Target != PPM {
			return nil, nil, errors.New("PPM rejected by Decode options")
		}
		img, comments, err = decodePPMWithComments(rr)
	case '7':
		// Raw PAM
		if o.Exact && o.Target != PAM {
			return nil, nil, errors.New("PAM rejected by Decode options")
		}
		img, comments, err = decodePAMWithComments(rr)
	default:
		// None of the above
		return nil, nil, fmt.Errorf("Unrecognized magic sequence %q", string(magic))
	}
	if err != nil {
		return nil, nil, err
	}

	// A PNM target accepts any of PBM, PGM, or PPM as is.
	nimg := img.(Image)
	if o.Target == PNM {
		return nimg, comments, nil
	}

	// If requested, promote the image to a richer format.  We've already
	// rejected the case of a mismatch when mismatches are forbidden.
	if nimg.Format() > o.Target {
		return nil, nil, fmt.Errorf("Cannot demote a %s image to a %s image", nimg.Format(), o.Target)
	}
	for nimg.Format() < o.Target {
		switch nimg.Format() {
		case PBM:
			mVal := o.PBMMaxValue
			if mVal < 256 {
				nimg = nimg.(*BW).PromoteToGrayM(uint8(mVal))
			} else {
				nimg = nimg.(*BW).PromoteToGrayM32(mVal)
			}
		case PGM:
			if nimg.MaxValue() < 256 {
				nimg = nimg.(*GrayM).PromoteToRGBM()
			} else {
				nimg = nimg.(*GrayM32).PromoteToRGBM64()
			}
		default:
			panic("Attempted to promote a format other than PBM or PGM")
		}
	}
	return nimg, comments, nil
}

// Decode reads a Netpbm image from r and returns it as an Image.  a
// bufio.Reader if you intend to read data following the image.
func Decode(r io.Reader, opts *DecodeOptions) (Image, error) {
	img, _, err := DecodeWithComments(r, opts)
	return img, err
}

// EncodeOptions represents a list of options for writing a Netpbm file.
type EncodeOptions struct {
	Format    Format   // Netpbm format
	MaxValue  uint16   // Maximum value for each color channel (ignored for PBM)
	Plain     bool     // true="plain" (ASCII); false="raw" (binary)
	TupleType string   // Image tuple type for a PAM image (RGB_ALPHA, etc.)
	Comments  []string // Header comments, with no leading "#" or trailing newlines
}

// Encode writes an arbitrary image in any of the Netpbm formats.  Given an
// opts.Format of PNM, use the image's Format if img is a Netpbm image or PPM
// if not.  Given an opts.MaxValue of 0, use the image's MaxValue if img is a
// Netpbm image or 255 if not.  Given a nil opts, assign Format as if it were
// PNM and MaxValue as if it were 0.
func Encode(w io.Writer, img image.Image, opts *EncodeOptions) error {
	// Start by copying opts if provided or initializing a new set of
	// EncodeOptions if not.
	var o EncodeOptions
	if opts != nil {
		o = *opts
	}

	// If Format is PNM (the zero value), replace it with an intelligently
	// selected Netpbm format.
	if o.Format == PNM {
		switch img := img.(type) {
		case Image:
			o.Format = img.Format()
		default:
			o.Format = PAM
		}
	}

	// If Format is PAM and TupleType is not specified, default to
	// RGB_ALPHA.
	if o.Format == PAM && o.TupleType == "" {
		o.TupleType = "RGB_ALPHA"
	}

	// If MaxValue is 0, replace it with an intelligently selected maximum
	// value.
	if o.MaxValue == 0 {
		switch img := img.(type) {
		case Image:
			o.MaxValue = img.MaxValue()
		default:
			o.MaxValue = 255
		}
	}

	// Encode the image using the specified format and options.
	switch o.Format {
	case PPM:
		return encodePPM(w, img, &o)
	case PGM:
		return encodePGM(w, img, &o)
	case PBM:
		return encodePBM(w, img, &o)
	case PAM:
		return encodePAM(w, img, &o)
	default:
		return fmt.Errorf("Invalid Netpbm format specified (%s)", o.Format)
	}
}

// writePlainData writes numbers read from a channel as base-10 strings, at
// most 70 characters per line.
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
