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
	// Read the image header and use it to prepare a grayscale image.
	br := bufio.NewReader(r)
	config, err := decodeConfigPGM(br)
	if err != nil {
		return nil, err
	}

	// Raw PGM images are nice because we can read directly into the image
	// data.
	var img image.Image
	var data []uint8
	if config.ColorModel.(netpbmHeader).Maxval < 256 {
		gray := image.NewGray(image.Rect(0, 0, config.Width, config.Height))
		data = gray.Pix
		img = gray
	} else {
		gray16 := image.NewGray16(image.Rect(0, 0, config.Width, config.Height))
		data = gray16.Pix
		img = gray16
	}
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

func init() {
	image.RegisterFormat("pgm", "P5", decodePGM, decodeConfigPGM)
}
