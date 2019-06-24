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
	Fmt    string // Netpbm format (e.g., "PPM")
	Plain  bool   // true=plain (ASCII); false=raw (binary)
	Width  int    // Expected width
	Height int    // Expected height
	Maxval int    // Expected maximum value
}

// AllConfigs describes all Netpbm formats we want to test.
var AllConfigs = []NetpbmType{
	{Fmt: "PBM", Plain: false, Width: 64, Height: 64, Maxval: 1},
	{Fmt: "PBM", Plain: true, Width: 63, Height: 65, Maxval: 1},
	{Fmt: "PGM", Plain: false, Width: 64, Height: 64, Maxval: 255},
	{Fmt: "PGM", Plain: true, Width: 63, Height: 65, Maxval: 777},
	{Fmt: "PPM", Plain: false, Width: 64, Height: 64, Maxval: 255},
	{Fmt: "PPM", Plain: true, Width: 63, Height: 65, Maxval: 777},
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
