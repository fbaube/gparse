package gparse

import "io"

// init checks that the structs that we *want*
// to implement Markupper *do* implement it.
func init() {

}

// Markupper is for markup fragments (XML, Markdown, etc.)
// that can be string'ified in two ways:
// - in an fromat that is nearly an echo of the input,
// but is in a normalized format (`Echo,EchoTo`)
// - in a format that is suitable for development
// and debugging, but which is (probably) not valid
// markup (`String,DumpTo`)
type Markupper interface {
	// Echo the input markup, but in a normalized format
	Echo() string
	EchoTo(io.Writer)
	// For development & debugging - probably not valid markup
	String() string
	DumpTo(io.Writer)
}
