// Present a collection of examples to demonstrate netpbm package usage.

package netpbm_test

import (
	"bytes"
	"compress/flate"
	"fmt"
	"github.com/spakin/netpbm"
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
		Target:      netpbm.PPM,  // Want to wind up with color
		Exact:       false,       // Can accept grayscale or B&W too
		PBMMaxValue: 42,          // B&W white --> (42, 42, 42)
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
