package gparse

import (
	"encoding/xml"
	"io"
)

// This file: Generic Golang XML Attributes.
// Struct `GAtt` is a renaming of struct `xml.Attr`.

// NOTE: In LwDITA, the `class` attribute can have more than one value,
// separated by space, like this:
//   <p class="a b c">Alice In Wonderland</p>
// Order does not matter.
// You should NOT use multiple `class`, such as `class="..." class="..."``

// GAtt is a generic golang XML attribute.
//
// Structure details of `xml.Attr`:
//   type Attr struct {
//     // xml.Name :: Space, Local string
//     Name  Name
//     Value string }
//
// NOTE The related struct `DAtt` drops `Value`,
// and adds `AttType,AttDflt string`
//
type GAtt xml.Attr

// GAttList is TODO? Replace with a map
type GAttList []*GAtt

// Echo implements Markupper (and inserts a leading space).
func (A GAtt) Echo() string {
	return " " + GName(A.Name).Echo() + "=\"" + A.Value + "\""
}

// Echo implements Markupper (and inserts a leading space).
func (AL GAttList) Echo() string {
	var s string
	for _, A := range AL {
		s += " " + GName(A.Name).Echo() + "=\"" + A.Value + "\""
	}
	return s
}

// EchoTo implements Markupper.
func (A GAtt) EchoTo(w io.Writer) {
	w.Write([]byte(A.Echo()))
}

// EchoTo implements Markupper.
func (AL GAttList) EchoTo(w io.Writer) {
	w.Write([]byte(AL.Echo()))
}

// String implements Markupper.
func (A GAtt) String() string {
	return A.Echo()
}

// String implements Markupper.
func (AL GAttList) String() string {
	return AL.Echo()
}

// DumpTo implements Markupper.
func (A GAtt) DumpTo(w io.Writer) {
	w.Write([]byte(A.String()))
}

// DumpTo implements Markupper.
func (AL GAttList) DumpTo(w io.Writer) {
	w.Write([]byte(AL.String()))
}

// GetAttVal returns the attribute's string value, or "" if not found.
func (p GAttList) GetAttVal(att string) string {
	for _, A := range p {
		if A.Name.Local == att {
			return A.Value
		}
	}
	return ""
}
