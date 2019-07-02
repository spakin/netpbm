// Test PPM files.

package netpbm

import (
	"bytes"
	"compress/flate"
	"testing"
)

// TestNetpbmDecodePGMPPMOpts determines if netpbm.Decode can decode a PGM file
// with PPM options.
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

// TestDecodePBMEncodePPM confirms that a PBM file can be re-encoded as PPM.
func TestDecodePBMEncodePPM(t *testing.T) {
	var w bytes.Buffer
	img := imageFromString(t, pbmRaw, PBM)
	opts := &EncodeOptions{Format: PPM}
	err := Encode(&w, img, opts)
	if err != nil {
		t.Fatal(err)
	}
}

// TestDecodePGMEncodePPM confirms that a PGM file can be re-encoded as PPM.
func TestDecodePGMEncodePPM(t *testing.T) {
	var w bytes.Buffer
	img := imageFromString(t, pgmRaw, PGM)
	opts := &EncodeOptions{Format: PPM}
	err := Encode(&w, img, opts)
	if err != nil {
		t.Fatal(err)
	}
}

// TestNetpbmDecodePlainPBMPPMOpts determines if netpbm.Decode can decode a
// plain PBM file with PPM options.
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


// TestAddRemoveAlphaPPM determins if we can add an alpha channel to a PPM
// file, remove it, and wind up with the same PPM image as we started with.
func TestAddRemoveAlphaPPM(t *testing.T) {
	addRemoveAlpha(t, ppmRaw, nil, nil)
}
