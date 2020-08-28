package gparse

import "io"

// MarkupStringer is an interface meant to be a souped-up version of
// [GoStringer](https://golang.org/pkg/fmt/#GoStringer) for markup
// fragments in any of the formats we process ("XML", "Markdown",
// "HTML", future TBS) that can be string'ified in two different ways:
// - `Echo` basically recreates the original input, although in
// a normalized format; the name is more specific than `String`.
// - `String` returns a format suitable for development and debugging,
// but (probably) not valid markup: `string String()`, `DumpTo(w)`
//
type MarkupStringer interface { // extends fmt.Stringer
	// Echo the input markup, (but) in a normalized format
	Echo() string
	EchoTo(io.Writer)
	// For development & debugging - probably not valid markup
	String() string
	DumpTo(io.Writer)
}
