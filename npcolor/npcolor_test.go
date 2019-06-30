// Test color conversions

package npcolor

import (
	"image/color"
	"math/rand"
	"reflect"
	"testing"
)

const numConversions = 100000 // Number of tests of each conversion type

// colorsNearEqual says whether all components of two arbitrary Netpbm colors
// are within a given threshold of each other.  Note that no error checking is
// performed to ensure that the types are comparable.
func colorsNearEqual(a, b interface{}, maxDelta uint64) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	nFields := aVal.NumField()
	for i := 0; i < nFields; i++ {
		aInt := aVal.Field(i).Uint()
		bInt := bVal.Field(i).Uint()
		var delta uint64
		if aInt > bInt {
			delta = aInt - bInt
		} else {
			delta = bInt - aInt
		}
		if delta > maxDelta {
			return false
		}
	}
	return true
}

// TestGrayMToFrom repeatedly chooses a random maximum value and grayscale
// value, creates a GrayM from those, converts that to RGBA, converts the
// result back to grayscale, and ensures that the result matches the original
// value.
func TestGrayMToFrom(t *testing.T) {
	for i := 0; i < numConversions; i++ {
		m := rand.Intn(255) + 1 // [1, 255]
		y := rand.Intn(m + 1)   // [0, m]
		gm1 := GrayM{Y: uint8(y), M: uint8(m)}
		r, g, b, a := gm1.RGBA()
		rgba := color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: uint8(a >> 8),
		}
		gmm := GrayMModel{M: uint8(m)}
		c := gmm.Convert(rgba)
		gm2, ok := c.(GrayM)
		if !ok {
			t.Fatalf("%#v is not a GrayM", c)
		}
		if !colorsNearEqual(gm1, gm2, 1) {
			t.Fatalf("Started with %v but ended with %v on trial %d:",
				gm1, gm2, i+1)
		}
	}
}

// TestGrayM32ToFrom repeatedly chooses a random maximum value and grayscale
// value, creates a GrayM32 from those, converts that to RGBA64, converts the
// result back to grayscale, and ensures that the result matches the original
// value.
func TestGrayM32ToFrom(t *testing.T) {
	for i := 0; i < numConversions; i++ {
		m := rand.Intn(65535) + 1 // [1, 65535]
		y := rand.Intn(m + 1)     // [0, m]
		gm1 := GrayM32{Y: uint16(y), M: uint16(m)}
		r, g, b, a := gm1.RGBA()
		rgba := color.RGBA64{
			R: uint16(r),
			G: uint16(g),
			B: uint16(b),
			A: uint16(a),
		}
		gmm := GrayM32Model{M: uint16(m)}
		c := gmm.Convert(rgba)
		gm2, ok := c.(GrayM32)
		if !ok {
			t.Fatalf("%#v is not a GrayM32", c)
		}
		if !colorsNearEqual(gm1, gm2, 0) {
			t.Fatalf("Started with %v but ended with %v on trial %d:",
				gm1, gm2, i+1)
		}
	}
}

// TestRGBMToFrom repeatedly chooses a random maximum value and RGB value,
// creates an RGBM from those, converts that to RGBA, converts the result back
// to RGB, and ensures that the result matches the original value.
func TestRGBMToFrom(t *testing.T) {
	for i := 0; i < numConversions; i++ {
		m := rand.Intn(255) + 1 // [1, 255]
		rm := rand.Intn(m + 1)  // [0, m]
		gm := rand.Intn(m + 1)  // [0, m]
		bm := rand.Intn(m + 1)  // [0, m]
		rgbm1 := RGBM{
			R: uint8(rm),
			G: uint8(gm),
			B: uint8(bm),
			M: uint8(m),
		}
		r, g, b, a := rgbm1.RGBA()
		rgba := color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: uint8(a >> 8),
		}
		gmm := RGBMModel{M: uint8(m)}
		c := gmm.Convert(rgba)
		rgbm2, ok := c.(RGBM)
		if !ok {
			t.Fatalf("%#v is not a RGBM", c)
		}
		if !colorsNearEqual(rgbm1, rgbm2, 1) {
			t.Fatalf("Started with %v but ended with %v on trial %d:",
				rgbm1, rgbm2, i+1)
		}
	}
}

// TestRGBM64ToFrom repeatedly chooses a random maximum value and RGB value,
// creates an RGBM64 from those, converts that to RGBA64, converts the result
// back to RGB, and ensures that the result matches the original value.
func TestRGBM64ToFrom(t *testing.T) {
	for i := 0; i < numConversions; i++ {
		m := rand.Intn(65535) + 1 // [1, 65535]
		rm := rand.Intn(m + 1)    // [0, m]
		gm := rand.Intn(m + 1)    // [0, m]
		bm := rand.Intn(m + 1)    // [0, m]
		rgbm1 := RGBM64{
			R: uint16(rm),
			G: uint16(gm),
			B: uint16(bm),
			M: uint16(m),
		}
		r, g, b, a := rgbm1.RGBA()
		rgba := color.RGBA64{
			R: uint16(r),
			G: uint16(g),
			B: uint16(b),
			A: uint16(a),
		}
		gmm := RGBM64Model{M: uint16(m)}
		c := gmm.Convert(rgba)
		rgbm2, ok := c.(RGBM64)
		if !ok {
			t.Fatalf("%#v is not a RGBM64", c)
		}
		if !colorsNearEqual(rgbm1, rgbm2, 0) {
			t.Fatalf("Started with %v but ended with %v on trial %d:",
				rgbm1, rgbm2, i+1)
		}
	}
}

// TestRGBM64Premult tests that RGBA returns a properly alpha-premultiplied
// value for an RGBM64 color.
func TestRGBM64Premult(t *testing.T) {
	for i := 0; i < numConversions; i++ {
		m := rand.Intn(65535) + 1 // [1, 65535]
		rm := rand.Intn(m + 1)    // [0, m]
		gm := rand.Intn(m + 1)    // [0, m]
		bm := rand.Intn(m + 1)    // [0, m]
		c := RGBM64{
			R: uint16(rm),
			G: uint16(gm),
			B: uint16(bm),
			M: uint16(m),
		}
		r, g, b, a := c.RGBA()
		if r > a || g > a || b > a {
			t.Fatalf("RGBA color [%d, %d, %d, %d] is not alpha-premultiplied", r, g, b, a)
		}
	}
}

// TestRGBMPremult tests that RGBA returns a properly alpha-premultiplied value
// for an RGBM color.
func TestRGBMPremult(t *testing.T) {
	for i := 0; i < numConversions; i++ {
		m := rand.Intn(255) + 1 // [1, 255]
		rm := rand.Intn(m + 1)  // [0, m]
		gm := rand.Intn(m + 1)  // [0, m]
		bm := rand.Intn(m + 1)  // [0, m]
		c := RGBM{
			R: uint8(rm),
			G: uint8(gm),
			B: uint8(bm),
			M: uint8(m),
		}
		r, g, b, a := c.RGBA()
		if r > a || g > a || b > a {
			t.Fatalf("RGBA color [%d, %d, %d, %d] is not alpha-premultiplied", r, g, b, a)
		}
	}
}

// TestGrayM32Premult tests that RGBA returns a properly alpha-premultiplied
// value for an GrayM32 color.
func TestGrayM32Premult(t *testing.T) {
	for i := 0; i < numConversions; i++ {
		m := rand.Intn(65535) + 1 // [1, 65535]
		ym := rand.Intn(m + 1)    // [0, m]
		c := GrayM32{
			Y: uint16(ym),
			M: uint16(m),
		}
		r, g, b, a := c.RGBA()
		if r > a || g > a || b > a {
			t.Fatalf("RGBA color [%d, %d, %d, %d] is not alpha-premultiplied", r, g, b, a)
		}
	}
}

// TestGrayMPremult tests that RGBA returns a properly alpha-premultiplied value
// for an GrayM color.
func TestGrayMPremult(t *testing.T) {
	for i := 0; i < numConversions; i++ {
		m := rand.Intn(255) + 1 // [1, 255]
		ym := rand.Intn(m + 1)  // [0, m]
		c := GrayM{
			Y: uint8(ym),
			M: uint8(m),
		}
		r, g, b, a := c.RGBA()
		if r > a || g > a || b > a {
			t.Fatalf("RGBA color [%d, %d, %d, %d] is not alpha-premultiplied", r, g, b, a)
		}
	}
}

// TestRGBAM64Premult tests that RGBA returns a properly alpha-premultiplied
// value for an RGBAM64 color.
func TestRGBAM64Premult(t *testing.T) {
	for i := 0; i < numConversions; i++ {
		m := rand.Intn(65535) + 1 // [1, 65535]
		rm := rand.Intn(m + 1)    // [0, m]
		gm := rand.Intn(m + 1)    // [0, m]
		bm := rand.Intn(m + 1)    // [0, m]
		am := rand.Intn(m + 1)    // [0, m]
		c := RGBAM64{
			R: uint16(rm),
			G: uint16(gm),
			B: uint16(bm),
			A: uint16(am),
			M: uint16(m),
		}
		r, g, b, a := c.RGBA()
		if r > a || g > a || b > a {
			t.Fatalf("RGBA color [%d, %d, %d, %d] is not alpha-premultiplied", r, g, b, a)
		}
	}
}

// TestRGBAMPremult tests that RGBA returns a properly alpha-premultiplied value
// for an RGBAM color.
func TestRGBAMPremult(t *testing.T) {
	for i := 0; i < numConversions; i++ {
		m := rand.Intn(255) + 1 // [1, 255]
		rm := rand.Intn(m + 1)  // [0, m]
		gm := rand.Intn(m + 1)  // [0, m]
		bm := rand.Intn(m + 1)  // [0, m]
		am := rand.Intn(m + 1)  // [0, m]
		c := RGBAM{
			R: uint8(rm),
			G: uint8(gm),
			B: uint8(bm),
			A: uint8(am),
			M: uint8(m),
		}
		r, g, b, a := c.RGBA()
		if r > a || g > a || b > a {
			t.Fatalf("RGBA color [%d, %d, %d, %d] is not alpha-premultiplied", r, g, b, a)
		}
	}
}
