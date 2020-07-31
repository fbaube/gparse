package gparse

import (
	"fmt"

	"github.com/fbaube/gtoken"
)

// XmlCheckForPreamble only prints something. It could return a flag,
// or even insert the standard XML preamble if one is not present.
func XmlCheckForPreambleToken(p []*gtoken.GToken) []*gtoken.GToken {
	if p == nil || len(p) == 0 {
		panic("Bad arg to XmlCheckForPreamble")
	}
	var pGT *gtoken.GToken
	pGT = p[0]
	var gotXmlDecl = (pGT.TTType == "PI") && (pGT.Keyword == "xml")
	if !gotXmlDecl {
		println("    --> XML preamble not found; could insert one; gtoken.xmlpreamble.L12")
	} else {
		fmt.Printf("    --> XML preamble found \n")
	}
	return p
}
