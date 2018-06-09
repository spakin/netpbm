netpbm
======

Introduction
------------

`netpbm` is a package for the [Go programming language](http://www.golang.org/) that implements image decoders and encoders for the [Netpbm image formats](http://netpbm.sourceforge.net/doc/#formats).  The package supports all of the following:

* All of the Netpbm image formats:

  - [PBM](http://netpbm.sourceforge.net/doc/pbm.html) (portable bitmap): black and white only
  - [PGM](http://netpbm.sourceforge.net/doc/pgm.html) (portable graymap): grayscale
  - [PPM](http://netpbm.sourceforge.net/doc/ppm.html) (portable pixmap): color
  - [PAM](http://netpbm.sourceforge.net/doc/pam.html) (portable arbitrary map): alpha

* Both "raw" (binary) and "plain" (ASCII) files

* Both 8-bit and 16-bit color channels

* Any maximum per-color-channel value (up to what the given number of bits can represent)

* Full compatibility with Go's [`image`](https://golang.org/pkg/image/) package

  - Implements the [`image.Image`](https://golang.org/pkg/image/#Image) interface
  - Additionally defines `Opaque`, `PixOffset`, `Set`, and `Subimage` methods (and color-model-specific variants of `At` and `Set`), like most of Go's standard image types

* Automatic promotion of image formats, if desired

That last feature means that a program that expects to read a grayscale image can also be given a black-and-white image, and a program that expects to read a color image can also be given either a grayscale or a black-and-white image.

Installation
------------

Instead of manually downloading and installing `netpbm` from GitHub, the recommended approach is to ensure your `GOPATH` environment variable is set properly then issue a

    go get github.com/spakin/netpbm

command.

Usage
-----

`netpbm` works just like the standard [`image/gif`](https://golang.org/pkg/image/gif/), [`image/jpeg`](https://golang.org/pkg/image/jpeg/), and [`image/png`](https://golang.org/pkg/image/png/) packages in that 

    import (
        _ "github.com/spakin/netpbm"
    )

will enable [`image.Decode`](https://golang.org/pkg/image/#Decode) to import Netpbm image formats.

Various package-specific functions, types, interfaces, and methods are available only with a normal (not "`_`") `import`.  A normal `import` is needed both to export Netpbm images and to exert more precise control over the Netpbm variants that are allowed to be imported.  See the [`netpbm` API documentation](http://godoc.org/github.com/spakin/netpbm) for details.

Author
------

[Scott Pakin](http://www.pakin.org/~scott/), *scott+npbm@pakin.org*
