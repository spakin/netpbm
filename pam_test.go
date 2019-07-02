// Test PAM files.

package netpbm

import (
	"bytes"
	"compress/flate"
	"testing"
)

// TestNetpbmDecodePGMPAMOpts determines if netpbm.Decode can decode a PGM file
// with PAM options.
func TestNetpbmDecodePGMPAMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pgmRaw))
	defer r.Close()
	img, err := Decode(r, &DecodeOptions{
		Target: PAM,
		Exact:  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PGM {
		t.Fatalf("Expected PGM but received %s", img.Format())
	}
}

// TestDecodePBMEncodePAM confirms that a PBM file can be re-encoded as PAM.
func TestDecodePBMEncodePAM(t *testing.T) {
	var w bytes.Buffer
	img := imageFromString(t, pbmRaw, PBM)
	opts := &EncodeOptions{Format: PAM}
	err := Encode(&w, img, opts)
	if err != nil {
		t.Fatal(err)
	}
}

// TestDecodePGMEncodePAM confirms that a PGM file can be re-encoded as PAM.
func TestDecodePGMEncodePAM(t *testing.T) {
	var w bytes.Buffer
	img := imageFromString(t, pgmRaw, PGM)
	opts := &EncodeOptions{Format: PAM}
	err := Encode(&w, img, opts)
	if err != nil {
		t.Fatal(err)
	}
}

// TestNetpbmDecodePlainPBMPAMOpts determines if netpbm.Decode can decode a
// plain PBM file with PAM options.
func TestNetpbmDecodePlainPBMPAMOpts(t *testing.T) {
	r := flate.NewReader(bytes.NewBufferString(pbmPlain))
	defer r.Close()
	img, err := Decode(r, &DecodeOptions{
		Target: PAM,
		Exact:  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Format() != PBM {
		t.Fatalf("Expected PBM but received %s", img.Format())
	}
}
