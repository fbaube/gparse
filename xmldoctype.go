package gparse

import (
	"fmt"
	S "strings"

	SU "github.com/fbaube/stringutils"
	"github.com/pkg/errors"
)

// This file contains LwDITA-specific stuff, but it is hard-coded and
// does not pull in other packages, so we leave it alone for now.

var knownRootTags = []string{"html", "map", "topic", "task", "concept", "reference"}

// Copied from gfile.go:
// [0] XML, BIN, TXT, MD
// [1] IMG, CNT (Content), TOC (Map), SCH(ema)
// [2] XML: per-DTD; BIN: fmt/filext; MD: flavor; SCH: fmt/filext

// type XmlDoctypeFamily string
//      XmlDoctypeFamilies are the broad groups of DOCTYPES.
//  var XmlDoctypeFamilies = []XmlDoctypeFamily {
//	"lwdita",
//	"dita",
//	"html5",
//	"html",
//	"other",
// }

// XmlDoctype is a parse of a DOCTYPE declaration.
// For [Lw]DITA, what interests us is
// PUBLIC "-//OASIS//DTD (PublicTextDesc)//EN" or sometimes
// PUBLIC "-//OASIS//ELEMENTS (PublicTextDesc)//EN" or sometimes
// :ul:
// :: PUBLIC | SYSTEM = Avbty
// :: - = Reg'n = Org'zn & DTD are not reg'd with ISO.
// :: OASIS = Org'zn
// :: DTD = Public Text Class (CAPACITY | CHARSET | DOCUMENT |
// DTD | ELEMENTS | ENTITIES | LPD | NONSGML | NOTATION |
// SHORTREF | SUBDOC | SYNTAX | TEXT )
// :: (*) = Public Text Description, incl. any version number
// :: EN = Public Text Language
// :: URL = optional, explicit
// -ul-
type XmlDoctype struct {
	raw string
	// PUBLIC
	Availability string
	// ??

	// NOTE this is dangerous
	XmlPublicID
	// and SHOULD BE replaced by this
	// rawXmlPublicID string

	// TopTag is the tag declated in the DOCTYPE
	TopTag string
	// Micodo is here because a DOCTYPE does indeed let us create one.
	Mtype []string
}

// NewXmlDoctypeInclMtype should work for a normal XML document DOCTYPE
// declaration, but it should also work on a DOCTYPE reference plucked
// out of a DTD file, which tells the user what DOCTYPE declaration to
// use to reference the DTD. In other words, the XML Catalog reference.
// Therefore, let this function parse a string that begins with PUBLIC
// or SYSTEM, and don't worry about the how exactly the string ends.
// Target strings of great interest:
// DOCTYPE topic PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN"
// DOCTYPE map   PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Map//EN"
// DOCTYPE html (HTML5)
// DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" (MAYBE!)
// NOTE that Mtype[] is non-nil but each element can be "", the empty string.
func NewXmlDoctypeInclMtype(s string) (*XmlDoctype, error) {
	if s == "" {
		return nil, nil
	}
	s = S.TrimSpace(s)
	var pDT = new(XmlDoctype)
	pDT.Mtype = make([]string, 3)

	pDT.raw = s
	if ExtraInfo {
		println("    --> Processing raw doctype:\n        ", s)
	}

	// Quick exit: HTML5
	if S.EqualFold(s, "html") {
		pDT.Mtype[0] = "xml"
		pDT.Mtype[1] = "cnt"
		pDT.TopTag = "html"
		pDT.Mtype[2] = "html5"
		pDT.PublicTextClass = "(HTML5)"
		return pDT, nil
	}
	if S.Contains(s, "HTML 4") {
		pDT.Mtype[0] = "html"
		pDT.Mtype[1] = "html4"
		pDT.TopTag = "html"
		pDT.Mtype[2] = "html"
		return pDT, nil
	}
	// topic PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN" "lw-topic.dtd"
	pDT.TopTag, s = SU.SplitOffFirstWord(s)
	if !SU.IsInSliceIgnoreCase(pDT.TopTag, knownRootTags) {
		println("    --> Unrecognized DOCTYPE element:", pDT.TopTag)
	}
	pDT.Availability, s = SU.SplitOffFirstWord(s)
	if pDT.Availability != "PUBLIC" && pDT.Availability != "SYSTEM" {
		return nil, errors.New("Bad DOCTYPE availability<" +
			pDT.Availability + "> (neither PUBLIC nor SYSTEM)")
	}
	// "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN" "lw-topic.dtd"
	var e error
	// It has to be a quoted string, but it might actually be two of them
	// (both Public ID and System ID) and they might even use different
	// types of quotes.
	// FIXME Handle cases of bad quoting.
	qtd1, qtd2, e := SU.SplitOffQuotedToken(s)
	if e != nil {
		return pDT, fmt.Errorf("gtoken.doctype.newInclM.SplitOffQuotedToken<%s>", s)
	}
	qtd2 = S.TrimSpace(qtd2)
	if qtd2 != "" {
		if !SU.IsXmlQuoted(qtd2) {
			return pDT, fmt.Errorf("gtoken.doctype.newInclM.SplitOffQuotedToken<%s>", s)
		}
		qtd2 = SU.MustXmlUnquote(qtd2)
	}
	if pDT.Availability == "SYSTEM" {
		if qtd2 != "" {
			return pDT, fmt.Errorf("gtoken.doctype.newInclM.SecondSYSTEMargument<%s>", qtd2)
		}
		pDT.SystemID = SystemID(qtd1)
	} else {
		// "PUBLIC"
		ppid, e := NewXmlPublicID(qtd1)
		if e != nil {
			return nil, errors.Wrapf(e, "gtoken.doctype.newInclM.NewXmlPublicID<%s>", qtd1)
		}
		pDT.XmlPublicID = *ppid
		pDT.SystemID = SystemID(qtd2)
	}

	sd := pDT.XmlPublicID.PublicTextClass
	if sd == "" {
		return pDT, nil
	}
	if S.Contains(sd, "DITA") {
		pDT.Mtype[0] = "dita"
		pDT.Mtype[1] = "[TBS]"
		pDT.Mtype[2] = pDT.TopTag
	}
	if S.Contains(sd, "XDITA") ||
		S.Contains(sd, "LW DITA") ||
		S.Contains(sd, "LIGHTWEIGHT DITA") {
		pDT.Mtype[0] = "lwdita"
		pDT.Mtype[1] = "xdita"
		pDT.Mtype[2] = pDT.TopTag
	}
	return pDT, nil
}

func (xd XmlDoctype) Echo() string { return xd.raw + "\n" }

func (xd XmlDoctype) String() string {
	if "" == xd.TopTag {
		panic("xd.TopTag")
	}
	/*
		// TODO Also print System ID, if present
		if xd.MicodoMainType == "html" && xd.MicodoSubType == "html5 "{
			return "<!DOCTYPE html>\n"
		}
		return fmt.Sprintf("<!DOCTYPE %s %s %s><!--mcd:%s:%s-->\n",
			xd.TopTag, xd.Availability, xd.XmlPublicID.String(),
			xd.MicodoMainType, xd.MicodoSubType)
	*/
	return "    gtoken.XmlDoctype.S: " + xd.Echo()
}
