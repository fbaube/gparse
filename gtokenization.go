package gparse

// This file: Structures for Generic Golang Tokens,
// based on the tokens returns by the Golang XML parser.
// Note that `GTokenization` does *not* implement `Markupper`.

import (
	"fmt"
	"io"
	// "github.com/dimchansky/utfbom"
)

// GTokenization is defined solely for the convenience methods defined below.
type GTokenization []*GToken

func (inGTzn GTokenization) DeleteNils() (outGTzn GTokenization) {
	if nil == inGTzn || len(inGTzn) == 0 {
		return nil
	}
	for _, pGT := range inGTzn {
		if nil != pGT {
			outGTzn = append(outGTzn, pGT)
		}
	}
	return outGTzn
}

// DumpTo writes out the `GToken`s to the `io.Writer`, one per line, and each
// line is prefixed with the token type. The output should parse the same as
// the input file, except perhaps for the treatment of all-whitespace CDATA.
func (GTzn GTokenization) DumpTo(w io.Writer) {
	if nil == GTzn || nil == w {
		return
	}
	// GTzn = GTzn.DeleteNils()
	var pGT *GToken
	var i int

	for i, pGT = range GTzn {
		if nil == pGT {
			continue
		}
		fmt.Fprintf(w, "<!--[%02d]%s--> %s \n",
			i, pGT.TTType.LongForm(), pGT.Echo())
	}
}

func (GTokzn GTokenization) HasDoctype() (bool, string) {
	if nil == GTokzn || len(GTokzn) == 0 {
		return false, ""
	}
	var pGT *GToken
	for _, pGT = range GTokzn {
		switch pGT.TTType {
		case "Dir":
			return true, pGT.Otherwords
		}
	}
	return false, ""
}

// GetFirstByTag checks the basic tag only, not any namespace.
func (gTkzn GTokenization) GetFirstByTag(s string) *GToken {
	if s == "" {
		return nil
	}
	for _, p := range gTkzn {
		if p.GName.Local == s && p.TTType == "SE" {
			return p
		}
	}
	return nil
}

// GetAllByTag returns a new GTokenization.
// It checks the basic tag only, not any namespace.
func (gTkzn GTokenization) GetAllByTag(s string) GTokenization {
	if s == "" {
		return nil
	}
	// fmt.Printf("GetAllByTag<%s> len:%d \n", s, len(gTkzn))
	var ret GTokenization
	ret = make(GTokenization, 0)
	for _, p := range gTkzn {
		if p.GName.Local == s && p.TTType == "SE" {
			// fmt.Printf("found a match [%d] %s (NS:%s)\n", i, p.GName.Local, p.GName.Space)
			ret = append(ret, p)
		}
	}
	return ret
}

func (gTkzn GTokenization) DString() {
	i := len(gTkzn)
	fmt.Printf("GTokenization len<%d>\n", i)
	for i, GT := range gTkzn {
		fmt.Printf("  [%2d] %s \n", i, GT.Echo())
	}
}
