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

// TestNetpbmDecodeRawPBMConfig determines if netpbm.DecodeConfig can
// decode the configuration of a raw PBM file.
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

// TestNetpbmDecodeRawPBM determines if netpbm.Decode can decode a
// raw PBM file.
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

// TestNetpbmDecodeRawPBMOpts determines if netpbm.Decode can decode a
// raw PBM file with non-default options.
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

// TestNetpbmEncodeRawPBM confirms that encoding and decoding do not alter
// a raw PBM file.
func TestNetpbmEncodeRawPBM(t *testing.T) {
	repeatDecodeEncode(t, pbmRaw, nil, nil)
}

// TestNetpbmEncodeRawPBMAsPNM confirms that encoding and decoding do not
// alter a raw PBM file when treated as PNM.
func TestNetpbmEncodeRawPBMAsPNM(t *testing.T) {
	dOpts := &DecodeOptions{Target: PNM}
	eOpts := &EncodeOptions{Format: PNM}
	repeatDecodeEncode(t, pbmRaw, dOpts, eOpts)
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

// TestNetpbmDecodePlainPBMConfig determines if netpbm.DecodeConfig can
// decode the configuration of a plain PBM file.
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

// TestNetpbmDecodePlainPBM determines if netpbm.Decode can decode a
// plain PBM file.
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

// TestNetpbmDecodePlainPBMOpts determines if netpbm.Decode can decode a
// plain PBM file with non-default options.
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
}

// TestNetpbmEncodePlainPBM confirms that encoding and decoding do not alter
// a plain PBM file.
func TestNetpbmEncodePlainPBM(t *testing.T) {
	repeatDecodeEncode(t, pbmPlain, nil, nil)
}

// TestNetpbmEncodePlainPBMAsPNM confirms that encoding and decoding do not
// alter a plain PBM file when treated as PNM.
func TestNetpbmEncodePlainPBMAsPNM(t *testing.T) {
	dOpts := &DecodeOptions{Target: PNM}
	eOpts := &EncodeOptions{Format: PNM}
	repeatDecodeEncode(t, pbmPlain, dOpts, eOpts)
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

// TestNetpbmDecodeRawPGMConfig determines if netpbm.DecodeConfig can
// decode the configuration of a raw PGM file.
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

// TestNetpbmDecodeRawPGM determines if netpbm.Decode can decode a
// raw PGM file.
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

// TestNetpbmDecodeRawPGMOpts determines if netpbm.Decode can decode a
// raw PGM file with non-default options.
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

// TestNetpbmEncodeRawPGM confirms that encoding and decoding do not alter
// a raw PGM file.
func TestNetpbmEncodeRawPGM(t *testing.T) {
	repeatDecodeEncode(t, pgmRaw, nil, nil)
}

// TestNetpbmEncodeRawPGMAsPNM confirms that encoding and decoding do not
// alter a raw PGM file when treated as PNM.
func TestNetpbmEncodeRawPGMAsPNM(t *testing.T) {
	dOpts := &DecodeOptions{Target: PNM}
	eOpts := &EncodeOptions{Format: PNM}
	repeatDecodeEncode(t, pgmRaw, dOpts, eOpts)
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

// TestNetpbmDecodePlainPGMConfig determines if netpbm.DecodeConfig can
// decode the configuration of a plain PGM file.
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

// TestNetpbmDecodePlainPGM determines if netpbm.Decode can decode a
// plain PGM file.
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

// TestNetpbmDecodePlainPGMOpts determines if netpbm.Decode can decode a
// plain PGM file with non-default options.
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
}

// TestNetpbmEncodePlainPGM confirms that encoding and decoding do not alter
// a plain PGM file.
func TestNetpbmEncodePlainPGM(t *testing.T) {
	repeatDecodeEncode(t, pgmPlain, nil, nil)
}

// TestNetpbmEncodePlainPGMAsPNM confirms that encoding and decoding do not
// alter a plain PGM file when treated as PNM.
func TestNetpbmEncodePlainPGMAsPNM(t *testing.T) {
	dOpts := &DecodeOptions{Target: PNM}
	eOpts := &EncodeOptions{Format: PNM}
	repeatDecodeEncode(t, pgmPlain, dOpts, eOpts)
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

// TestNetpbmDecodeRawPPMConfig determines if netpbm.DecodeConfig can
// decode the configuration of a raw PPM file.
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

// TestNetpbmDecodeRawPPM determines if netpbm.Decode can decode a
// raw PPM file.
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

// TestNetpbmDecodeRawPPMOpts determines if netpbm.Decode can decode a
// raw PPM file with non-default options.
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

// TestNetpbmEncodeRawPPM confirms that encoding and decoding do not alter
// a raw PPM file.
func TestNetpbmEncodeRawPPM(t *testing.T) {
	repeatDecodeEncode(t, ppmRaw, nil, nil)
}

// TestNetpbmEncodeRawPPMAsPNM confirms that encoding and decoding do not
// alter a raw PPM file when treated as PNM.
func TestNetpbmEncodeRawPPMAsPNM(t *testing.T) {
	dOpts := &DecodeOptions{Target: PNM}
	eOpts := &EncodeOptions{Format: PNM}
	repeatDecodeEncode(t, ppmRaw, dOpts, eOpts)
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

// TestNetpbmDecodePlainPPMConfig determines if netpbm.DecodeConfig can
// decode the configuration of a plain PPM file.
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

// TestNetpbmDecodePlainPPM determines if netpbm.Decode can decode a
// plain PPM file.
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

// TestNetpbmDecodePlainPPMOpts determines if netpbm.Decode can decode a
// plain PPM file with non-default options.
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
}

// TestNetpbmEncodePlainPPM confirms that encoding and decoding do not alter
// a plain PPM file.
func TestNetpbmEncodePlainPPM(t *testing.T) {
	repeatDecodeEncode(t, ppmPlain, nil, nil)
}

// TestNetpbmEncodePlainPPMAsPNM confirms that encoding and decoding do not
// alter a plain PPM file when treated as PNM.
func TestNetpbmEncodePlainPPMAsPNM(t *testing.T) {
	dOpts := &DecodeOptions{Target: PNM}
	eOpts := &EncodeOptions{Format: PNM}
	repeatDecodeEncode(t, ppmPlain, dOpts, eOpts)
}

// TestDecodeRawPAMConfig determines if image.DecodeConfig can decode
// the configuration of a raw PAM file.
func TestDecodeRawPAMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawColor))
	defer r.Close()
	cfg, str, err := image.DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pam" {
		t.Fatalf("Expected \"pam\" but received %q", str)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestDecodeRawPAM determines if image.Decode can decode a raw PAM file.
func TestDecodeRawPAM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawColor))
	defer r.Close()
	img, str, err := image.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pam" {
		t.Fatalf("Expected pam but received %s", str)
	}
	nimg, ok := img.(Image)
	if !ok {
		t.Fatal("Image is not a Netpam image")
	}
	if nimg.MaxValue() != 255 {
		t.Fatalf("Expected a maximum value of 255 but received %d", nimg.MaxValue())
	}
}

// TestNetpbmDecodeRawPAMConfig determines if netpbm.DecodeConfig can
// decode the configuration of a raw PAM file.
func TestNetpbmDecodeRawPAMConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawColor))
	defer r.Close()
	cfg, err := DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestNetpbmDecodeRawPAM determines if netpbm.Decode can decode a
// raw PAM file.
func TestNetpbmDecodeRawPAM(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawColor))
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

// TestNetpbmDecodeRawPAMOpts determines if netpbm.Decode can decode a
// raw PAM file with non-default options.
func TestNetpbmDecodeRawPAMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawColor))
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

// TestNetpbmEncodeRawPAM confirms that encoding and decoding do not alter
// a raw PAM file.
func TestNetpbmEncodeRawPAM(t *testing.T) {
	repeatDecodeEncode(t, pamRawColor, nil, nil)
}

// TestNetpbmEncodeRawPAMAsPNM confirms that encoding and decoding do not
// alter a raw PAM file when treated as PNM.
func TestNetpbmEncodeRawPAMAsPNM(t *testing.T) {
	dOpts := &DecodeOptions{Target: PNM}
	eOpts := &EncodeOptions{Format: PNM}
	repeatDecodeEncode(t, pamRawColor, dOpts, eOpts)
}

// TestDecodeRawPAMGrayAlphaConfig determines if image.DecodeConfig can decode
// the configuration of a raw PAM file.
func TestDecodeRawPAMGrayAlphaConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawGrayAlpha))
	defer r.Close()
	cfg, str, err := image.DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pam" {
		t.Fatalf("Expected \"pam\" but received %q", str)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestDecodeRawPAMGrayAlpha determines if image.Decode can decode a raw PAM file.
func TestDecodeRawPAMGrayAlpha(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawGrayAlpha))
	defer r.Close()
	img, str, err := image.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pam" {
		t.Fatalf("Expected pam but received %s", str)
	}
	nimg, ok := img.(Image)
	if !ok {
		t.Fatal("Image is not a Netpam image")
	}
	if nimg.MaxValue() != 255 {
		t.Fatalf("Expected a maximum value of 255 but received %d", nimg.MaxValue())
	}
}

// TestNetpbmDecodeRawPAMGrayAlphaConfig determines if netpbm.DecodeConfig can
// decode the configuration of a raw PAM file.
func TestNetpbmDecodeRawPAMGrayAlphaConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawGrayAlpha))
	defer r.Close()
	cfg, err := DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestNetpbmDecodeRawPAMGrayAlpha determines if netpbm.Decode can decode a
// raw PAM file.
func TestNetpbmDecodeRawPAMGrayAlpha(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawGrayAlpha))
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

// TestNetpbmDecodeRawPAMGrayAlphaOpts determines if netpbm.Decode can decode a
// raw PAM file with non-default options.
func TestNetpbmDecodeRawPAMGrayAlphaOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawGrayAlpha))
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

// TestNetpbmEncodeRawPAMGrayAlpha confirms that encoding and decoding do not alter
// a raw PAM file.
func TestNetpbmEncodeRawPAMGrayAlpha(t *testing.T) {
	repeatDecodeEncode(t, pamRawGrayAlpha, nil, nil)
}

// TestNetpbmEncodeRawPAMGrayAlphaAsPNM confirms that encoding and decoding do not
// alter a raw PAM file when treated as PNM.
func TestNetpbmEncodeRawPAMGrayAlphaAsPNM(t *testing.T) {
	dOpts := &DecodeOptions{Target: PNM}
	eOpts := &EncodeOptions{Format: PNM}
	repeatDecodeEncode(t, pamRawGrayAlpha, dOpts, eOpts)
}

// TestDecodeRawPAMAlphaConfig determines if image.DecodeConfig can decode
// the configuration of a raw PAM file.
func TestDecodeRawPAMAlphaConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawColorAlpha))
	defer r.Close()
	cfg, str, err := image.DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pam" {
		t.Fatalf("Expected \"pam\" but received %q", str)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestDecodeRawPAMAlpha determines if image.Decode can decode a raw PAM file.
func TestDecodeRawPAMAlpha(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawColorAlpha))
	defer r.Close()
	img, str, err := image.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	if str != "pam" {
		t.Fatalf("Expected pam but received %s", str)
	}
	nimg, ok := img.(Image)
	if !ok {
		t.Fatal("Image is not a Netpam image")
	}
	if nimg.MaxValue() != 255 {
		t.Fatalf("Expected a maximum value of 255 but received %d", nimg.MaxValue())
	}
}

// TestNetpbmDecodeRawPAMAlphaConfig determines if netpbm.DecodeConfig can
// decode the configuration of a raw PAM file.
func TestNetpbmDecodeRawPAMAlphaConfig(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawColorAlpha))
	defer r.Close()
	cfg, err := DecodeConfig(r)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != 64 || cfg.Height != 64 {
		t.Fatalf("Expected a 64x64 image but received %dx%d", cfg.Width, cfg.Height)
	}
}

// TestNetpbmDecodeRawPAMAlpha determines if netpbm.Decode can decode a
// raw PAM file.
func TestNetpbmDecodeRawPAMAlpha(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawColorAlpha))
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

// TestNetpbmDecodeRawPAMAlphaOpts determines if netpbm.Decode can decode a
// raw PAM file with non-default options.
func TestNetpbmDecodeRawPAMAlphaOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pamRawColorAlpha))
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

// TestNetpbmEncodeRawPAMAlpha confirms that encoding and decoding do not alter
// a raw PAM file.
func TestNetpbmEncodeRawPAMAlpha(t *testing.T) {
	repeatDecodeEncode(t, pamRawColorAlpha, nil, nil)
}

// TestNetpbmEncodeRawPAMAlphaAsPNM confirms that encoding and decoding do not
// alter a raw PAM file when treated as PNM.
func TestNetpbmEncodeRawPAMAlphaAsPNM(t *testing.T) {
	dOpts := &DecodeOptions{Target: PNM}
	eOpts := &EncodeOptions{Format: PNM}
	repeatDecodeEncode(t, pamRawColorAlpha, dOpts, eOpts)
}
