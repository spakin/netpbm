// This program generates tests for multiple image formats.
package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// A NetpbmType represents a Netpbm image type.
type NetpbmType struct {
	Fmt     string // Netpbm format (e.g., "PPM")
	BaseFmt string // Base Netpbm format (e.g., "PPM") in the case of a PAM file
	Plain   bool   // true=plain (ASCII); false=raw (binary)
	Width   int    // Expected width
	Height  int    // Expected height
	Maxval  int    // Expected maximum value
	Image   string // Test image (variable declared in netpbm_test.go)
	Suffix  string // Suffix to append to test names
}

// AllConfigs describes all Netpbm formats we want to test.
var AllConfigs = []NetpbmType{
	{Fmt: "PBM", BaseFmt: "PBM", Plain: false, Width: 64, Height: 64, Maxval: 1, Image: "pbmRaw"},
	{Fmt: "PBM", BaseFmt: "PBM", Plain: true, Width: 63, Height: 65, Maxval: 1, Image: "pbmPlain"},
	{Fmt: "PGM", BaseFmt: "PGM", Plain: false, Width: 64, Height: 64, Maxval: 255, Image: "pgmRaw"},
	{Fmt: "PGM", BaseFmt: "PGM", Plain: true, Width: 63, Height: 65, Maxval: 777, Image: "pgmPlain"},
	{Fmt: "PPM", BaseFmt: "PPM", Plain: false, Width: 64, Height: 64, Maxval: 255, Image: "ppmRaw"},
	{Fmt: "PPM", BaseFmt: "PPM", Plain: true, Width: 63, Height: 65, Maxval: 777, Image: "ppmPlain"},
	{Fmt: "PAM", BaseFmt: "PPM", Plain: false, Width: 64, Height: 64, Maxval: 255, Image: "pamRawColor"},
	{Fmt: "PAM", BaseFmt: "PGM", Plain: false, Width: 64, Height: 64, Maxval: 255, Image: "pamRawGrayAlpha", Suffix: "GrayAlpha"},
	{Fmt: "PAM", BaseFmt: "PPM", Plain: false, Width: 64, Height: 64, Maxval: 255, Image: "pamRawColorAlpha", Suffix: "Alpha"},
}

func main() {
	// Construct a template.
	allTmpl, err := template.New("all.tmpl").Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	}).ParseFiles("helpers/all.tmpl")
	if err != nil {
		panic(err)
	}

	// Write the output to a file.
	fd, err := os.Create("generated_test.go")
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	fmt.Fprintln(fd, `/*
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
)`)

	// Apply the template to a wide variety of configurations.
	for _, cfg := range AllConfigs {
		err = allTmpl.Execute(fd, cfg)
		if err != nil {
			panic(err)
		}
	}
}
