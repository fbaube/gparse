package gparse

import (
	"fmt"
	S "strings"

	SU "github.com/fbaube/stringutils"
)

// XmlPreamble is a parse of the optional PI "<?xml ..." at the top of the file.
// XML major version MUST be 1. Note that strictly speaking, it is not a PI.
//  <?xml version="version_number"         <= required, 1.x
//       encoding="encoding_declaration"   <= optional, assume "UTF-8"
//     standalone="standalone_status" ?>   <= optional, can be "yes", assume "no"
type XmlPreamble struct {
	Raw          string
	MinorVersion string // e.g. expect "0" for XML 1.0
	// Encoding: not sure what the valid values are or what form they are.
	Encoding     string
	IsStandalone bool // "yes"  or "no"
}

// NewXmlPreamble parses an XML preamble, which MUST be the first line in the file.
// Example: <?xml version="1.0" encoding='UTF-8' standalone="yes"?>
func NewXmlPreamble(s string) (*XmlPreamble, error) {
	if s == "" {
		return nil, nil
	}
	i := S.Index(s, "?>")
	if i == -1 || !S.HasPrefix(s, "<?xml ") {
		return nil, fmt.Errorf("gtoken.preamble.new.oops:", s)
	}
	s = s[:i]
	p := new(XmlPreamble)
	s = S.TrimPrefix(s, "<?xml ")
	p.Raw = S.TrimSpace(s)
	var props, sides []string
	var prop, varbl, value string
	// Break at spaces to get one to three properties.
	props = S.Split(s, " ")
	for _, prop = range props {
		sides = S.Split(prop, "=")
		varbl = sides[0]
		if !SU.IsXmlQuoted(sides[1]) {
			return p, fmt.Errorf("gtoken.preamble.new.badquotes<%s>", prop)
		}
		value = SU.MustXmlUnquote(sides[1])

		switch varbl {
		case "encoding":
			p.Encoding = value
		case "version":
			p.MinorVersion = S.TrimPrefix(value, "1.")
			if "0" != p.MinorVersion {
				panic("Bad XML version number<" + p.MinorVersion + ">")
			}
		case "standalone":
			// We may safely ignore bogus values.
			p.IsStandalone = (value == "yes")
		}
	}
	return p, nil
}

func (xp XmlPreamble) Echo() string { return xp.Raw + "\n" }

func (xp XmlPreamble) String() string {
	var xmlver, encodg, stdaln string
	xmlver = fmt.Sprintf("<?xml version=\"1.%s\"", xp.MinorVersion)
	if xp.Encoding != "" {
		encodg = fmt.Sprintf(" encoding=\"%s\"", xp.Encoding)
	}
	if xp.IsStandalone {
		stdaln = " standalone=\"yes\""
	}
	return xmlver + encodg + stdaln + ">?\n"
}

// XmlCheckForPreamble only prints something. It could return a flag,
// or even insert the standard XML preamble if one is not present.
func XmlCheckForPreambleToken(p []*GToken) []*GToken {
	if p == nil || len(p) == 0 {
		panic("Bad arg to XmlCheckForPreamble")
	}
	if !ExtraInfo {
		return p
	}
	var pGT *GToken
	pGT = p[0]
	var gotXmlDecl = (pGT.TTType == "PI") && (pGT.Keyword == "xml")
	if !gotXmlDecl {
		println("    --> XML preamble not found; could insert one; gtoken.xmlpreamble.L12")
	} else {
		fmt.Printf("    --> XML preamble found \n")
	}
	return p
}
