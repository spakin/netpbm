// This file provides image support for both "raw" (binary) and
// "plain" (ASCII) Portable PixMap (PPM) files.

package netpbm

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"io"
)

// A Color is an image.NRGBA that knows its maximum value.
type Color struct {
	*image.NRGBA       // Color image representation
	Maxval       uint8 // Value representing 100% white
}

// NewColor returns a new 32-bit color image with the given bounds and maximum
// value.
func NewColor(r image.Rectangle, m uint8) *Color {
	return &Color{NRGBA: image.NewNRGBA(r), Maxval: m}
}

// A Color64 is an image.NRGBA64 that knows its maximum value.
type Color64 struct {
	*image.NRGBA64        // Color image representation
	Maxval         uint16 // Value representing 100% white
}

// NewColor64 returns a new 64-bit color image with the given bounds and
// maximum value.
func NewColor64(r image.Rectangle, m uint16) *Color64 {
	return &Color64{NRGBA64: image.NewNRGBA64(r), Maxval: m}
}

// decodeConfigPPM reads and parses a PPM header, either "raw" (binary) or
// "plain" (ASCII).
func decodeConfigPPM(r io.Reader) (image.Config, error) {
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
		return image.Config{}, err
	}

	// Define the color model using the color channel's maximum value.
	if header.Maxval < 256 {
		header.Model = color.NRGBAModel
	} else {
		header.Model = color.NRGBA64Model
	}

	// Store and return the image configuration.
	var cfg image.Config
	cfg.Width = header.Width
	cfg.Height = header.Height
	cfg.ColorModel = header
	return cfg, nil
}

// decodePPM reads a complete "raw" (binary) PPM image.
func decodePPM(r io.Reader) (image.Image, error) {
	// Read the image header, and use it to prepare a color image.
	br := bufio.NewReader(r)
	config, err := decodeConfigPPM(br)
	if err != nil {
		return nil, err
	}

	// Create either a Color or a Color64 image.
	var img image.Image // Image to return
	var data []uint8    // RGBA image data
	var rgbData []uint8 // RGB file data
	nPixels := config.Width * config.Height
	maxVal := config.ColorModel.(netpbmHeader).Maxval // 100% white value
	if maxVal < 256 {
		rgb := image.NewNRGBA(image.Rect(0, 0, config.Width, config.Height))
		data = rgb.Pix
		rgbData = make([]uint8, nPixels*3)
		img = Color{NRGBA: rgb, Maxval: uint8(maxVal)}
	} else {
		rgb := image.NewNRGBA64(image.Rect(0, 0, config.Width, config.Height))
		data = rgb.Pix
		rgbData = make([]uint8, nPixels*3*2)
		img = Color64{NRGBA64: rgb, Maxval: uint16(maxVal)}
	}

	// Read RGB (no A) data into a holding buffer.
	rgbDataLeft := rgbData // RGB data left to read
	for len(rgbDataLeft) > 0 {
		nRead, err := br.Read(rgbDataLeft)
		if err != nil && err != io.EOF {
			return img, err
		}
		if nRead == 0 {
			return img, errors.New("Failed to read binary PPM data")
		}
		rgbDataLeft = rgbDataLeft[nRead:]
	}

	// Spread out RGB data into RGBA.
	nCopy := 3                       // Copy this many bytes from the input...
	nAlpha := 1                      // ...then generate this many bytes of alpha.
	opaque := []uint8{uint8(maxVal)} // Maximum opacity
	if maxVal >= 256 {
		nCopy *= 2
		nAlpha *= 2
		opaque = []uint8{uint8(maxVal >> 8), uint8(maxVal)}
	}
	for p, s, d := 0, 0, 0; p < nPixels; p++ {
		copy(data[d:d+nCopy], rgbData[s:s+nCopy])
		copy(data[d+nCopy:d+nCopy+nAlpha], opaque)
		s += nCopy
		d += nCopy + nAlpha
	}
	return img, nil
}

// decodePPMPlain reads a complete "plain" (ASCII) PPM image.
func decodePPMPlain(r io.Reader) (image.Image, error) {
	// Read the image header, and use it to prepare a color image.
	br := bufio.NewReader(r)
	config, err := decodeConfigPPM(br)
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
			err = errors.New("Failed to parse ASCII PPM data")
		}
		return img, err
	}

	// Create either a Color or a Color64 image.
	var data []uint8                                  // Image data
	maxVal := config.ColorModel.(netpbmHeader).Maxval // 100% white value
	if maxVal < 256 {
		rgb := image.NewNRGBA(image.Rect(0, 0, config.Width, config.Height))
		data = rgb.Pix
		img = Color{NRGBA: rgb, Maxval: uint8(maxVal)}
	} else {
		rgb := image.NewNRGBA64(image.Rect(0, 0, config.Width, config.Height))
		data = rgb.Pix
		img = Color64{NRGBA64: rgb, Maxval: uint16(maxVal)}
	}

	// Read ASCII base-10 integers until no more remain.
	if maxVal < 256 {
		for i := 0; i < len(data); {
			for d := 0; d < 3; d++ {
				val := nr.GetNextInt()
				switch {
				case nr.Err() != nil:
					return badness()
				case val < 0 || val > maxVal:
					return badness()
				default:
					data[i] = uint8(val)
					i++
				}
			}
			data[i] = uint8(maxVal)
			i++
		}
	} else {
		for i := 0; i < len(data); {
			for d := 0; d < 3; d++ {
				val := nr.GetNextInt()
				switch {
				case nr.Err() != nil:
					return badness()
				case val < 0 || val > maxVal:
					return badness()
				default:
					data[i] = uint8(val >> 8)
					data[i+1] = uint8(val)
					i += 2
				}
			}
			data[i] = uint8(maxVal >> 8)
			data[i+1] = uint8(maxVal)
			i += 2
		}
	}
	return img, nil
}

// Indicate that we can decode both raw and plain PPM files.
func init() {
	image.RegisterFormat("ppm", "P6", decodePPM, decodeConfigPPM)
	image.RegisterFormat("ppm", "P3", decodePPMPlain, decodeConfigPPM)
}
