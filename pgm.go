// This file provides image support for both "raw" (binary) and
// "plain" (ASCII) Portable GrayMap (PGM) files.

package netpbm

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"io"
)

// A Gray is an image.Gray that knows its maximum value.
type Gray struct {
	*image.Gray       // Grayscale image representation
	Maxval      uint8 // Value representing 100% white
}

// NewGray returns a new Gray with the given bounds and maximum value.
func NewGray(r image.Rectangle, m uint8) *Gray {
	return &Gray{Gray: image.NewGray(r), Maxval: m}
}

// A Gray16 is an image.Gray16 that knows its maximum value.
type Gray16 struct {
	*image.Gray16        // Grayscale image representation
	Maxval        uint16 // Value representing 100% white
}

// NewGray16 returns a new Gray16 with the given bounds and maximum value.
func NewGray16(r image.Rectangle, m uint16) *Gray16 {
	return &Gray16{Gray16: image.NewGray16(r), Maxval: m}
}

// decodeConfigPGM reads and parses a PGM header, either "raw" (binary) or
// "plain" (ASCII).
func decodeConfigPGM(r io.Reader) (image.Config, error) {
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
		return image.Config{}, err
	}

	// Define the color model using the gray channel's maximum value.
	if header.Maxval < 256 {
		header.Model = color.GrayModel
	} else {
		header.Model = color.Gray16Model
	}

	// Store and return the image configuration.
	var cfg image.Config
	cfg.Width = header.Width
	cfg.Height = header.Height
	cfg.ColorModel = header
	return cfg, nil
}

// decodePGM reads a complete "raw" (binary) PGM image.
func decodePGM(r io.Reader) (image.Image, error) {
	// Read the image header, and use it to prepare a grayscale image.
	br := bufio.NewReader(r)
	config, err := decodeConfigPGM(br)
	if err != nil {
		return nil, err
	}

	// Create either a Gray or a Gray16 image.
	var img image.Image                               // Image to return
	var data []uint8                                  // Image data
	maxVal := config.ColorModel.(netpbmHeader).Maxval // 100% white value
	if maxVal < 256 {
		gray := image.NewGray(image.Rect(0, 0, config.Width, config.Height))
		data = gray.Pix
		img = Gray{Gray: gray, Maxval: uint8(maxVal)}
	} else {
		gray16 := image.NewGray16(image.Rect(0, 0, config.Width, config.Height))
		data = gray16.Pix
		img = Gray16{Gray16: gray16, Maxval: uint16(maxVal)}
	}

	// Raw PGM images are nice because we can read directly into the image
	// data.
	for len(data) > 0 {
		nRead, err := br.Read(data)
		if err != nil && err != io.EOF {
			return img, err
		}
		if nRead == 0 {
			return img, errors.New("Failed to read binary PGM data")
		}
		data = data[nRead:]
	}
	return img, nil
}

// decodePGMPlain reads a complete "plain" (ASCII) PGM image.
func decodePGMPlain(r io.Reader) (image.Image, error) {
	// Read the image header, and use it to prepare a grayscale image.
	br := bufio.NewReader(r)
	config, err := decodeConfigPGM(br)
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
			err = errors.New("Failed to parse ASCII PGM data")
		}
		return img, err
	}

	// Create either a Gray or a Gray16 image.
	var data []uint8                                  // Image data
	maxVal := config.ColorModel.(netpbmHeader).Maxval // 100% white value
	if maxVal < 256 {
		gray := image.NewGray(image.Rect(0, 0, config.Width, config.Height))
		data = gray.Pix
		img = Gray{Gray: gray, Maxval: uint8(maxVal)}
	} else {
		gray16 := image.NewGray16(image.Rect(0, 0, config.Width, config.Height))
		data = gray16.Pix
		img = Gray16{Gray16: gray16, Maxval: uint16(maxVal)}
	}

	// Read ASCII base-10 integers until no more remain.
	totalVals := config.Width * config.Height
	for i := 0; i < totalVals; {
		val := nr.GetNextInt()
		switch {
		case nr.Err() != nil:
			return badness()
		case val < 0 || val > maxVal:
			return badness()
		case maxVal < 256:
			data[i] = uint8(val)
			i++
		case maxVal < 65536:
			data[i] = uint8(val >> 8)
			data[i+1] = uint8(val)
			i += 2
		default:
			return badness()
		}
	}
	return img, nil
}

// Indicate that we can decode both raw and plain PGM files.
func init() {
	image.RegisterFormat("pgm", "P5", decodePGM, decodeConfigPGM)
	image.RegisterFormat("pgm", "P2", decodePGMPlain, decodeConfigPGM)
}
