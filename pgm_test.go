// Test PGM files.

package netpbm

import (
	"bytes"
	"compress/flate"
	"testing"
)

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
