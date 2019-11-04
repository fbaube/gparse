package gparse

import "io"

// init checks that the structs that we *want*
// to implement Markupper *do* implement it.
func init() {
	/*
		var MU Markupper

		GN := GName{"ns:", "lcl"}
		MU = GN
		fmt.Printf("GNam OK as Markupper\n") // %s // %#v \n", MU.Echo(), MU.(GName))

		GA := GAtt{xml.Name{"myns:", "myprop"}, "itsval"}
		MU = GA
		fmt.Printf("GAtt OK as Markupper\n") // %s // %#v \n", MU.Echo(), MU.(GAtt))

		GAL := GAttList{GA, GA} // {&GA, &GA}
		GAL[1].Value += "2"
		GAL[1].Name.Local += "2"
		MU = GAL
		fmt.Printf("GAtL OK as Markupper\n") // %s // %#v \n", MU.Echo(), MU.(GAttList))

		GT := GToken{nil, "SE", GN, GAL, "Kwd", "Othwds"}
		MU = GT
		fmt.Printf("GTkn OK as Markupper\n") // %s // %#v \n", MU.Echo(), MU.(GToken))

		println("(obligatory Markupper)", MU)
		// GL := gfile.GLink{}
		// MU = GL
		// fmt.Printf("GLnk OK: %s // %#v \n", MU.Echo(), MU.(GToken))
	*/
}

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
