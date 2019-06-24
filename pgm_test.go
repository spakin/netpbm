// Test PGM files

package netpbm

import (
	"bytes"
	"compress/flate"
	"testing"
)

// TestNetpbmDecodeRawPGMConfig determines if netpbm.DecodeConfig can decode
// the configuration of a raw PGM file.
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

// TestNetpbmDecodeRawPGM determines if netpbm.Decode can decode a raw PGM
// file.
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

// TestNetpbmDecodeRawPGMOpts determines if netpbm.Decode can decode a raw PGM
// file with non-default options.
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

// TestNetpbmDecodePBMPGMOpts determines if netpbm.Decode can decode a PBM file
// with PGM options.
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

// TestNetpgmEncodePGM confirms that encoding and decoding do not alter a raw
// PGM file.
func TestNetpgmEncodePGM(t *testing.T) {
	repeatDecodeEncode(t, pgmRaw, nil, nil)
}

// TestNetppmEncodePGMAsPNM confirms that encoding and decoding do not alter a
// raw PGM file when treated as PNM.
func TestNetppmEncodePGMAsPNM(t *testing.T) {
	dOpts := &DecodeOptions{Target: PNM}
	eOpts := &EncodeOptions{Format: PNM}
	repeatDecodeEncode(t, pgmRaw, dOpts, eOpts)
}

// TestDecodePBMEncodePGM confirms that a PBM file can be re-encoded as PGM.
func TestDecodePBMEncodePGM(t *testing.T) {
	var w bytes.Buffer
	img := imageFromString(t, pbmRaw, PBM)
	opts := &EncodeOptions{Format: PGM}
	err := Encode(&w, img, opts)
	if err != nil {
		t.Fatal(err)
	}
}

// TestDecodePPMEncodePGM confirms that a PPM file can be re-encoded as PGM.
func TestDecodePPMEncodePGM(t *testing.T) {
	var w bytes.Buffer
	img := imageFromString(t, ppmRaw, PPM)
	opts := &EncodeOptions{Format: PGM}
	err := Encode(&w, img, opts)
	if err != nil {
		t.Fatal(err)
	}
}

// TestNetpbmDecodePlainPGMConfig determines if netpbm.DecodeConfig can decode
// the configuration of a plain PGM file.
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

// TestNetpbmDecodePlainPGM determines if netpbm.Decode can decode a plain PGM
// file.
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

// TestNetpbmDecodePlainPGMOpts determines if netpbm.Decode can decode a plain
// PGM file with non-default options.
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

// TestNetpgmEncodePlainPGM confirms that encoding and decoding do not alter a
// plain PGM file.
func TestNetpgmEncodePlainPGM(t *testing.T) {
	repeatDecodeEncode(t, pgmPlain, nil, nil)
}

// TestNetppmEncodePlainPGMAsPNM confirms that encoding and decoding do not
// alter a plain PGM file when treated as PNM.
func TestNetppmEncodePlainPGMAsPNM(t *testing.T) {
	dOpts := &DecodeOptions{Target: PNM}
	eOpts := &EncodeOptions{Format: PNM}
	repeatDecodeEncode(t, pgmPlain, dOpts, eOpts)
}
