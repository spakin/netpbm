// Test PBM files

package netpbm

import (
	"bytes"
	"compress/flate"
	"image"
	"testing"
)

// Determine if image.DecodeConfig can decode the configuration of a raw PBM
// file.
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

// Determine if image.Decode can decode a raw PBM file.
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

// Determine if netpbm.DecodeConfig can decode the configuration of a raw PBM
// file.
func TestNetpbmDecodeRawPBMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmRaw))
	defer r.Close()
	cfg, err := DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// Determine if netpbm.Decode can decode a raw PBM file.
func TestNetpbmDecodeRawPBM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmRaw))
	defer r.Close()
	img, err := Decode(r, nil)
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PBM {
		t.Fatalf("Expected PBM but received %s", img.Format())
	}
	if img.MaxValue() != 1 {
		t.Fatalf("Expected a maximum value of 1 but received %d", img.MaxValue())
	}
}

// Determine if netpbm.Decode can decode a raw PBM file with non-default
// options.
func TestNetpbmDecodeRawPBMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmRaw))
	defer r.Close()
	img, err := Decode(r, &DecodeOptions{
		Target: PBM,
		Exact:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PBM {
		t.Fatalf("Expected PBM but received %s", img.Format())
	}
}

// Determine if image.DecodeConfig can decode the configuration of a plain PBM
// file.
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

// Determine if image.Decode can decode a plain PBM file.
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

// Determine if netpbm.DecodeConfig can decode the configuration of a plain PBM
// file.
func TestNetpbmDecodePlainPBMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmPlain))
	defer r.Close()
	cfg, err := DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != 63 || cfg.Height != 65 {
		t.Fatalf("Expected a 63x65 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// Determine if netpbm.Decode can decode a plain PBM file.
func TestNetpbmDecodePlainPBM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmPlain))
	defer r.Close()
	img, err := Decode(r, nil)
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PBM {
		t.Fatalf("Expected PBM but received %s", img.Format())
	}
	if img.MaxValue() != 1 {
		t.Fatalf("Expected a maximum value of 1 but received %d", img.MaxValue())
	}
}

// Determine if netpbm.Decode can decode a plain PBM file with non-default
// options.
func TestNetpbmDecodePlainPBMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmPlain))
	defer r.Close()
	img, err := Decode(r, &DecodeOptions{
		Target: PBM,
		Exact:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PBM {
		t.Fatalf("Expected PBM but received %s", img.Format())
	}
	if img.MaxValue() != 1 {
		t.Fatalf("Expected a maximum value of 1 but received %d", img.MaxValue())
	}
}
