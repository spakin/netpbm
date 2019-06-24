// Test PBM files

package netpbm

import (
	"bytes"
	"compress/flate"
	"testing"
)

// TestDecodePGMEncodePBM confirms that a PGM file can be re-encoded as PBM.
func TestDecodePGMEncodePBM(t *testing.T) {
	var w bytes.Buffer
	img := imageFromString(t, pgmRaw, PGM)
	opts := &EncodeOptions{Format: PBM}
	err := Encode(&w, img, opts)
	if err != nil {
		t.Fatal(err)
	}
}

// TestDecodePPMEncodePBM confirms that a PPM file can be re-encoded as PBM.
func TestDecodePPMEncodePBM(t *testing.T) {
	var w bytes.Buffer
	img := imageFromString(t, ppmRaw, PPM)
	opts := &EncodeOptions{Format: PBM}
	err := Encode(&w, img, opts)
	if err != nil {
		t.Fatal(err)
	}
}

// TestDecodePBMComments confirms that we can decode a PBM file
// containing comments.
func TestDecodePBMComments(t *testing.T) {
	// Read the image.
	r := flate.NewReader(bytes.NewBufferString(pbmRawComments))
	defer r.Close()
	_, cs, err := DecodeWithComments(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Confirm that the comments are as expected.
	exp := []string{"This file contains", "a variety of comments...", "", "   ...in all sorts", "of tricky forms."}
	if len(cs) != len(exp) {
		t.Fatalf("Expected %#v but received %#v", exp, cs)
	}
	for i, e := range exp {
		if e != cs[i] {
			t.Fatalf("Expected %q but received %q", e, cs[i])
		}
	}
}
