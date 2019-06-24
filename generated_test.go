/*
This file declares test functions that work on all "real" Netpbm types (PBM,
PGM, PPM, and PAM, but not PNM).

This is a generated file.  DO NOT EDIT.  Edit helpers/all.tmpl instead.
*/

package netpbm

import (
	"bytes"
	"compress/flate"
	"image"
	"testing"
)

// TestDecodeRawPBMConfig determines if image.DecodeConfig can decode
// the configuration of a raw PBM file.
func TestDecodeRawPBMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmRaw))
	defer r.Close()
	cfg, str, err := image.DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pbm" {
		t.Fatalf("Expected \"pbm\" but received %q", str)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestDecodeRawPBM determines if image.Decode can decode a raw PBM file.
func TestDecodeRawPBM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmRaw))
	defer r.Close()
	img, str, err := image.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pbm" {
		t.Fatalf("Expected pbm but received %s", str)
	}
	nimg, ok := img.(Image)
	if !ok {
		t.Fatal("Image is not a Netpbm image")
	}
	if nimg.MaxValue() != 1 {
		t.Fatalf("Expected a maximum value of 1 but received %d", nimg.MaxValue())
	}
}

// TestDecodePlainPBMConfig determines if image.DecodeConfig can decode
// the configuration of a plain PBM file.
func TestDecodePlainPBMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmPlain))
	defer r.Close()
	cfg, str, err := image.DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pbm" {
		t.Fatalf("Expected \"pbm\" but received %q", str)
	}
	if cfg.Width != 63 || cfg.Height != 65 {
		t.Fatalf("Expected a 63x65 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestDecodePlainPBM determines if image.Decode can decode a plain PBM file.
func TestDecodePlainPBM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmPlain))
	defer r.Close()
	img, str, err := image.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pbm" {
		t.Fatalf("Expected pbm but received %s", str)
	}
	nimg, ok := img.(Image)
	if !ok {
		t.Fatal("Image is not a Netpbm image")
	}
	if nimg.MaxValue() != 1 {
		t.Fatalf("Expected a maximum value of 1 but received %d", nimg.MaxValue())
	}
}

// TestDecodeRawPGMConfig determines if image.DecodeConfig can decode
// the configuration of a raw PGM file.
func TestDecodeRawPGMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmRaw))
	defer r.Close()
	cfg, str, err := image.DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pgm" {
		t.Fatalf("Expected \"pgm\" but received %q", str)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestDecodeRawPGM determines if image.Decode can decode a raw PGM file.
func TestDecodeRawPGM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmRaw))
	defer r.Close()
	img, str, err := image.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pgm" {
		t.Fatalf("Expected pgm but received %s", str)
	}
	nimg, ok := img.(Image)
	if !ok {
		t.Fatal("Image is not a Netpgm image")
	}
	if nimg.MaxValue() != 255 {
		t.Fatalf("Expected a maximum value of 255 but received %d", nimg.MaxValue())
	}
}

// TestDecodePlainPGMConfig determines if image.DecodeConfig can decode
// the configuration of a plain PGM file.
func TestDecodePlainPGMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmPlain))
	defer r.Close()
	cfg, str, err := image.DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pgm" {
		t.Fatalf("Expected \"pgm\" but received %q", str)
	}
	if cfg.Width != 63 || cfg.Height != 65 {
		t.Fatalf("Expected a 63x65 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestDecodePlainPGM determines if image.Decode can decode a plain PGM file.
func TestDecodePlainPGM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmPlain))
	defer r.Close()
	img, str, err := image.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pgm" {
		t.Fatalf("Expected pgm but received %s", str)
	}
	nimg, ok := img.(Image)
	if !ok {
		t.Fatal("Image is not a Netpgm image")
	}
	if nimg.MaxValue() != 777 {
		t.Fatalf("Expected a maximum value of 777 but received %d", nimg.MaxValue())
	}
}

// TestDecodeRawPPMConfig determines if image.DecodeConfig can decode
// the configuration of a raw PPM file.
func TestDecodeRawPPMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(ppmRaw))
	defer r.Close()
	cfg, str, err := image.DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "ppm" {
		t.Fatalf("Expected \"ppm\" but received %q", str)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestDecodeRawPPM determines if image.Decode can decode a raw PPM file.
func TestDecodeRawPPM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(ppmRaw))
	defer r.Close()
	img, str, err := image.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "ppm" {
		t.Fatalf("Expected ppm but received %s", str)
	}
	nimg, ok := img.(Image)
	if !ok {
		t.Fatal("Image is not a Netppm image")
	}
	if nimg.MaxValue() != 255 {
		t.Fatalf("Expected a maximum value of 255 but received %d", nimg.MaxValue())
	}
}

// TestDecodePlainPPMConfig determines if image.DecodeConfig can decode
// the configuration of a plain PPM file.
func TestDecodePlainPPMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(ppmPlain))
	defer r.Close()
	cfg, str, err := image.DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "ppm" {
		t.Fatalf("Expected \"ppm\" but received %q", str)
	}
	if cfg.Width != 63 || cfg.Height != 65 {
		t.Fatalf("Expected a 63x65 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestDecodePlainPPM determines if image.Decode can decode a plain PPM file.
func TestDecodePlainPPM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(ppmPlain))
	defer r.Close()
	img, str, err := image.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "ppm" {
		t.Fatalf("Expected ppm but received %s", str)
	}
	nimg, ok := img.(Image)
	if !ok {
		t.Fatal("Image is not a Netppm image")
	}
	if nimg.MaxValue() != 777 {
		t.Fatalf("Expected a maximum value of 777 but received %d", nimg.MaxValue())
	}
}
