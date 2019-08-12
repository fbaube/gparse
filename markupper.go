package gparse

import (
	"encoding/xml"
	"fmt"
	"io"
	// "github.com/fbaube/gfile"
)

// init checks that the structs that we *want*
// to implement Markupper *do* implement it.
func init() {
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
	println(MU)
	// GL := gfile.GLink{}
	// MU = GL
	// fmt.Printf("GLnk OK: %s // %#v \n", MU.Echo(), MU.(GToken))
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
