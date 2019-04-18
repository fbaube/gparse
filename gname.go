package gparse

// This file: Structures for Generic Golang XML Names.
// Struct `GName` is a renaming of struct `xml.Name`.

import (
	"encoding/xml"
	"io"
	S "strings"
)

// GName is a generic golang XML name.
//
// NOTE If `GName.Name` (i.e. the namespace part, not the `Local`
// part) is non-nil, then ALWAYS include a trailing semicolon in it.
// This *greatly* simplifies output generation.
//
// Structure details of `xml.Name`:
//   type Name struct { Space, Local string }
//
type GName xml.Name

// Echo implements Markupper.
func (N GName) Echo() string {
	// if N.Space == "" {
	// 	return N.Local
	// }
	// Note that we assume the colon is there at the end of `N.Space`
	return N.Space + N.Local
}

// EchoTo implements Markupper.
func (N GName) EchoTo(w io.Writer) {
	w.Write([]byte(N.Echo()))
}

// String implements Markupper.
func (N GName) String() string {
	return N.Echo()
}

// DumpTo implements Markupper.
func (N GName) DumpTo(w io.Writer) {
	w.Write([]byte(N.String()))
}

// NewGName adds a colon to a non-empty namespace if it is not there already.
func NewGName(ns, local string) *GName {
	p := new(GName)
	if ns != "" && !S.HasSuffix(ns, ":") {
		ns += ":"
	}
	p.Space = ns
	p.Local = local
	return p
}
