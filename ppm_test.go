// Test PPM files

package netpbm

import (
	"bytes"
	"compress/flate"
	"image"
	"testing"
)

// Determine if image.DecodeConfig can decode the configuration of a raw PPM
// file.
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

// Determine if image.Decode can decode a raw PPM file.
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
		t.Fatal("Image is not a Netpbm image")
	}
	if nimg.MaxValue() != 255 {
		t.Fatalf("Expected a maximum value of 266 but received %d", nimg.MaxValue())
	}
}

// Determine if netpbm.DecodeConfig can decode the configuration of a raw PPM
// file.
func TestNetpbmDecodeRawPPMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(ppmRaw))
	defer r.Close()
	cfg, err := DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// Determine if netpbm.Decode can decode a raw PPM file.
func TestNetpbmDecodeRawPPM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(ppmRaw))
	defer r.Close()
	img, err := Decode(r, nil)
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PPM {
		t.Fatalf("Expected PPM but received %s", img.Format())
	}
	if img.MaxValue() != 255 {
		t.Fatalf("Expected a maximum value of 255 but received %d", img.MaxValue())
	}
}

// Determine if netpbm.Decode can decode a raw PPM file with non-default
// options.
func TestNetpbmDecodeRawPPMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(ppmRaw))
	defer r.Close()
	img, err := Decode(r, &DecodeOptions{
		Target: PPM,
		Exact:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PPM {
		t.Fatalf("Expected PPM but received %s", img.Format())
	}
}

// Determine if netpbm.Decode can decode a PGM file with PPM options.
func TestNetpbmDecodePGMPPMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmRaw))
	defer r.Close()
	img, err := Decode(r, &DecodeOptions{
		Target: PPM,
		Exact:  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PPM {
		t.Fatalf("Expected PPM but received %s", img.Format())
	}
}

// Determine if netpbm.Decode can decode a plain PBM file with PPM options.
func TestNetpbmDecodePlainPBMPPMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmPlain))
	defer r.Close()
	img, err := Decode(r, &DecodeOptions{
		Target: PPM,
		Exact:  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PPM {
		t.Fatalf("Expected PPM but received %s", img.Format())
	}
}

// Determine if image.DecodeConfig can decode the configuration of a plain PPM
// file.
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

// Determine if image.Decode can decode a plain PPM file.
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
		t.Fatal("Image is not a Netpbm image")
	}
	if nimg.MaxValue() != 777 {
		t.Fatalf("Expected a maximum value of 777 but received %d", nimg.MaxValue())
	}
}

// Determine if netpbm.DecodeConfig can decode the configuration of a plain PPM
// file.
func TestNetpbmDecodePlainPPMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(ppmPlain))
	defer r.Close()
	cfg, err := DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != 63 || cfg.Height != 65 {
		t.Fatalf("Expected a 63x65 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// Determine if netpbm.Decode can decode a plain PPM file.
func TestNetpbmDecodePlainPPM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(ppmPlain))
	defer r.Close()
	img, err := Decode(r, nil)
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PPM {
		t.Fatalf("Expected PPM but received %s", img.Format())
	}
	if img.MaxValue() != 777 {
		t.Fatalf("Expected a maximum value of 777 but received %d", img.MaxValue())
	}
}

// Determine if netpbm.Decode can decode a plain PPM file with non-default
// options.
func TestNetpbmDecodePlainPPMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(ppmPlain))
	defer r.Close()
	img, err := Decode(r, &DecodeOptions{
		Target: PPM,
		Exact:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PPM {
		t.Fatalf("Expected PPM but received %s", img.Format())
	}
	if img.MaxValue() != 777 {
		t.Fatalf("Expected a maximum value of 777 but received %d", img.MaxValue())
	}
}
