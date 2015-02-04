// Test PGM files

package netpbm

import (
	"bytes"
	"compress/flate"
	"image"
	"testing"
)

// Determine if image.DecodeConfig can decode the configuration of a raw PGM
// file.
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

// Determine if image.Decode can decode a raw PGM file.
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
		t.Fatal("Image is not a Netpbm image")
	}
	if nimg.MaxValue() != 255 {
		t.Fatalf("Expected a maximum value of 266 but received %d", nimg.MaxValue())
	}
}

// Determine if netpbm.DecodeConfig can decode the configuration of a raw PGM
// file.
func TestNetpbmDecodeRawPGMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmRaw))
	defer r.Close()
	cfg, err := DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// Determine if netpbm.Decode can decode a raw PGM file.
func TestNetpbmDecodeRawPGM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmRaw))
	defer r.Close()
	img, err := Decode(r, nil)
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PGM {
		t.Fatalf("Expected PGM but received %s", img.Format())
	}
	if img.MaxValue() != 255 {
		t.Fatalf("Expected a maximum value of 255 but received %d", img.MaxValue())
	}
}

// Determine if netpbm.Decode can decode a raw PGM file with non-default
// options.
func TestNetpbmDecodeRawPGMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmRaw))
	defer r.Close()
	img, err := Decode(r, &DecodeOptions{
		Target: PGM,
		Exact:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PGM {
		t.Fatalf("Expected PGM but received %s", img.Format())
	}
}

// Determine if netpbm.Decode can decode a PBM file with PGM options.
func TestNetpbmDecodePBMPGMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmRaw))
	defer r.Close()
	img, err := Decode(r, &DecodeOptions{
		Target: PGM,
		Exact:  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PGM {
		t.Fatalf("Expected PGM but received %s", img.Format())
	}
}

// Determine if image.DecodeConfig can decode the configuration of a plain PGM
// file.
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

// Determine if image.Decode can decode a plain PGM file.
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
		t.Fatal("Image is not a Netpbm image")
	}
	if nimg.MaxValue() != 777 {
		t.Fatalf("Expected a maximum value of 777 but received %d", nimg.MaxValue())
	}
}

// Determine if netpbm.DecodeConfig can decode the configuration of a plain PGM
// file.
func TestNetpbmDecodePlainPGMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmPlain))
	defer r.Close()
	cfg, err := DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != 63 || cfg.Height != 65 {
		t.Fatalf("Expected a 63x65 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// Determine if netpbm.Decode can decode a plain PGM file.
func TestNetpbmDecodePlainPGM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmPlain))
	defer r.Close()
	img, err := Decode(r, nil)
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PGM {
		t.Fatalf("Expected PGM but received %s", img.Format())
	}
	if img.MaxValue() != 777 {
		t.Fatalf("Expected a maximum value of 777 but received %d", img.MaxValue())
	}
}

// Determine if netpbm.Decode can decode a plain PGM file with non-default
// options.
func TestNetpbmDecodePlainPGMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmPlain))
	defer r.Close()
	img, err := Decode(r, &DecodeOptions{
		Target: PGM,
		Exact:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PGM {
		t.Fatalf("Expected PGM but received %s", img.Format())
	}
	if img.MaxValue() != 777 {
		t.Fatalf("Expected a maximum value of 777 but received %d", img.MaxValue())
	}
}
