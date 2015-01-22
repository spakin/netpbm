// This file provides image support for both "raw" (binary) and
// "plain" (ASCII) Portable BitMap (PBM) files.

package netpbm

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"io"
	"unicode"
)

// decodeConfigPBM reads and parses a PBM header, either "raw" (binary) or
// "plain" (ASCII).
func decodeConfigPBM(r io.Reader) (image.Config, error) {
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
		return image.Config{}, err
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
	return cfg, nil
}

// decodePBM reads a complete "raw" (binary) PBM image.
func decodePBM(r io.Reader) (image.Image, error) {
	// Read the image header, and use it to prepare a paletted image.
	br := bufio.NewReader(r)
	config, err := decodeConfigPBM(br)
	if err != nil {
		return nil, err
	}
	img := image.NewPaletted(image.Rect(0, 0, config.Width, config.Height), config.ColorModel.(color.Palette))

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
				if bitNum%config.Width == 0 {
					// Ignore row padding.
					break
				}
			}
		}
	}
	return img, nil
}

// decodePBMPlain reads a complete "plain" (ASCII) PBM image.
func decodePBMPlain(r io.Reader) (image.Image, error) {
	// Read the image header, and use it to prepare a paletted image.
	br := bufio.NewReader(r)
	config, err := decodeConfigPBM(br)
	if err != nil {
		return nil, err
	}
	img := image.NewPaletted(image.Rect(0, 0, config.Width, config.Height), config.ColorModel.(color.Palette))

	// Define a simple error handler.
	nr := newNetpbmReader(br)
	badness := func() (image.Image, error) {
		// Something went wrong.  Either we have an error code to
		// explain what or we make up a generic error message.
		err := nr.Err()
		if err == nil {
			err = errors.New("Failed to parse ASCII PBM data")
		}
		return img, err
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
	return img, nil
}

// Indicate that we can decode both raw and plain PBM files.
func init() {
	image.RegisterFormat("pbm", "P4", decodePBM, decodeConfigPBM)
	image.RegisterFormat("pbm", "P1", decodePBMPlain, decodeConfigPBM)
}
