// Present a collection of examples to demonstrate netpbm package usage.

package netpbm_test

import (
	"bytes"
	"compress/flate"
	"fmt"
	"image"
	"strings"

	"github.com/spakin/netpbm"
	"github.com/spakin/netpbm/npcolor"
)

func ExampleDecode() {
	// In this example, we read a PBM (black and white) "file" (really a
	// hard-wired string to make the example self-contained) into a PPM
	// (color) image.  Because netpbm.Decode returns a netpbm.Image, not a
	// generic image.Image we can query its Netpbm format and maximum
	// value.  We also output an arbitrary white pixel to show how B&W can
	// be promoted to color.

	const pbmRaw = "\n0\xe1RVp\x0eru\f\xf1\x0f\xb2Rp\xf7\xf4\rP\b\xf0\xf3Up\xcb\xcc)I-R\bK-*\xce\xcc\xcfS0\xd43\xe423Q02\xe1\xfa\x0f\x02\xff\xfe\xff\x06\xd3\x1f\xec\x1b\xc1\xf4a{\x060\xddo\xc3\x01\xa6\xe5l*\xc0\xca\xea,~\x80\xe8?\u007f\n\u007f\x82\xe8\x9f\xff\x1eC\xe8\xff\xcf?\x02\xa9\v\xc6\xff\xfa?\xff\xff\xdf\xc0\xfd|A\xff\xe3\xff\x17\xe2\xff\x1fc\x90?\xfe\x83\xf5\xff\x97\xbeB\xfb\xe3M\xff\xff2\xcc\u007fd/\xbf\xff/\x03\xb7\xfc\xa1:\xf9\xff\f\xfc\xff\xed\xfbj\xec\xff\x89\xff\a\xd2\x15\xf5 \xe3\xed\xe4\x12\xc0\xd6\xd8\xd81\x82i\x81\u007f\xec`\x9a\xf1\u007f?\xd8\x19`\x1e\x1a\x00\x04\x00\x00\xff\xff"
	r := flate.NewReader(bytes.NewBufferString(pbmRaw))
	img, err := netpbm.Decode(r, &netpbm.DecodeOptions{
		Target:      netpbm.PPM, // Want to wind up with color
		Exact:       false,      // Can accept grayscale or B&W too
		PBMMaxValue: 42,         // B&W white --> (42, 42, 42)
	})
	if err != nil {
		panic(err)
	}
	r.Close()
	fmt.Printf("Image is of type %s.\n", img.Format())
	fmt.Printf("Maximum channel value is %d.\n", img.MaxValue())
	fmt.Printf("Color at (32, 20) is %#v.\n", img.At(32, 20))
	// Output:
	// Image is of type PPM.
	// Maximum channel value is 42.
	// Color at (32, 20) is npcolor.RGBM{R:0x2a, G:0x2a, B:0x2a, M:0x2a}.
}

func ExampleEncode() {
	// In this example, we create an 800x800 color image with a maximum
	// per-channel value of 800 and fill it with a gradient pattern.  We
	// then write the image in "plain" (ASCII) PPM format to a string and
	// output the first few lines of that string.  More typical usage would
	// be to write to a file and to output in "raw" (binary) format.

	// Create an image with a gradient pattern.
	const edge = 800 // Width, height, and maximum channel value
	img := netpbm.NewRGBM64(image.Rect(0, 0, edge, edge), edge)
	for r := 0; r < edge; r++ {
		for c := 0; c < edge; c++ {
			rgbm := npcolor.RGBM64{
				R: uint16(r),
				G: 0,
				B: uint16(c),
				M: edge,
			}
			img.SetRGBM64(c, r, rgbm)
		}
	}

	// Write the image to a string.
	var ppmBytes bytes.Buffer
	err := netpbm.Encode(&ppmBytes, img, &netpbm.EncodeOptions{
		Format:   img.Format(),   // Original format
		MaxValue: img.MaxValue(), // Original maximum value
		Plain:    true,           // ASCII output to clarify output
		Comments: []string{"Sample PPM file"},
	})
	if err != nil {
		panic(err)
	}
	ppmLines := strings.Split(ppmBytes.String(), "\n")

	// Output the first few lines of the PPM file.
	for i, line := range ppmLines {
		fmt.Println(line)
		if i == 6 {
			break
		}
	}

	// Output:
	// P3
	// # Sample PPM file
	// 800 800
	// 800
	// 0 0 0 0 0 1 0 0 2 0 0 3 0 0 4 0 0 5 0 0 6 0 0 7 0 0 8 0 0 9 0 0 10 0
	// 0 11 0 0 12 0 0 13 0 0 14 0 0 15 0 0 16 0 0 17 0 0 18 0 0 19 0 0 20 0
	// 0 21 0 0 22 0 0 23 0 0 24 0 0 25 0 0 26 0 0 27 0 0 28 0 0 29 0 0 30 0
}
